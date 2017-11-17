package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Xillez/CloudTechAssign3/mongodb"
	"github.com/Xillez/CloudTechAssign3/types"
	"github.com/Xillez/CloudTechAssign3/utils"

	"gopkg.in/mgo.v2"
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

	// Setup session with database
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
		t.Error("Failed to fetch from the server! | Error: " + err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&fetchedWebhook)
	if err != nil {
		t.Error("Failed to decode response from \"" + server.URL + "\" | Error: " + err.Error())
	}

	if fetchedWebhook.TargetCurrency != testWebhookPos.TargetCurrency {
		t.Error("Names not equal")
	}

	// Clean up after testing
	//_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}

// Positive test, ProcAddWebhook
func Test_Pos_ProcAddWebhook(t *testing.T) {
	// Initialize database
	DB.Init()

	fetchedWebhook := types.WebhookInfo{}

	// Setup session with database
	session, err := mgo.Dial(DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Make server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAddWebHook(r, w), w)
	}))
	defer server.Close()

	// Convert to byteReader
	byteSlice, _ := json.Marshal(testWebhookPos)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)
	//fmt.Println(string(byteSlice))
	//l, _ := ioutil.ReadAll(reader)
	//fmt.Println(string(l))

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

	slice, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(slice))

	/*if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("The webhook's id aren't the same")
	}*/

	// Try to fetch added webhook
	errFind := session.DB(DB.DatabaseName).C(DB.WebCollName).Find(bson.M{"_id": bson.ObjectIdHex(string(slice))}).One(&fetchedWebhook)
	if errFind != nil {
		t.Error("Failed to fetch added webhook | FindError: " + errFind.Error())
	}

	//fmt.Printf("%v\n", testWebhookPos)
	//fmt.Printf("%v\n", fetchedWebhook)

	//fmt.Printf("%v\n", bson.ObjectIdHex(string(slice)))

	if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("The webhook's id aren't the same")
	}

	// Clean up after testing
	//_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}

// Positive test, ProcDelWebhook
/*func Test_Pos_ProcDelWebhook(t *testing.T) {
	// Setup session with database
	session, err := mgo.Dial(DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Make sure database is empty
	_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()

	// Initialize database
	DB.Init()

	// Make server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAddWebHook(r, w), w)
	}))
	defer server.Close()

	// Convert to byteReader
	byteSlice, _ := json.Marshal(testWebhookPos)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)
	//fmt.Println(string(byteSlice))
	//l, _ := ioutil.ReadAll(reader)
	//fmt.Println(string(l))

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

	slice, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(slice))

	/*if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("The webhook's id aren't the same")
	}*

	// Try to fetch added webhook
	errFind := session.DB(DB.DatabaseName).C(DB.WebCollName).Find(bson.M{"_id": bson.ObjectIdHex(string(slice))}).One(&fetchedWebhook)
	if errFind != nil {
		t.Error("Failed to fetch added webhook | FindError: " + errFind.Error())
	}

	//fmt.Printf("%v\n", testWebhookPos)
	//fmt.Printf("%v\n", fetchedWebhook)

	//fmt.Printf("%v\n", bson.ObjectIdHex(string(slice)))

	if fetchedWebhook.ID.Hex() != testWebhookPos.ID.Hex() {
		t.Error("The webhook's id aren't the same")
	}

	// Clean up after testing
	//_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}*/
