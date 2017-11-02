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

	dayUpdated := time.Now().Format("2006-01-02")
	for true /*dayUpdated != time.Now().Format("2006-01-02")*/ {
		if dayUpdated != time.Now().Format("2006-01-02") {
			log.Println(Info + "--------------- Updating currency ---------------")
			db.UpdateCurr()
			log.Println(Info + "--------------- Invoking webhooks ---------------")
			db.InvokeWebhooks(true)
		}

		log.Println(Info + "--------------- Sleeping ---------------")
		timeLoc, _ := time.LoadLocation("CET")

		//fmt.Println("-------- TIME UNTIL TOMORROW 17:00!: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 00, 00, 00, timeLoc)).String())

		// Sleep until 17:00 the next day
		time.Sleep(time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 00, 00, 00, timeLoc)))
	}
}
