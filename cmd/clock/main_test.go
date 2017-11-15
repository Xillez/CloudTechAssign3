package main

import (
	"testing"
	"github.com/Xillez/CloudTechAssign3/mongodb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/Xillez/CloudTechAssign3/types"
)

var testdb = &mongodb.MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
var testWebhookPos = types.WebhookInfo{bson.NewObjectId(), "http://webhook:8080/something", "EUR", "NOK", 0.5, 2.0}
var testCurrPos = types.CurrencyInfo{bson.NewObjectId(), "EUR", "2017-01-01", map[string]float64{"NOK": 2.0}}


func Test_TestingSetup(t *testing.T) {
	DB = &mongodb.MongoDB{"mongodb://localhost", "Testing", "testWeb", "testCurr"}
}

func Test_keepAlive(t *testing.T){
	isWorking := false
	keepAlive(&isWorking, false)
	if !isWorking{
		t.Error("Webhooks were not updated")
	}

	session, err := mgo.Dial(testdb.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Clean up after testing
	_ = session.DB(testdb.DatabaseName).C(testdb.WebCollName).DropCollection()

}

func Test_updateWebhooks(t *testing.T){
	isWorking := false
	updateWebhooks(&isWorking, false)
	if !isWorking{
		t.Error("Webhooks were not updated")
	}
}