package main

import (
	"fmt"
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
var DB = &mongodb.MongoDB{
	"mongodb://admin:assign3@ds157185.mlab.com:57185/assignment3",
	"assignment3",
	"webhook",
	"curr"}

func updateWebhooks(isWorking *bool, shouldRun4Ever bool) {
	dayUpdated := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	didUpdate := false
	for true {
		log.Println(logInfo + "Checking if update to currency collection is needed")
		if dayUpdated != time.Now().Format("2006-01-02") && !(time.Now().Hour() < 17) {
			didUpdate = false

			log.Println(logInfo + "--------------- Updating currency ---------------")
			DB.UpdateCurr()
			log.Println(logInfo + "--------------- Invoking webhooks ---------------")
			DB.InvokeWebhooks(true)

			log.Println(logInfo + "Updating and Invoking was successfull")
			dayUpdated = time.Now().Format("2006-01-02")
			didUpdate = true
		}

		// Report time to sleep
		timeLoc, _ := time.LoadLocation("CET")
		if didUpdate {
			(*isWorking) = true
			log.Println(logInfo + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 17, 0, 0, 0, timeLoc)).String())
		} else {
			(*isWorking) = false
			log.Println(logInfo + "TIME TO CURRENCY UPDATE: " + time.Until(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 17, 0, 0, 0, timeLoc)).String())
		}

		if shouldRun4Ever {

			// Sleep for 1 min
			time.Sleep(time.Minute)
		} else {
			log.Println(logInfo + "Was ordered to only run once!")
			return
		}
	}
}

func keepAlive(isWorking *bool, shouldRun4Ever bool) {
	for true {
		log.Println(logInfo + "---------- keepAlive() ---------- ")
		log.Println(logInfo + "Trying to http.Get from web")
		_, err := http.Get("https://sleepy-eyrie-58724.herokuapp.com/")
		if err != nil {
			if isWorking != nil {
				fmt.Println("Dereference! isWorking != nil")
				(*isWorking) = false
			}
			log.Println(logError + "Failed http.Get from web!")
		} else {
			fmt.Println("Dereference! No check!")
			(*isWorking) = true
			log.Println(logInfo + "Successful http.Get to web!")
		}

		log.Println(logInfo + "Sleeping for 15 min!")
		log.Println(logInfo + "---------- keepAlive() END ----------")

		if shouldRun4Ever {
			// Sleep for 15 min
			time.Sleep(time.Minute * 15)
		} else {
			log.Println(logInfo + "Was ordered to only run once!")
			return
		}
	}
}

func main() {
	ch := make(chan bool, 1)
	ch <- true

	// run once for first input, then wait forever
	for range ch {
		var isWorking bool
		go keepAlive(&isWorking, true)
		go updateWebhooks(&isWorking, true)
	}
}
