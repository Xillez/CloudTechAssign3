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
	"regexp"
	"strconv"
	"time"

	"github.com/Xillez/CloudTechAssign3/mongodb"
	"github.com/Xillez/CloudTechAssign3/types"
	"github.com/Xillez/CloudTechAssign3/utils"

	"gopkg.in/mgo.v2/bson"
)

var Warn = "[WARNING]: "
var Error = "[ERROR]: "
var Info = "[INFO]: "

var DB = &mongodb.MongoDB{"mongodb://admin:assign3@ds157185.mlab.com:57185/assignment3", "assignment3", "webhook", "curr"}

/* ---------- Root Handler functions ---------- */

func testing_to_see_if_heroku_will_FUCKIN_UPDATE() {

}

func procGetWebHook(url string, w http.ResponseWriter) utils.CustError {
	log.Println(Info + "--------------- Got GetWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "application/json")
	webhook := types.WebhookInfo{}

	log.Println(Info + "Splitting incomming URL")
	// Split url into usable components
	splitURL, err := utils.GetSplitURL(url, 3)
	if err.Status != 0 {
		return err
	}

	log.Println(Info + "Validating GetWebhook request")
	// Validate user input just in case of attack
	validID, errReg := regexp.MatchString("^[0-9a-f]{24}$", splitURL[2])
	if errReg != nil {
		// Somehting went wrong! Inform the user!
		log.Println(Warn + "Validating failed! Informing user!")
		return utils.CustError{http.StatusInternalServerError, utils.ErrorStr[10] + ": \"" + splitURL[2] + "\""}
	} else if !validID {
		// Somehting went wrong! Inform the user!
		log.Println(Warn + "Validating failed! Informing user!")
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": Given webhook id is invalid!"}
	}

	log.Println(Info + "Trying to get requested webhook")
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

	log.Println(Info + "Encoding webhook to browser")
	// Dump the fetched webhook to browser
	errEncode := json.NewEncoder(w).Encode(&display)
	if errEncode != nil {
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Encoding failed!"}
	}

	log.Println(Info + "GetWebhook request finnished successfully")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

func procAddWebHook(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(Info + "--------------- Got AddWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	webhook := types.WebhookInfo{}

	// Decode request
	log.Println(Info + "Decoding AddWebhook request")
	errDec := json.NewDecoder(r.Body).Decode(&webhook)
	if errDec != nil {
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Webhook registration decoding failed!"}
	}

	log.Println(Info + "Validating data from recieved AddWebhook request")
	// Validate vital data
	if webhook.BaseCurrency != "EUR" {
		log.Println(Warn + "Validation failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(webhook.TargetCurrency) != 3 {
		log.Println(Warn + "Validating failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}
	if webhook.MinValue < 0 && webhook.MinValue > webhook.MaxValue {
		log.Println(Warn + "Validating failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"minTriggerValue\" need to be between 0 - \"maxTriggerValue\""}
	}
	if webhook.MaxValue < webhook.MinValue {
		log.Println(Warn + "Validating failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"maxTriggerValue\" has to be >= \"minTriggerValue\""}
	}

	log.Println(Info + "converting retrieved webhook data into storage format")
	webhookID := bson.NewObjectId()
	// Add webhook info to database
	err := DB.AddWebhook(types.WebhookInfo{
		webhookID,
		webhook.URL,
		webhook.BaseCurrency,
		webhook.TargetCurrency,
		webhook.MinValue,
		webhook.MaxValue})
	if err.Status != 0 {
		fmt.Println(Error + strconv.Itoa(err.Status) + " | " + err.Msg)
		return err
	}

	log.Println(Info + "Writing stored entry id to browser")
	fmt.Fprintf(w, webhookID.Hex())

	log.Println(Info + "AddWebhook request finished successfully!")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

func procDelWebHook(url string, w http.ResponseWriter) utils.CustError {
	log.Println(Info + "--------------- Got DelWebhook Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")

	// Split url to extract id string
	log.Println(Info + "Splitting incoming URl")
	splitURL, err := utils.GetSplitURL(url, 3)
	if err.Status != 0 {
		return err
	}

	// Validate vital data
	log.Println(Info + "Validating DelWebhook request data")
	validID, errReg := regexp.MatchString("^[0-9a-f]{24}$", splitURL[2])
	if errReg != nil {
		log.Println(Warn + "Validation failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, utils.ErrorStr[10] + ": \"" + splitURL[2] + "\""}
	} else if !validID {
		log.Println(Warn + "Validation failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": Given webhook id is invalid!"}
	}

	log.Println(Info + "Deleting requested webhook")
	// Try to delete webhook entry with given id
	errDel := DB.DelWebhook(splitURL[2])
	if errDel.Status != 0 {
		return errDel
	}

	log.Println(Info + "Informing user that the webhook was deleted")
	fmt.Fprintf(w, "Webhook deleted!")

	log.Println(Info + "DelWebhook request finished successfully")
	// Nothing bad happened
	return utils.CustError{0, utils.ErrorStr[0]}
}

/* ---------- Latest Handler functions ---------- */

func procLatest(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(Info + "--------------- Got Latest Currency Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	currReq := types.CurrencyReq{}
	currInfo := types.CurrencyInfo{}

	// Decode request
	log.Println(Info + "Decoding latest request body")
	errDec := json.NewDecoder(r.Body).Decode(&currReq)
	if errDec != nil {
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Latest currency exchange request decoding failed!"}
	}

	// Validate vital data
	log.Println(Info + "Validating data from latest request")
	if currReq.BaseCurrency != "EUR" {
		log.Println(Warn + "Validation failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(currReq.TargetCurrency) != 3 {
		// Somehting went wrong! Inform the user!
		log.Println(Warn + "Validation failed! Informing user!")
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	log.Println(Info + "Getting currency from database")
	// Retrieve currency info with current date, if error, try with yesterday's dates
	errRetCurr := DB.GetCurrByDate(time.Now().Format("2006-01-02"), &currInfo)
	if errRetCurr.Status != 0 {
		log.Println(Warn + "Failed getting for current day, trying to get from yester day instead (it's before 17:00)")
		errRetCurr := DB.GetCurrByDate(time.Now().AddDate(0, 0, -1).Format("2006-01-02"), &currInfo)
		if errRetCurr.Status != 0 {
			return errRetCurr
		}
	}

	log.Println(Info + "Writing the latest rate of the requested currency")
	fmt.Fprintf(w, strconv.FormatFloat(currInfo.Rates[currReq.TargetCurrency], 'f', -1, 64))

	log.Println(Info + "Latest request finished successfully")
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
		log.Println(Info + "Trying to get currency with date: " + date + " from databse!")
		err := DB.GetCurrByDate(date, &curr[i])
		if err.Status != 0 {
			return 0.0, utils.CustError{http.StatusInternalServerError, utils.ErrorStr[1] + " | Error: Failed to get average!"}
		}

		// Add up given TargetCurrency's rate for averaging
		log.Println(Info + "Adding all " + strconv.Itoa(len(curr)) + " day of rates for TargetCurrency")
		avg += curr[i].Rates[targetCurr]
	}

	log.Println(Info + "Averaging finished successfully")
	// Nothing bad happened
	return avg / float64(len(curr)), utils.CustError{0, utils.ErrorStr[0]}
}

func procAverage(r *http.Request, w http.ResponseWriter) utils.CustError {
	log.Println(Info + "--------------- Got Average Request ---------------")
	http.Header.Add(w.Header(), "content-type", "text/plain")
	curr := types.CurrencyReq{}

	// Decode request
	log.Println(Info + "Decode averge request data")
	errDec := json.NewDecoder(r.Body).Decode(&curr)
	if errDec != nil {
		log.Println(Error + "Decoding failed, Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusInternalServerError, "Failed!"}
	}

	// Validate vital data
	log.Println(Info + "Validating data")
	if curr.BaseCurrency != "EUR" {
		log.Println(Warn + "Validation failed! Informing user!")
		// Somehting went wrong! Inform the user!
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"baseCurrency\" has to be \"EUR\"!"}
	}
	if len(curr.TargetCurrency) != 3 {
		// Somehting went wrong! Inform the user!
		log.Println(Warn + "Validation failed! Informing user!")
		return utils.CustError{http.StatusNotAcceptable, utils.ErrorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	log.Println(Info + "Fetching average")
	avg, err := fetchAverage(curr.BaseCurrency, curr.TargetCurrency)
	if err.Status != 0 {
		return err
	}

	log.Println(Info + "Writing average to browser")
	fmt.Fprintf(w, strconv.FormatFloat(avg, 'f', -1, 64))

	// Nothing bad happened
	log.Println(Info + "Average finished successfully")
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
		utils.CheckPrintErr(utils.CustError{http.StatusBadRequest, utils.ErrorStr[9]}, w)
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
		utils.CheckPrintErr(utils.CustError{http.StatusNotImplemented, utils.ErrorStr[9]}, w)
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
		utils.CheckPrintErr(utils.CustError{http.StatusNotImplemented, utils.ErrorStr[9]}, w)
		return
	}
}

// HandlerEvalTrigger - Evaluation Trigger handler, used when service is evaluated
func handlerEvalTrigger(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		DB.InvokeWebhooks(false)
	default:
		utils.CheckPrintErr(utils.CustError{http.StatusNotImplemented, utils.ErrorStr[9]}, w)
		return
	}
}

func main() {
	err := DB.Init()
	if err != nil {
		panic(err)
	}

	// Registration of webhhooks goes over "/root"
	// While viewing and deletion of webhooks goes over "/root/"
	http.HandleFunc("/exchange", handlerRoot)
	http.HandleFunc("/exchange/", handlerRoot)
	http.HandleFunc("/exchange/latest", handlerLatest)
	http.HandleFunc("/exchange/average", handlerAverage)
	http.HandleFunc("/exchange/evaluationtrigger", handlerEvalTrigger)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
