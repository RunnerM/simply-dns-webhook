package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bobesa/go-domain-util/domainutil"
	"io"
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

// CreateUpdateRecordBody api type
type CreateUpdateRecordBody struct {
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
func (c *SimplyClient) AddTxtRecord(FQDNName string, Value string, credentials Credentials) (int, error) {
	// Trim one trailing dot
	fqdnName := cutTrailingDotIfExist(FQDNName)
	TXTRecordBody := CreateUpdateRecordBody{
		Type:     "TXT",
		Name:     domainutil.Subdomain(fqdnName),
		Data:     Value,
		Priority: 1,
		Ttl:      3600,
	}
	postBody, _ := json.Marshal(TXTRecordBody)
	req, err := http.NewRequest("POST", apiUrl+"/my/products/"+domainutil.Domain(fqdnName)+"/dns/records", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		fmt.Println("Error on request: ", err, " response: ", response.StatusCode)
		return 0, err
	}
	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println("Error on read: ", err)
		return 0, err
	}
	var data CreateRecordResponse

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		fmt.Println("Error on unmarshalling: ", err)
		return 0, err
	}
	return data.Record.Id, nil
}

// RemoveTxtRecord Remove TXT record from symply
func (c *SimplyClient) RemoveTxtRecord(RecordId int, DnsName string, credentials Credentials) bool {
	dnsName := cutTrailingDotIfExist(DnsName)
	req, err := http.NewRequest("DELETE", apiUrl+"/my/products/"+domainutil.Domain(dnsName)+"/dns/records/"+strconv.Itoa(RecordId), nil)
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		_ = fmt.Errorf("error on request(DELETE record): %v response: %d", err, response.StatusCode)
		return false
	} else {
		return true
	}
}

// GetTxtRecord Fetch TXT record by data returns id
func (c *SimplyClient) GetTxtRecord(FQDNName string, credentials Credentials) (int, string, error) {
	fqdnName := cutTrailingDotIfExist(FQDNName)
	req, err := http.NewRequest("GET", apiUrl+"/my/products/"+domainutil.Domain(fqdnName)+"/dns/records", nil)
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		_ = fmt.Errorf("error on request(GET record): %v response: %d", err, response.StatusCode)
		return 0, "", err
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		_ = fmt.Errorf("error on read: %v", err)
		return 0, "", err
	}

	var records RecordResponse
	err = json.Unmarshal(responseData, &records)
	var recordId int
	var recordData string

	if err == nil {
		for i := 0; i < len(records.Records); i++ {
			if records.Records[i].Type == "TXT" && records.Records[i].Name == domainutil.Subdomain(fqdnName) {
				recordId = records.Records[i].RecordId
				recordData = records.Records[i].Data
				return recordId, recordData, nil
			}
		}
	} else {
		_ = fmt.Errorf("error on fecthing records: %v", err)
		return 0, "", err
	}
	return 0, "", nil
}

// GetExactTxtRecord Fetch TXT record by data returns id of exact record
func (c *SimplyClient) GetExactTxtRecord(TxtData string, FQDNName string, credentials Credentials) (int, error) {
	fqdnName := cutTrailingDotIfExist(FQDNName)
	req, err := http.NewRequest("GET", apiUrl+"/my/products/"+domainutil.Domain(fqdnName)+"/dns/records", nil)
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil || response.StatusCode != 200 {
		_ = fmt.Errorf("error on request(GET record): %v response: %d", err, response.StatusCode)
		return 0, err
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		_ = fmt.Errorf("error on read: %v", err)
		return 0, err
	}

	var records RecordResponse
	err = json.Unmarshal(responseData, &records)
	var recordId int

	if err == nil {
		for i := 0; i < len(records.Records); i++ {
			if records.Records[i].Data == TxtData && records.Records[i].Name == domainutil.Subdomain(fqdnName) {
				recordId = records.Records[i].RecordId

				return recordId, nil
			}
		}
	} else {
		_ = fmt.Errorf("error on fecthing records: %v", err)
		return 0, err
	}
	return 0, nil
}

func (c *SimplyClient) UpdateTXTRecord(RecordId int, FQDNName string, Value string, credentials Credentials) (bool, error) {
	// Trim one trailing dot
	fqdnName := cutTrailingDotIfExist(FQDNName)
	TXTRecordBody := CreateUpdateRecordBody{
		Type:     "TXT",
		Name:     domainutil.Subdomain(fqdnName),
		Data:     Value,
		Priority: 1,
		Ttl:      3600,
	}
	putBody, _ := json.Marshal(TXTRecordBody)
	req, err := http.NewRequest("PUT", apiUrl+"/my/products/"+domainutil.Domain(fqdnName)+"/dns/records/"+strconv.Itoa(RecordId), bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(credentials.AccountName, credentials.ApiKey)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		_ = fmt.Errorf("error on request(PUT Record): %v response: %d", err, response.StatusCode)
		return false, err
	}
	return true, nil
}

func cutTrailingDotIfExist(FQDNName string) string {
	fqdnName := FQDNName
	if last := len(fqdnName) - 1; last >= 0 && fqdnName[last] == '.' {
		fqdnName = fqdnName[:last]
	}
	return fqdnName
}
