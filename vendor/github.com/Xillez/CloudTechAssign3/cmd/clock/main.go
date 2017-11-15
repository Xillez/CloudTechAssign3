package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Xillez/CloudTechAssign3/mongodb"
)

var testCurr string

var Warn = "[WARNING]: "
var Error = "[ERROR]: "
var Info = "[INFO]: "

//var db = &mongodb.MongoDB{"mongodb://localhost", "Currencies", "webhook", "curr"}
var DB = &mongodb.MongoDB{
	"mongodb://admin:assign3@ds157185.mlab.com:57185/assignment3",
	"assignment3",
	"webhook",
	"curr"}

func updateWebhooks(isWorking *bool, shouldRun4Ever bool) {
	dayUpdated := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	didUpdate := false
	for true {
		if dayUpdated != time.Now().Format("2006-01-02") && !(time.Now().Hour() < 17) {
			didUpdate = false
			log.Println(Info + "--------------- Updating currency ---------------")
			DB.UpdateCurr()
			log.Println(Info + "--------------- Invoking webhooks ---------------")
			DB.InvokeWebhooks(true)

			dayUpdated = time.Now().Format("2006-01-02")
			didUpdate = true
		}

		// Report time to sleep
		timeLoc, _ := time.LoadLocation("CET")
		if didUpdate {
			*isWorking = true
			log.Println(Info + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 0, 0, 0, timeLoc)).String())
		} else {
			*isWorking = false
			log.Println(Info + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, timeLoc)).String())
		}

		if shouldRun4Ever {
			// Sleep for 1 min
			time.Sleep(time.Minute)
		}else {
			return
		}
	}
}

func keepAlive(isWorking *bool, shouldRun4Ever bool) {
	for true {
		log.Println(Info + "---------- keepAlive() ---------- ")
		log.Println(Info + "Trying to http.Get from web")
		_, err := http.Get("https://powerful-brushlands-93106.herokuapp.com/")
		if err != nil {
			if isWorking != nil {*isWorking = false}
			log.Println(Error + "Failed http.Get from web!")
		} else {
			*isWorking = true
			log.Println(Info + "Successful http.Get to web!")
		}

		log.Println(Info + "Sleeping for 15 min!")
		log.Println(Info + "---------- keepAlive() END ----------")

		if shouldRun4Ever {
			// Sleep for 15 min
			time.Sleep(time.Minute * 15)
		}else {
			return
		}
	}
}

func main() {
	ch := make(chan bool, 1)
	ch <- true

	// run once for first input, then wait forever
	for range ch{
		var isWorking *bool
		go keepAlive(isWorking, true)
		go updateWebhooks(isWorking, true)
	}
}