package mongodb

import (
	"strconv"
	"testing"

	"github.com/Xillez/CloudTechAssign3/types"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"time"
	"net/http"
)

var testdb = &MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
var testWebhookPos = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "EUR", "NOK", 0.5, 2.0}
var testCurrPos = types.CurrencyInfo{bson.NewObjectId(), "EUR", "2017-01-01", map[string]float64{"NOK": 2.0}}

// Positive test, db.AddWebhook
func Test_Pos_AddWebhook(t *testing.T) {
	// Initialize database
	testdb.Init()

	fetchedWebhook := types.WebhookInfo{}
	// Dial database
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
	// Dial database
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
	// Dial database
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

// Positive test, db.GetAllWebhooks
func Test_Pos_GetAllWebhooks(t *testing.T) {
	var fetchedWebhook []types.WebhookInfo
	// Dial database
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

// Positive test, db.GetCurrency
func Test_Pos_GetCurrByDate(t *testing.T) {
	fetchedCurrency := types.CurrencyInfo{}
	// Dial database
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

// Positive test, db.GetAllCurr
func Test_Pos_GetAllCurr(t *testing.T) {
	var fetchedCurrency []types.CurrencyInfo
	// Dial database
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
	// Dial database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test currency manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// Counts the added webhooks
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

// Positive test, db.Count
func Test_Neg_Count(t *testing.T) {
	// Fetch the added webhook
	_, errGet := testdb.Count("")
	if errGet.Status != http.StatusInternalServerError {
		t.Error("Count did not return error")
	}
}

// Positive test, db.UpdateCurr()
func Test_Pos_UpdateCurr(t *testing.T){
	// Dial database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test currency manually
	errInsert := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errInsert != nil {
		t.Error("Failed to add element in the database for testing")
	}

	var errUpdate = testdb.UpdateCurr()
	if errUpdate.Status != 0 {
		t.Error(errUpdate.Msg)
	}

	targetDate := time.Now().Format("2006-01-02")

	q := bson.M{"date":targetDate}
	result := types.CurrencyInfo{}

	//Find test currency manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Find(q).One(&result)
	if errFind != nil {
		t.Error("Failed to find element in the database for testing")
	}

	if result.Date != targetDate {
		endStr := testCurrPos.Date + " vs " + result.Date
		t.Error("Database not updated, dates not the same: " + endStr)
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
}

// Positive test, db.DelWebhook
func Test_Pos_DelWebook(t *testing.T){
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errInsert := session.DB(testdb.DatabaseName).C(testdb.WebCollName).Insert(testWebhookPos)
	if errInsert != nil {
		t.Error("Failed to add webhook in the database for testing")
	}

	var errUpdate = testdb.DelWebhook(testWebhookPos.ID.Hex())
	if errUpdate.Status != 0 {
		t.Error(errUpdate.Msg)
	}

	q := bson.M{"_id":testWebhookPos.ID}
	result := types.CurrencyInfo{}

	//Find test currency manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Find(q).One(&result)
	if errFind == nil {
		t.Error("Failed to find webhook in the database for testing")
	}

	if result.ID.String() == testWebhookPos.ID.String() {
		t.Error("Database did not delete object, found ObjectID: " + result.ID.String())
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()
}

func Test_InvokeWebhooks(t *testing.T){
	// Dial database
	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Add test webhook manually
	errInsert := session.DB(testdb.DatabaseName).C(testdb.WebCollName).Insert(testWebhookPos)
	if errInsert != nil {
		t.Error("Failed to add webhook in the database for testing")
	}

	// Add test currency manually
	errFind := session.DB(testdb.DatabaseName).C(testdb.CurrCollName).Insert(testCurrPos)
	if errFind != nil {
		t.Error("Failed to add element in the database for testing")
	}

	// invokes all webhooks
	custErr := testdb.InvokeWebhooks(false)
	if custErr.Status != 0 {
		t.Error(custErr.Msg)
	}

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.CurrCollName).DropCollection()
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()
}