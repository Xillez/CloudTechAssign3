package main

import (
	"Assignment2/mongodb"
	"Assignment2/types"
	"Assignment2/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var testWebhookPos = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "EUR", "NOK", 0.5, 2.0}
var testWebhookNeg = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "CAD", "NOK", 0.5, 2.0}
var testCurrPos = types.CurrencyInfo{bson.NewObjectId(), "EUR", "2017-01-01", map[string]float64{"NOK": 2.0}}
var client = http.Client{}

func Test_TestingSetup(t *testing.T) {
	DB = &mongodb.MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
	DB.Init()
}

// Positive test, ProcAddWebhook
func Test_Pos_ProcAddWebhook(t *testing.T) {
	recorder := httptest.NewRecorder()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAddWebHook(r, recorder), recorder)
	}))
	defer server.Close()

	byteSlice, _ := json.Marshal(testWebhookPos)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)

	req, err := http.NewRequest("POST", server.URL+"/exchange", reader)
	if err != nil {
		t.Error("Failed to construct new request to server! | Error: " + err.Error())
	}

	_, err = client.Do(req)
	if err != nil {
		t.Error("Failed to POST to  \"" + server.URL + "\" | Error: " + err.Error())
	}

	if strings.Contains(recorder.Body.String(), testWebhookPos.ID.Hex()) {
		t.Error("Names not equal")
	}
}

// Negative test, ProcAddWebhook
func Test_Neg_ProcAddWebhook(t *testing.T) {
	recorder := httptest.NewRecorder()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAddWebHook(r, recorder), recorder)
	}))
	defer server.Close()

	byteSlice, _ := json.Marshal(testWebhookNeg)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)

	req, err := http.NewRequest("POST", server.URL+"/exchange", reader)
	if err != nil {
		t.Error("Failed to construct new request to server! | Error: " + err.Error())
	}

	_, err = client.Do(req)
	if err != nil {
		t.Error("Failed to POST to  \"" + server.URL + "\" | Error: " + err.Error())
	}

	if recorder.Result().StatusCode != http.StatusNotAcceptable {
		t.Error("Got wrong error code! Code: " + strconv.Itoa(recorder.Result().StatusCode) + " | Expected: " + strconv.Itoa(http.StatusNotAcceptable))
	}
}

// Positive test, ProcGetWebhook
/*func Test_Pos_ProcGetWebhook(t *testing.T) {
	recorder := httptest.NewRecorder()
	fetchedWebhook := types.WebhookInfo{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procGetWebHook(r.URL.Path, recorder), recorder)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/exchange/" + testWebhookPos.ID.Hex())
	if err != nil {
		t.Error("Failed to from the server! | Error: " + err.Error())
	}

	if json.NewDecoder(resp.Body).Decode(&fetchedWebhook) != nil {
		t.Error("Failed to POST to  \"" + server.URL + "\" | Error: " + err.Error())
	}

	if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("Names not equal")
	}
}
*/
