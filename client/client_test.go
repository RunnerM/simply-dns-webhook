package client

import (
	"fmt"
	"testing"
)

var fixture SimplyClient

type testData struct {
	domain      string
	data        string
	accountname string
	apikey      string
}

//Plot in your own api details for testing.
func TestAll(t *testing.T) {
	data := testData{
		domain:      ".com", //add your credentials here to test.
		data:        "",
		accountname: "",
		apikey:      "",
	}
	testAdd(t, data)
	//testGet(t, data)
	//testRemove(t, data, 0)

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
func testRemove(t *testing.T, data testData, id int) {
	res := fixture.RemoveTxtRecord(id, data.domain, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if res != true {
		t.Fail()
	}

}
func testGet(t *testing.T, data testData) int {
	id, _ := fixture.GetTxtRecord(data.data, data.domain, Credentials{
		AccountName: data.accountname,
		ApiKey:      data.apikey,
	})
	if id == 0 {
		t.Fail()
	}
	return id
}
