package main

import (
	"Assignment2/mongodb"
	"log"
	"time"
)

var testCurr string

var Warn = "[WARNING]: "
var Error = "[ERROR]: "
var Info = "[INFO]: "

//var db = &mongodb.MongoDB{"mongodb://localhost", "Currencies", "webhook", "curr"}
var db = &mongodb.MongoDB{"mongodb://admin:admin@ds151452.mlab.com:51452/webhook_curr", "webhook_curr", "webhook", "curr"}

func main() {

	dayUpdated := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	didUpdate := false
	durToSleep := time.Duration(0)
	for true {
		if dayUpdated != time.Now().Format("2006-01-02") && !(time.Now().Hour() < 17) {
			log.Println(Info + "--------------- Updating currency ---------------")
			db.UpdateCurr()
			log.Println(Info + "--------------- Invoking webhooks ---------------")
			db.InvokeWebhooks(true)

			dayUpdated = time.Now().Format("2006-01-02")
			didUpdate = true
		}

		log.Println(Info + "--------------- Sleeping ---------------")
		timeLoc, _ := time.LoadLocation("CET")
		if didUpdate {
			durToSleep = time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 00, 00, 00, timeLoc))
		} else {
			durToSleep = time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 00, 00, 00, timeLoc))
		}
		log.Println("[INFO] SLEEP DURATION: " + durToSleep.String())

		// Sleep until 17:00 the next day
		time.Sleep(durToSleep)
	}
}
