package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
)

func Test_correctURL(t *testing.T) {
	RunningTest = true

	server := httptest.NewServer(http.HandlerFunc(handlerProjectinfo))
	defer server.Close()

	var fileData ExportInfo
	var httpData ExportInfo

	file, err := ioutil.ReadFile("output.json")
	if err != nil {
		t.Error("Test failed to read file \"output.json\" | Error: " + err.Error())
	} else {
		err = json.Unmarshal(file, &fileData)
		if err != nil {
			t.Error("Test failed to unmarshal data loaded from \"output.json\" | Error: " + err.Error())
		}
	}

	resp, err := http.Get(server.URL + "/projectinfo/v1/github.com/Xillez/Test")
	if err != nil {
		t.Error("Test failed to fetch from \"" + server.URL + "/projectinfo/v1/github.com/Xillez/Test\" | Error: " + err.Error())
	} else {
		err = json.NewDecoder(resp.Body).Decode(&httpData)
		if err != nil {
			t.Error("Test failed to decode data fetched from \"" + server.URL + "\" | Error: " + err.Error())
		}
	}

	sort.Strings(httpData.Langs)
	sort.Strings(fileData.Langs)

	if httpData.Name != fileData.Name {
		t.Error("Names not equal")
	}
	if httpData.Owner != fileData.Owner {
		t.Error("Owner not equal")
	}
	if httpData.Contrib != fileData.Contrib {
		t.Error("Top contributor not equal")
	}
	if httpData.Commits != fileData.Commits {
		t.Error("Nr of commit from top contributor not equal")
	}
	for i := 0; i < len(httpData.Langs); i++ {
		if httpData.Langs[i] != fileData.Langs[i] {
			t.Error("\"" + httpData.Langs[i] + "\" <- Does not equal -> \"" + fileData.Langs[i] + "\"")
		}
	}
}

func Test_partMissingURL(t *testing.T) {
	RunningTest = false
	server := httptest.NewServer(http.HandlerFunc(handlerProjectinfo))
	defer server.Close()

	resp, err := http.Get(server.URL + "/projectinfo/v1/github.com/Xillez")
	if err != nil {
		t.Error("Test failed to fetch from \"" + server.URL + "/projectinfo/v1/github.com/Xillez\" | Error: " + err.Error())
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Statuscode isn't \"" + http.StatusText(http.StatusBadRequest) + "\" / 400")
	}
}

func Test_noncompleteURL(t *testing.T) {
	RunningTest = false
	server := httptest.NewServer(http.HandlerFunc(handlerProjectinfo))
	defer server.Close()

	resp, err := http.Get(server.URL + "/projectinfo/v1/github.com/Xilez/Test")

	if err != nil {
		t.Error("Test failed to fetch from \"" + server.URL + "/projectinfo/v1/github.com/Xillez/Test\" | Error: " + err.Error())
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Error("Statuscode isn't \"" + http.StatusText(http.StatusNotFound) + "\" / " + strconv.Itoa(http.StatusNotFound))
	}
}
