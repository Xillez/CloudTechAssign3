package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Xillez/CloudTechAssign3/mongodb"
)

var testCurr string

const logWarn = "[WARNING]: "
const logError = "[ERROR]: "
const logInfo = "[INFO]: "

//var db = &mongodb.MongoDB{"mongodb://localhost", "Currencies", "webhook", "curr"}
var db = &mongodb.MongoDB{"mongodb://admin:assign3@ds157185.mlab.com:57185/assignment3", "assignment3", "webhook", "curr"}

func keepAlive() {
	for true {
		log.Println(logInfo + "---------- keepAlive() ---------- ")
		log.Println(logInfo + "Trying to http.Get from web")
		_, err := http.Get("https://powerful-brushlands-93106.herokuapp.com/")
		if err != nil {
			log.Println(logError + "Failed http.Get from web!")
		} else {
			log.Println(logInfo + "Successful http.Get to web!")
		}

		log.Println(logInfo + "Sleeping for 15 min!")
		log.Println(logInfo + "---------- keepAlive() END ----------")
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
			log.Println(logInfo + "--------------- Updating currency ---------------")
			db.UpdateCurr()
			log.Println(logInfo + "--------------- Invoking webhooks ---------------")
			db.InvokeWebhooks(true)

			dayUpdated = time.Now().Format("2006-01-02")
			didUpdate = true
		}

		// Report time to sleep
		timeLoc, _ := time.LoadLocation("CET")
		if didUpdate {
			log.Println(logInfo + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 0, 0, 0, timeLoc)).String())
		} else {
			log.Println(logInfo + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, timeLoc)).String())
		}

		// Sleep for 1 min
		time.Sleep(time.Minute)
	}
}
