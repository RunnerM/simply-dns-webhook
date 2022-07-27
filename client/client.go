package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// SimplyClient base type
type SimplyClient struct {
	Baseurl     string
	ApiKey      string
	AccountName string
	Domain      string
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

// AddTxtRecord Add txt record to simply
func (c *SimplyClient) AddTxtRecord(SubDomain string, Value string) int {
	TXTRecordBody := CreateRecordBody{
		Type:     "TXT",
		Name:     SubDomain,
		Data:     Value,
		Priority: 1,
		Ttl:      3600,
	}
	postBody, _ := json.Marshal(TXTRecordBody)

	req, error := http.NewRequest("POST", c.Baseurl+"/my/products/"+c.Domain+"/dns/records", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(c.AccountName, c.ApiKey)
	client := &http.Client{}
	response, error := client.Do(req)

	if error != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", error, " response: ", response.StatusCode)
	}
	responseData, error := ioutil.ReadAll(response.Body)

	if error != nil {
		fmt.Println("Error on read: ", error)
	}
	var data CreateRecordResponse

	error = json.Unmarshal(responseData, &data)
	if error != nil {
		panic(error)
	}
	return data.Record.Id
}

// RemoveTxtRecord Remove TXT record from symply
func (c *SimplyClient) RemoveTxtRecord(RecordId int) bool {
	req, error := http.NewRequest("DELETE", c.Baseurl+"/my/products/"+c.Domain+"/dns/records/"+strconv.Itoa(RecordId), nil)
	req.SetBasicAuth(c.AccountName, c.ApiKey)
	client := &http.Client{}
	response, error := client.Do(req)

	if error != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", error, " response: ", response.StatusCode)
		return false
	} else {
		return true
	}
}

// GetTxtRecord Fetch TXT record by data returns id
func (c *SimplyClient) GetTxtRecord(TxtData string) int {
	req, error := http.NewRequest("GET", c.Baseurl+"/my/products/"+c.Domain+"/dns/records", nil)
	req.SetBasicAuth(c.AccountName, c.ApiKey)
	client := &http.Client{}
	response, error := client.Do(req)

	if error != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", error, " response: ", response.StatusCode)
	}
	responseData, error := ioutil.ReadAll(response.Body)

	if error != nil {
		fmt.Println("Error on read: ", error)
	}

	var records RecordResponse

	error = json.Unmarshal(responseData, &records)
	var recordId int

	if error == nil {
		for i := 0; i < len(records.Records); i++ {
			if records.Records[i].Data == TxtData {
				recordId = records.Records[i].RecordId
			}
		}
	} else {
		panic(error)
	}

	return recordId
}
