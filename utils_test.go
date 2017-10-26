package main

import (
	"testing"
)

func Test_something1(t *testing.T) {
	RunningTest = true

	/*server := httptest.NewServer(http.HandlerFunc(handlerProjectinfo))
	defer server.Close()

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

	if httpData.Name != fileData.Name {
		t.Error("Names not equal")
	}*/
}
