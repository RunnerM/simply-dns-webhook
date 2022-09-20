package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	apiUrl = "https://api.simply.com/2"
)

// SimplyClient base type
type SimplyClient struct {
}

// RecordResponse api type
type RecordResponse struct {
	Records []struct {
		RecordId int    `json:"record_id"`
		Name     string `json:"name"`
		Ttl      int    `json:"ttl"`
		Data     string `json:"data"`
		Type     string `json:"type"`
		Priority int    `json:"priority"`
	} `json:"records"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// CreateRecordBody api type
type CreateRecordBody struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	Priority int    `json:"priority"`
	Ttl      int    `json:"ttl"`
}

// CreateRecordResponse api type
type CreateRecordResponse struct {
	Record struct {
		Id int `json:"id"`
	} `json:"record"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Credentials struct {
	AccountName string `json:"status"`
	ApiKey      string `json:"message"`
}

// AddTxtRecord Add txt record to simply
func (c *SimplyClient) AddTxtRecord(Domain string, Value string, credentials Credentials) (int, error) {
	TXTRecordBody := CreateRecordBody{
		Type:     "TXT",
		Name:     domainutil.Subdomain(Domain),
		Data:     Value,
		Priority: 1,
		Ttl:      3600,
	}
	postBody, _ := json.Marshal(TXTRecordBody)
	fmt.Println("adding dns record: ")
	fmt.Println(postBody)

	req, err := http.NewRequest("POST", apiUrl+"/my/products/"+domainutil.Domain(Domain)+"/dns/records", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	fmt.Println("name: ", credentials.AccountName, "  key: ", credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", err, " response: ", response.StatusCode)
		return 0, err
	}
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error on read: ", err)
	}
	var data CreateRecordResponse

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		panic(err)
	}
	return data.Record.Id, nil
}

// RemoveTxtRecord Remove TXT record from symply
func (c *SimplyClient) RemoveTxtRecord(RecordId int, DnsName string, credentials Credentials) bool {
	req, err := http.NewRequest("DELETE", apiUrl+"/my/products/"+domainutil.Domain(DnsName)+"/dns/records/"+strconv.Itoa(RecordId), nil)
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	fmt.Println(credentials.AccountName, "  ", credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", err, " response: ", response.StatusCode)
		return false
	} else {
		return true
	}
}

// GetTxtRecord Fetch TXT record by data returns id
func (c *SimplyClient) GetTxtRecord(TxtData string, DnsName string, credentials Credentials) int {
	req, err := http.NewRequest("GET", apiUrl+"/my/products/"+domainutil.Domain(DnsName)+"/dns/records", nil)
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", err, " response: ", response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error on read: ", err)
	}

	var records RecordResponse

	err = json.Unmarshal(responseData, &records)
	var recordId int

	if err == nil {
		for i := 0; i < len(records.Records); i++ {
			if records.Records[i].Data == TxtData {
				recordId = records.Records[i].RecordId
			}
		}
	} else {
		panic(err)
	}

	return recordId
}
