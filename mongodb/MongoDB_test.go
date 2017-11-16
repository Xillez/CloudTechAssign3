package mongodb

import (
	"strconv"
	"testing"

	"github.com/Xillez/CloudTechAssign3/types"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

var testdb = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
var testWebhookPos = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "EUR", "NOK", 0.5, 2.0}
var testCurrPos = types.CurrencyInfo{bson.NewObjectId(), "EUR", "2017-01-01", map[string]float64{"NOK": 2.0}}

// Positive test, db.AddWebhook
func Test_Pos_AddWebhook(t *testing.T) {
	// Initialize database
	testdb.Init()

	fetchedWebhook := types.WebhookInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add webhook
	errAdd := testdb.AddWebhook(testWebhookPos)
	if errAdd.Status != 0 {
		t.Error("Add function got an error: " + strconv.Itoa(errAdd.Status) + " | " + errAdd.Msg)
	}

	// Fetch the added webhook
	errFind := session.DB(testdb.DatabaseName).C(testdb.WebCollName).Find(bson.M{"_id": testWebhookPos.ID}).One(&fetchedWebhook)
	if errFind != nil {
		t.Error("Failed to fecth added element in the database")
	}

	// Check if the URL is the same
	if fetchedWebhook.URL != testWebhookPos.URL {
		t.Error("Added webhook's URL isn't the same!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()
}

// Positive test, db.AddCurr
func Test_Pos_AddCurr(t *testing.T) {
	// Initialize database
	testdb.Init()

	fetchedCurr := types.CurrencyInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test currency
	errAdd := testdb.AddCurr(testCurrPos)
	if errAdd.Status != 0 {
		t.Error("Add function got an error: " + strconv.Itoa(errAdd.Status) + " | " + errAdd.Msg)
	}

	// Fetch the added currency
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Find(bson.M{"_id": testCurrPos.ID}).One(&fetchedCurr)
	if errFind != nil {
		t.Error("Failed to fecth added element in the database")
	}

	// Check if the ID's are the same
	if fetchedCurr.ID != testCurrPos.ID {
		t.Error("Added currency's id isn't the same!")
	}
	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
}

// Positive test, db.GetWebhook
func Test_Pos_GetWebhook(t *testing.T) {
	fetchedWebhook := types.WebhookInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.WebCollName).Insert(testWebhookPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Fetch the added webhook
	errGet := testdb.GetWebhook(testWebhookPos.ID.Hex(), &fetchedWebhook)
	if errGet.Status != 0 {
		t.Error("GetWebhook function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	// Check to see if the ID's are the same
	if fetchedWebhook.ID != testWebhookPos.ID {
		t.Error("Fetched webhook's id isn't the same!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()
}

// Positive test, db.GetCurrency
func Test_Pos_GetCurrency(t *testing.T) {
	fetchedCurrency := types.CurrencyInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Fetch the added webhook
	errGet := testdb.GetCurrByDate(testCurrPos.Date, &fetchedCurrency)
	if errGet.Status != 0 {
		t.Error("GetCurrByDate function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	// Check to see if the ID's are the same
	if fetchedCurrency.ID != testCurrPos.ID {
		t.Error("Fetched currencies's id isn't the same!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
}

// Positive test, db.GetAllWebhooks
func Test_Pos_GetAllWebhooks(t *testing.T) {
	fetchedWebhook := []types.WebhookInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.WebCollName).Insert(testWebhookPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Fetch the added webhook
	errGet := testdb.GetAllWebhooks(&fetchedWebhook)
	if errGet.Status != 0 {
		t.Error("GetWebhook function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	count := len(fetchedWebhook)

	if count != 1 {
		t.Error("Failed to get expected nr of elements: " + strconv.Itoa(count) + " | Expected 1!")
	}

	// Check to see if the ID's are the same
	if fetchedWebhook[0].ID != testWebhookPos.ID {
		t.Error("Fetched webhook's id isn't the same!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()
}

// Positive test, db.GetAllCurr
func Test_Pos_GetAllCurr(t *testing.T) {
	fetchedCurrency := []types.CurrencyInfo{}
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Fetch the added webhook
	errGet := testdb.GetAllCurr(&fetchedCurrency)
	if errGet.Status != 0 {
		t.Error("GetCurrByDate function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	count := len(fetchedCurrency)

	if count != 1 {
		t.Error("Failed to get expected nr of elements: " + strconv.Itoa(count) + "| Expected 1!")
	}

	// Check to see if the ID's are the same
	if fetchedCurrency[0].ID != testCurrPos.ID {
		t.Error("Fetched currencies's id isn't the same!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
}

// Positive test, db.Count
func Test_Pos_Count(t *testing.T) {
	// Setup session with database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Fetch the added webhook
	count, errGet := testdb.Count(testdb.CurrCollName)
	if errGet.Status != 0 {
		t.Error("GetCurrByDate function got an error: " + strconv.Itoa(errGet.Status) + " | " + errGet.Msg)
	}

	count, errCount := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Count()
	if errCount != nil {
		t.Error("Faild to count elements in table: " + testdb.CurrCollName + " in the database!")
	}

	if count != 1 {
		t.Error("Failed to get expected nr of elements: " + strconv.Itoa(count) + "| Expected 1!")
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
}
