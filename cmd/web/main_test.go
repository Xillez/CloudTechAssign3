package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Xillez/CloudTechAssign3/mongodb"
	"github.com/Xillez/CloudTechAssign3/types"
	"github.com/Xillez/CloudTechAssign3/utils"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//var testdb = &mongodb.MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
var testWebhookPos = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "EUR", "NOK", 0.5, 2.0}
var testWebhookNeg = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "CAD", "NOK", 0.5, 2.0}
var testCurrPos = types.CurrencyInfo{bson.NewObjectId(), "EUR", "2017-01-01", map[string]float64{"NOK": 2.0}}
var client = http.Client{}

func Test_TestingSetup(t *testing.T) {
	DB = &mongodb.MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
}

// Positive test, ProcGetWebhook
func Test_Pos_ProcGetWebhook(t *testing.T) {
	// Initialize database
	DB.Init()

	fetchedWebhook := types.WebhookDisp{}

	// Dial database
	session, err := mgo.Dial(DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add testwebhok to database
	err = session.DB(DB.DatabaseName).C(DB.WebCollName).Insert(&testWebhookPos)
	if err != nil {
		t.Error("Failed adding test webhook to database")
	}

	// Make server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procGetWebHook(r.URL.Path, w), w)
	}))
	defer server.Close()

	// Get webhook from processing function
	resp, err := http.Get(server.URL + "/exchange/" + testWebhookPos.ID.Hex())
	if err != nil {
		t.Error("Failed to from the server! | Error: " + err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&fetchedWebhook)
	if err != nil {
		t.Error("Failed to decode response from \"" + server.URL + "\" | Error: " + err.Error())
	}

	if fetchedWebhook.TargetCurrency != testWebhookPos.TargetCurrency {
		t.Error("Names not equal")
	}

	// Clean up after testing
	_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}

// Positive test, ProcAddWebhook
func Test_Pos_ProcAddWebhook(t *testing.T) {
	// Initialize database
	DB.Init()

	fetchedWebhook := types.WebhookInfo{}

	// Dial database
	session, err := mgo.Dial(DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Make server for testing
	recorder := httptest.NewRecorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAddWebHook(r, recorder), recorder)
	}))
	defer server.Close()

	// Convert to byteReader
	byteSlice, _ := json.Marshal(testWebhookPos)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)

	// Make the request
	req, err := http.NewRequest("POST", server.URL+"/exchange", reader)
	if err != nil {
		t.Error("Failed to construct new request to server! | Error: " + err.Error())
	}

	// Set content as applicaiton/json
	req.Header.Add("content-type", "application/json")

	// Set POST request
	resp, err := client.Do(req)
	if err != nil {
		t.Error("Failed to POST to  \"" + server.URL + "\" | Error: " + err.Error())
	}

	// Check to see if ID is the same
	if strings.Contains(resp., testWebhookPos.ID.Hex()) {
		t.Error("Names not equal")
	}

	// Try to fetch added webhook
	errFind := session.DB(DB.DatabaseName).C(DB.WebCollName).Find(bson.M{"_id": testWebhookPos.ID}).One(&fetchedWebhook)
	if errFind != nil {
		t.Error("Failed to fetch added webhook")
	}

	if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("The webhook's id aren't the same")
	}

	// Clean up after testing
	_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}
