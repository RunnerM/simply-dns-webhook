package client

import (
	"fmt"
	"testing"
)

var fixture SimplyClient

type testData struct {
	domain      string
	data        string
	data2       string
	accountname string
	apikey      string
	basedomain  string
}

// Plot in your own api details for testing.
func TestAll(t *testing.T) {
	data := testData{ //add your credentials here to test.
		domain:      "_acme-challenge.foo.com",
		data:        "test_txt_data",
		data2:       "test_txt_data_2",
		accountname: "",
		apikey:      "",
	}
	testAdd(t, data)
	id := testGet(t, data)
	testUpdate(t, data, id)
	testRemove(t, data, id)

}

func testAdd(t *testing.T, data testData) {
	id, err := fixture.AddTxtRecord(data.domain, data.data, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if err != nil {
		t.Fail()
	}
	if id == 0 {
		t.Fail()
	}
	fmt.Println(id)
}

func testUpdate(t *testing.T, data testData, id int) {
	res, err := fixture.UpdateTXTRecord(id, data.domain, data.data2, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if err != nil {
		t.Fail()
	}
	if res != true {
		t.Fail()
	}
	fmt.Println(id)
}

func testRemove(t *testing.T, data testData, id int) {
	res2, _ := fixture.GetExactTxtRecord(data.data2, data.domain, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})

	if res2 != id {
		t.Fail()
	}

	res := fixture.RemoveTxtRecord(id, data.domain, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if res != true {
		t.Fail()
	}

}
func testGet(t *testing.T, data testData) int {
	id, recData, _ := fixture.GetTxtRecord(data.domain, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if id == 0 {
		t.Fail()
	}
	if recData == "" {
		t.Fail()
	}
	return id
}
