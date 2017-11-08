package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Xillez/CloudTechAssign2/mongodb"
)

var testCurr string

var Warn = "[WARNING]: "
var Error = "[ERROR]: "
var Info = "[INFO]: "

//var db = &mongodb.MongoDB{"mongodb://localhost", "Currencies", "webhook", "curr"}
var db = &mongodb.MongoDB{"mongodb://admin:admin@ds151452.mlab.com:51452/webhook_curr", "webhook_curr", "webhook", "curr"}

func keepAlive() {
	for true {
		log.Println("[INFO]: ---------- keepAlive() ---------- ")
		log.Println("[INFO]: Trying to http.Get from web")
		_, err := http.Get("https://powerful-brushlands-93106.herokuapp.com/")
		if err != nil {
			log.Println("[ERROR]: Failed http.Get from web!")
		} else {
			log.Println("[INFO]: Successful http.Get to web!")
		}

		log.Println("[INFO] Sleeping for 15 min!")
		log.Println("[INFO]: ---------- keepAlive() END ----------")
		time.Sleep(time.Minute * 15)
	}
}

func main() {
	go keepAlive()

	dayUpdated := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	didUpdate := false
	for true {
		if dayUpdated != time.Now().Format("2006-01-02") && !(time.Now().Hour() < 17) {
			didUpdate = false
			log.Println(Info + "--------------- Updating currency ---------------")
			db.UpdateCurr()
			log.Println(Info + "--------------- Invoking webhooks ---------------")
			db.InvokeWebhooks(true)

			dayUpdated = time.Now().Format("2006-01-02")
			didUpdate = true
		}

		// Report time to sleep
		timeLoc, _ := time.LoadLocation("CET")
		if didUpdate {
			log.Println("[INFO] TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 0, 0, 0, timeLoc)).String())
		} else {
			log.Println("[INFO] TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, timeLoc)).String())
		}

		// Sleep for 1 min
		time.Sleep(time.Minute)
	}
}
