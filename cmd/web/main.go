// Program is made in collaboration with / got help from
//          - Jonas J. Solsvik
// 			- Zohaib Butt
// 			- Eldar Hauge Torkelsen

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Xillez/CloudTechAssign3/mongodb"
	"github.com/Xillez/CloudTechAssign3/types"
	"github.com/Xillez/CloudTechAssign3/utils"

	"gopkg.in/mgo.v2/bson"
)

const logWarn = "[WARNING]: "
const logError = "[ERROR]: "
const logInfo = "[INFO]: "

var DB = &mongodb.MongoDB{"mongodb://admin:assign3@ds157185.mlab.com:57185/assignment3", "assignment3", "webhook", "curr"}

/* ---------- Root Handler functions ---------- */

func procGetWebHook(url string, w http.ResponseWriter) utils.CustError {
	log.Println(logInfo + "--------------- Got GetWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "application/json")
	webhook := types.WebhookInfo{}

	log.Println(logInfo + "Splitting incomming URL")
	// Split url into usable components
	splitURL, err := utils.GetSplitURL(url, 3)
	if err.Status != 0 {
		return err
	}

	log.Println(logInfo + "Validating GetWebhook request")
	// Validate user input just in case of attack
	if bson.IsObjectIdHex(splitURL[2]) {
		// Somehting went wrong! Inform the user!
		log.Println(logWarn + "Validating failed! Informing user! Id not acceptable: " + splitURL[2])
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": Given webhook id is invalid!"}
	}

	log.Println(logInfo + "Trying to get requested webhook")
	// Retrieve webhook from database
	err = DB.GetWebhook(splitURL[2], &webhook)
	if err.Status != 0 {
		return err
	}

	// Convert WebhookInfo over to WebhookDisp format
	display := types.WebhookDisp{
		webhook.URL,
		webhook.BaseCurrency,
		webhook.TargetCurrency,
		webhook.MinValue,
		webhook.MaxValue,
	}

	log.Println(logInfo + "Encoding webhook to browser")
	// Dump the fetched webhook to browser
	errEncode := json.NewEncoder(w).Encode(&display)
	if errEncode != nil {
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Encoding failed!"}
	}

	log.Println(logInfo + "GetWebhook request finnished successfully")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

func procAddWebHook(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(logInfo + "--------------- Got AddWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	webhook := types.WebhookInfo{}

	// Decode request
	log.Println(logInfo + "Decoding AddWebhook request")
	errDec := json.NewDecoder(r.Body).Decode(&webhook)
	if errDec != nil {
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Webhook registration decoding failed!"}
	}

	log.Println(logInfo + "Validating data from recieved AddWebhook request")
	// Validate vital data
	if webhook.URL == "" {
		log.Println(logWarn + "Validation failed! Informing user! URL is not present!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"webhookURL\" can't be empty!"}
	}
	if webhook.BaseCurrency != "EUR" {
		log.Println(logWarn + "Validation failed! Informing user! \"baseCurrency\" not \"EUR\"")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(webhook.TargetCurrency) != 3 {
		log.Println(logWarn + "Validating failed! Informing user! \"targetCurrency\" has to be 3 long!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}
	if webhook.MinValue < 0 && webhook.MinValue > webhook.MaxValue {
		log.Println(logWarn + "Validating failed! Informing user! \"minTriggerValue\" not between 0 - \"maxTriggerValue\"")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"minTriggerValue\" need to be between 0 - \"maxTriggerValue\""}
	}
	if webhook.MaxValue < webhook.MinValue {
		log.Println(logWarn + "Validating failed! Informing user! \"maxTriggerValue\" has to be >= \"minTriggerValue\"")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"maxTriggerValue\" has to be >= \"minTriggerValue\""}
	}

	log.Println(logInfo + "converting retrieved webhook data into storage format")

	// Set new id if no id is given, if given, use it
	webhookID := bson.NewObjectId()
	if bson.IsObjectIdHex(webhook.ID.Hex()) {
		webhookID = webhook.ID
	}
	// Add webhook info to database
	err := DB.AddWebhook(types.WebhookInfo{
		webhookID,
		webhook.URL,
		webhook.BaseCurrency,
		webhook.TargetCurrency,
		webhook.MinValue,
		webhook.MaxValue})
	if err.Status != 0 {
		return err
	}

	log.Println(logInfo + "Writing stored entry id to browser")
	fmt.Fprintf(w, webhookID.Hex())

	log.Println(logInfo + "AddWebhook request finished successfully!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

func procDelWebHook(url string, w http.ResponseWriter) utils.CustError {
	log.Println(logInfo + "--------------- Got DelWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")

	// Split url to extract id string
	log.Println(logInfo + "Splitting incoming URL")
	splitURL, err := utils.GetSplitURL(url, 3)
	if err.Status != 0 {
		return err
	}

	// Validate vital data
	log.Println(logInfo + "Validating DelWebhook request data")
	if bson.IsObjectIdHex(splitURL[2]) {
		log.Println(logWarn + "Validation failed! Given webhook id is invalid!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": Given webhook id is invalid!"}
	}

	log.Println(logInfo + "Deleting requested webhook")
	// Try to delete webhook entry with given id
	errDel := DB.DelWebhook(splitURL[2])
	if errDel.Status != 0 {
		return errDel
	}

	log.Println(logInfo + "Informing user that the webhook was deleted")
	fmt.Fprintf(w, "Webhook deleted!")

	log.Println(logInfo + "DelWebhook request finished successfully")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

/* ---------- Latest Handler functions ---------- */

func procLatest(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(logInfo + "--------------- Got Latest Currency Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	currReq := types.CurrencyReq{}
	currInfo := types.CurrencyInfo{}

	// Decode request
	log.Println(logInfo + "Decoding latest request body")
	errDec := json.NewDecoder(r.Body).Decode(&currReq)
	if errDec != nil {
		log.Println(logInfo + "Failed decoding latest request body")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Latest currency exchange request decoding failed!"}
	}

	// Validate vital data
	log.Println(logInfo + "Validating data from latest request")
	if currReq.BaseCurrency != "EUR" {
		log.Println(logWarn + "Validation failed! \"baseCurrency\" has to be \"EUR\"!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(currReq.TargetCurrency) != 3 {
		// Somehting went wrong! Inform the user!
		log.Println(logWarn + "Validation failed! \"targetCurrency\" has to be 3 long!")
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	log.Println(logInfo + "Getting currency from database")
	// Retrieve currency info with current date, if error, try with yesterday's dates
	errRetCurr := DB.GetCurrByDate(time.Now().Format("2006-01-02"), &currInfo)
	if errRetCurr.Status != 0 {
		log.Println(logWarn + "Failed getting for current day, trying to get from yester day instead (it's before 17:00)")
		errRetCurr := DB.GetCurrByDate(time.Now().AddDate(0, 0, -1).Format("2006-01-02"), &currInfo)
		if errRetCurr.Status != 0 {
			return errRetCurr
		}
	}

	log.Println(logInfo + "Writing the latest rate of the requested currency")
	fmt.Fprintf(w, strconv.FormatFloat(currInfo.Rates[currReq.TargetCurrency], 'f', -1, 64))

	log.Println(logInfo + "Latest request finished successfully")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}

}

/* ---------- Average Handler functions ---------- */

func fetchAverage(baseCurr string, targetCurr string) (float64, utils.CustError) {
	curr := make([]types.CurrencyInfo, 3, 3)
	avg := 0.0

	for i := 0; i < len(curr); i++ {
		// Find date
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		if time.Now().Hour() < 17 {
			date = time.Now().AddDate(0, 0, -(i + 1)).Format("2006-01-02")
		}

		// Get currency document with specified date
		log.Println(logInfo + "Trying to get currency with date: " + date + " from databse!")
		err := DB.GetCurrByDate(date, &curr[i])
		if err.Status != 0 {
			return 0.0, utils.CustError{http.StatusInternalServerError, utils.ErrorStr[1] + " | Error: Failed to get average!"}
		}

		// Add up given TargetCurrency's rate for averaging
		log.Println(logInfo + "Adding all " + strconv.Itoa(len(curr)) + " day of rates for TargetCurrency")
		avg += curr[i].Rates[targetCurr]
	}

	log.Println(logInfo + "Averaging finished successfully")
	// Nothing bad happened
	return avg / float64(len(curr)), utils.CustError{0, utils.ErrorStr[0]}
}

func procAverage(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(logInfo + "--------------- Got Average Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	curr := types.CurrencyReq{}

	// Decode request
	log.Println(logInfo + "Decode averge request data")
	errDec := json.NewDecoder(r.Body).Decode(&curr)
	if errDec != nil {
		log.Println(logError + "Decoding failed, Informing user!")
		// Somehtin(g went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Failed!"}
	}

	// Validate vital data
	log.Println(logInfo + "Validating data")
	if curr.BaseCurrency != "EUR" {
		log.Println(logWarn + "Validation failed! \"baseCurrency\" has to be \"EUR\"!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(curr.TargetCurrency) != 3 {
		// Somehting went wrong! Inform the user!
		log.Println(logWarn + "Validation failed! \"targetCurrency\" has to be 3 long!")
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	log.Println(logInfo + "Fetching average")
	avg, err := fetchAverage(curr.BaseCurrency, curr.TargetCurrency)
	if err.Status != 0 {
		return err
	}

	log.Println(logInfo + "Writing average to browser")
	fmt.Fprintf(w, strconv.FormatFloat(avg, 'f', -1, 64))

	// Nothing bad happened
	log.Println(logInfo + "Average finished successfully")
	return utils.CustError{0, utils.ErrorStr[0]}
}

// procDialogFlow - Takes JSON structure from the DialogFlow bot and responds
// with the latest exchange rate as the currency-rate variable.
func procDialogFlow(r *http.Request, w http.ResponseWriter) utils.CustError {

	log.Println(logInfo + "--------------- Got DialogFlow Request ---------------")

	jsonReq := types.DialogFlowReq{}

	log.Println(logInfo + "Trying to decode request")
	err := json.NewDecoder(r.Body).Decode(&jsonReq)
	if err != nil {
		log.Printf(logWarn+"Could not decode JSON. Error: %v", err)
		return utils.CustError{http.StatusBadRequest, "Invalid JSON format"}
	}

	log.Printf(logInfo+"Post request for DialogFlow.\n Contains:\n %v.", jsonReq)

	// Checking request
	log.Println(logInfo + "Validating json request")
	if len(jsonReq.Result.Parameters.Currency) != 2 {
		log.Println(logInfo + "Validation failed! \"Currency\" isn't 2 long!")
		return utils.CustError{http.StatusBadRequest, "Expected two Parameters"}
	}

	if jsonReq.Result.Parameters.Currency[0] != "EUR" {
		log.Println(logInfo + "Validation failed! First \"Currency\" isn't \"EUR\"!")
		return utils.CustError{http.StatusNotImplemented, "'EUR' is the only base currency for now"}
	}

	currInfo := types.CurrencyInfo{}

	// Getting latest from database
	log.Println(logInfo + "Trying to get currency by date")
	errRetCurr := DB.GetCurrByDate(time.Now().Format("2006-01-02"), &currInfo)
	if errRetCurr.Status != 0 {

		log.Println(logInfo + "Failed getting for current day, trying to get from yesterday instead (it's before 17:00)")

		errRetCurr := DB.GetCurrByDate(time.Now().AddDate(0, 0, -1).Format("2006-01-02"), &currInfo)
		if errRetCurr.Status != 0 {
			return errRetCurr
		}
	}

	// Checking target currency is valid
	if currInfo.Rates[string(jsonReq.Result.Parameters.Currency[1])] == 0 {
		return utils.CustError{http.StatusBadRequest, "Did not recognice the target currency"}
	}

	// Creating response
	http.Header.Add(w.Header(), "content-type", "text/plain")

	resp := types.DialogFlowResp{}

	resp.Speech = strconv.FormatFloat(currInfo.Rates[jsonReq.Result.Parameters.Currency[1]], 'f', -1, 64)
	resp.DisplayText = resp.Speech + " disp"

	// Send respons
	log.Println(logInfo + "Trying to encode response from dialogflow to browser")
	errEncode := json.NewEncoder(w).Encode(&resp)
	if errEncode != nil {
		log.Println(logError + "Failed to encode response to browser!")
		return utils.CustError{http.StatusInternalServerError, "Could not encoding result"}
	}

	// Nothing bad happened
	log.Println(logInfo + "Average finished successfully")
	return utils.CustError{0, utils.ErrorStr[0]}
}

// HandlerRoot - Root path handler, handles registration, viewing and
// deleting of webhooks
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if utils.CheckPrintErr(procGetWebHook(r.URL.Path, w), w) {
			return
		}
	case "POST":
		if utils.CheckPrintErr(procAddWebHook(r, w), w) {
			return
		}
	case "DELETE":
		if utils.CheckPrintErr(procDelWebHook(r.URL.Path, w), w) {
			return
		}
	default:
		utils.CheckPrintErr(utils.CustError{http.StatusMethodNotAllowed, utils.ErrorStr[9]}, w)
		return
	}
}

// HandlerLatest - Latest Currency Exchange handler
func handlerLatest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if utils.CheckPrintErr(procLatest(r, w), w) {
			return
		}
	default:
		utils.CheckPrintErr(utils.CustError{http.StatusMethodNotAllowed, utils.ErrorStr[9]}, w)
		return
	}
}

// HandlerAverage - 3-day Average Currency Exchange handler
func handlerAverage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if utils.CheckPrintErr(procAverage(r, w), w) {
			return
		}
	default:
		utils.CheckPrintErr(utils.CustError{http.StatusMethodNotAllowed, utils.ErrorStr[9]}, w)
		return
	}
}

// HandlerEvalTrigger - Evaluation Trigger handler, used when service is evaluated
func handlerEvalTrigger(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		DB.InvokeWebhooks(false)
	default:
		utils.CheckPrintErr(utils.CustError{http.StatusMethodNotAllowed, utils.ErrorStr[9]}, w)
		return
	}
}

// HandlerDialogFlow - DialogFlow entry point, used when dialogflow want to communicate
func handlerDialogFlow(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		utils.CheckPrintErr(procDialogFlow(r, w), w)
	} else {
		utils.CheckPrintErr(utils.CustError{http.StatusMethodNotAllowed, "Expected a POST request"}, w)
	}
}

func main() {
	err := DB.Init()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal(logError + "$PORT was not set")
	}

	// Registration of webhhooks goes over "/root"
	// While viewing and deletion of webhooks goes over "/root/"
	http.HandleFunc("/exchange", handlerRoot)
	http.HandleFunc("/exchange/", handlerRoot)
	http.HandleFunc("/exchange/latest", handlerLatest)
	http.HandleFunc("/exchange/average", handlerAverage)
	http.HandleFunc("/exchange/evaluationtrigger", handlerEvalTrigger)
	http.HandleFunc("/exchange/dialogFlowEntryPoint", handlerDialogFlow)
	log.Printf(logInfo+"Listening on port: %v", port)
	http.ListenAndServe(":"+port, nil)
}
