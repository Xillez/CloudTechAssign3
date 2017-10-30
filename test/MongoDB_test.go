package main

import (
	"net/http"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var dbTest *MongoDB

// Positive test, MongoDB.AddWebhook
func Test_Pos_AddWebhook(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	testWebhook := WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", bson.NewObjectId(), 0.5, 2.0}
	fetchedWebhook := WebhookInfo{}
	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Add webhook
	errAdd := dbTest.AddWebhook(testWebhook)
	if errAdd.Status != 0 {
		t.Error("Add function got an error: " + strconv.Itoa(errAdd.Status) + " | " + errAdd.Msg)
	}

	// Fetch the added webhook
	errFind := session.DB(dbTest.DatabaseName).C(dbTest.WebCollName).Find(bson.M{"_id": testWebhook.ID}).One(&fetchedWebhook)
	if errFind != nil {
		t.Error("Failed to fecth added element in the database")
	}

	// Check if the URL is the same
	if fetchedWebhook.URL != testWebhook.URL {
		t.Error("Added webhook's URL isn't the same!")
	}
}

// Positive test, MongoDB.AddCurr
func Test_Pos_AddCurr(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	testCurr := CurrencyInfo{bson.NewObjectId(), "EUR", "NOK", 2.0}
	fetchedCurr := CurrencyInfo{}
	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Add test currency
	errAdd := dbTest.AddCurr(testCurr)
	if errAdd.Status != 0 {
		t.Error("Add function got an error: " + strconv.Itoa(errAdd.Status) + " | " + errAdd.Msg)
	}

	// Fetch the added currency
	errFind := session.DB(dbTest.DatabaseName).C(dbTest.CurrCollName).Find(bson.M{"_id": testCurr.ID}).One(&fetchedCurr)
	if errFind != nil {
		t.Error("Failed to fecth added element in the database")
	}

	// Check if the ID's are the same
	if fetchedCurr.ID != testCurr.ID {
		t.Error("Added currency's id isn't the same!")
	}
}

// Negative test, MongoDb.AddWebhook
func Test_Neg_AddWebhook(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Invalid webhook, "BaseCurrency" is more than 3 letters long!
	testWebhook := WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", bson.NewObjectId(), 0.5, 2.0}

	// Try and add it
	errAdd := dbTest.AddWebhook(testWebhook)

	// Make sure the addition of invalid webhook fails
	if errAdd.Status == http.StatusNotAcceptable {
		t.Error("Added an invalid webhook entry! Status: " + strconv.Itoa(errAdd.Status) + " | \"" + http.StatusText(errAdd.Status) + "\"")
	}
}

// Negative test, MongoDb.AddCurr
func Test_Neg_AddCurr(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Invalid currency, "BaseCurrency" is more than 3 letters long!
	testCurr := CurrencyInfo{bson.NewObjectId(), "EUR", "NOK", 1.0}

	// Add test currency
	errAdd := dbTest.AddCurr(testCurr)

	// Make sure addition of invalid currency fails
	if errAdd.Status == http.StatusNotAcceptable {
		t.Error("Added an invalid currency entry! Status: " + strconv.Itoa(errAdd.Status) + " | \"" + http.StatusText(errAdd.Status) + "\"")
	}
}

// Positive test, MongoDB.GetWebhook
func Test_Pos_GetWebhook(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	testWebhookID := bson.NewObjectId()
	testWebhook := WebhookInfo{testWebhookID, "SomeAddress", bson.NewObjectId(), 2.0, 2.5}
	fetchedWebhook := WebhookInfo{}
	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(dbTest.DatabaseName).C(dbTest.WebCollName).Insert(testWebhook)
	if errFind != nil {
		t.Error("Failed to add element in the database")
	}

	// Fetch the added webhook
	errGet := dbTest.GetWebhook(testWebhookID.String(), &fetchedWebhook)
	if errGet.Status != 0 {
		t.Error("GetWebhook function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	// Check to see if the ID's are the same
	if fetchedWebhook.ID != testWebhookID {
		t.Error("Fetched webhook's id isn't the same!")
	}
}

// Negative test, MongoDb.GetWebhook
func Test_Neg_GetWebhook(t *testing.T) {
	dbTest = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}

	testWebhookID := bson.NewObjectId()
	testWebhook := WebhookInfo{testWebhookID, "SomeAddress", bson.NewObjectId(), 2.0, 2.5}
	fetchedWebhook := WebhookInfo{}
	// Dial database
	session, err := mgo.Dial(dbTest.DatabaseURL)

	// Check for session error
	if err != nil {
		panic(err)
	}

	// Postpone closing connection until we return
	defer session.Close()

	// Add test webhook
	errFind := session.DB(dbTest.DatabaseName).C(dbTest.WebCollName).Insert(testWebhook)
	if errFind != nil {
		t.Error("Failed to add element in the database")
	}

	// Fetch the added webhook
	errGet := dbTest.GetWebhook(testWebhookID.String(), &fetchedWebhook)
	if errGet.Status == http.StatusNotAcceptable {
		t.Error("Fetched an invalid webhook entry! Status: " + strconv.Itoa(errGet.Status) + " | \"" + http.StatusText(errGet.Status) + "\"")
	}
}
