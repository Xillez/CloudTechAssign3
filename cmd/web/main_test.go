package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

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
var testCurrReqPos = types.CurrencyReq{"EUR", "NOK"}
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
	_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
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
	errClean := session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
	if errClean != nil {
		fmt.Println(errClean.Error())
	}
}

// Positive test, ProcDelWebhook
func Test_Pos_ProcDelWebhook(t *testing.T) {
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

	// Try inserting webhook for testing
	errInsert := session.DB(DB.DatabaseName).C(DB.WebCollName).Insert(&testWebhookPos)
	if errInsert != nil {
		t.Error("Failed inserting test webhook! | Error: " + errInsert.Error())
	}

	// Make server for testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procDelWebHook(r.URL.Path, w), w)
	}))
	defer server.Close()

	// Make the request
	//fmt.Println(server.URL + "/exchange/" + testWebhookPos.ID.Hex())

	req, err := http.NewRequest("DELETE", server.URL+"/exchange/"+testWebhookPos.ID.Hex(), nil)
	if err != nil {
		t.Error("Failed to construct new request to server! | Error: " + err.Error())
	}

	// Set content as applicaiton/json
	req.Header.Add("content-type", "application/json")

	// Setup POST request
	_, err = client.Do(req)
	if err != nil {
		t.Error("Failed to POST to  \"" + server.URL + "\" | Error: " + err.Error())
	}

	// Trying to count nr webhooks, expecting 0
	nr, errCount := session.DB(DB.DatabaseName).C(DB.WebCollName).Count()
	if errCount != nil {
		t.Error("Failed counting from database! | Error: " + errCount.Error())
	}

	// Checking to see if database is empty
	if nr != 0 {
		t.Error("Webhook wasn't deleted! | Nr: " + strconv.Itoa(nr))
	}

	// Clean up after testing, if needed
	_ = session.DB(DB.DatabaseName).C(DB.WebCollName).DropCollection()
}

// Positive test, ProcAverage
func Test_Pos_ProcAverage(t *testing.T) {
	// Setup session with database
	session, err := mgo.Dial(DB.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Make sure currency collection is empty
	_ = session.DB(DB.DatabaseName).C(DB.CurrCollName).DropCollection()

	// Initialize database after droping it
	DB.Init()

	// Save old date and ID for later
	/*testCurrPos_OldDate := testCurrPos.Date
	testCurrPos_OldID := testCurrPos.ID*/

	for i := 0; i < 3; i++ {
		// Set date to the current day - i
		testCurrPos.Date = time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if time.Now().Hour() < 17 {
			testCurrPos.Date = time.Now().AddDate(0, 0, -(i + 1)).Format("2006-01-02")
		}
		// Set a new ID for each iteration
		testCurrPos.ID = bson.NewObjectId()

		fmt.Printf("%v\n", testCurrPos)

		// Add testCurrency to database
		err = session.DB(DB.DatabaseName).C(DB.WebCollName).Insert(&testCurrPos)
		if err != nil {
			t.Error("Failed adding test currency nr:  " + strconv.Itoa(i) + " with date " + testCurrPos.Date + " to database")
		}
	}

	// Make server for testing
	/*server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.CheckPrintErr(procAverage(r, w), w)
	}))
	defer server.Close()

	// Convert to byteReader
	byteSlice, _ := json.Marshal(testCurrReqPos)
	byteReader := bytes.NewReader(byteSlice)
	reader := ioutil.NopCloser(byteReader)
	//fmt.Println(string(byteSlice))
	//l, _ := ioutil.ReadAll(reader)
	//fmt.Println(string(l))

	// Make the request
	req, err := http.NewRequest("POST", server.URL+"/exchange/average", reader)
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

	slice, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(slice))
	avg, errParseFloat := strconv.ParseFloat(string(slice), 64)
	if err != nil {
		t.Error("Failed to ReadAll response from \"" + server.URL + "\" | Error: " + err.Error())
	}
	if errParseFloat != nil {
		t.Error("Failed to ParseFloat within response from \"" + server.URL + "\" | Error: " + errParseFloat.Error())
	}

	if avg != testCurrPos.Rates["NOK"] {
		t.Error("Average of NOK over 3 days with similar rate each day isn't what it's supposed to be! | Expected : " + strconv.FormatFloat(testCurrPos.Rates["NOK"], 'f', -1, 64) + " | Got: " + strconv.FormatFloat(avg, 'f', -1, 64))
	}

	// Clean up testCurrPos.Date and ID
	testCurrPos.Date = testCurrPos_OldDate
	testCurrPos.ID = testCurrPos_OldID*/

	// Clean up after testing
	_ = session.DB(DB.DatabaseName).C(DB.CurrCollName).DropCollection()
}
