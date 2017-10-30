// Program is made in collaboration with / got help from
//          - Jonas J. Solsvik
// 			- Zohaib Butt
// 			- Eldar Hauge Torkelsen

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//var baseURL = "https://api.github.com/repos/"
var fixerURL = "https://api.fixer.io/"

/* ---------- Root Handler functions ---------- */

func procGetWebHook(url string, w http.ResponseWriter) CustError {
	http.Header.Add(w.Header(), "content-type", "application/json")
	webhook := WebhookInfo{}
	curr := CurrencyInfo{}

	// Slip
	splitURL, err := getSplitURL(url, 2)
	if err.Status != 0 {
		return err
	}

	validID, errReg := regexp.MatchString("^ObjectIdHex\u0028\"[0-9a-f]{24}\"\u0029$", splitURL[1])
	if errReg != nil {
		return CustError{http.StatusInternalServerError, errorStr[10] + ": \"" + splitURL[1] + "\""}
	} else if !validID {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": Given webhook id is invalid!"}
	}

	// Retrieve webhook from database
	err = db.db.GetWebhook(splitURL[1], &webhook)
	if err.Status != 0 {
		return err
	}

	// Retrieve related currency
	err = db.db.GetCurrByID(splitURL[1], &curr)
	if err.Status != 0 {
		return err
	}

	errEncode := json.NewEncoder(w).Encode(&WebhookReg{
		webhook.URL,
		curr.BaseCurrency,
		curr.TargetCurrency,
		webhook.MinValue,
		webhook.MaxValue,
	})
	if errEncode != nil {
		return CustError{http.StatusInternalServerError, "Encoding failed! Error: " + errEncode.Error()}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

func procAddWebHook(r *http.Request, w http.ResponseWriter) CustError {
	http.Header.Add(w.Header(), "content-type", "text/plain")
	webhook := WebhookReg{}
	curr := CurrencyInfo{}

	// Decode request
	errDec := json.NewDecoder(r.Body).Decode(&webhook)
	if errDec != nil {
		return CustError{http.StatusInternalServerError, "Webhook registration decoding failed!"}
	}

	// Validate vital data
	if len(webhook.BaseCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"baseCurrency\" has to be 3 long!"}
	}
	if len(webhook.TargetCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}
	if webhook.MinValue < 0 && webhook.MinValue > webhook.MaxValue {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"minTriggerValue\" need to be between 0 - \"maxTriggerValue\""}
	}
	if webhook.MaxValue < webhook.MinValue {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"maxTriggerValue\" has to be >= \"minTriggerValue\""}
	}

	// Update all currencies, so their fresh for immediate viewing
	currUpdate()

	// Retrieve currency info
	errRetCurr := db.db.GetCurrByTarget(webhook.TargetCurrency, &curr)
	if errRetCurr.Status != 0 {
		return errRetCurr
	}

	webhookID := bson.NewObjectId()
	// Add webhook info to database
	err := db.db.AddWebhook(WebhookInfo{
		webhookID,
		webhook.URL,
		curr.ID,
		webhook.MinValue,
		webhook.MaxValue})
	if err.Status != 0 {
		return err
	}

	fmt.Fprintf(w, webhookID.String())

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

func procDelWebHook(url string, w http.ResponseWriter) CustError {
	http.Header.Add(w.Header(), "content-type", "text/plain")
	// Split url to extract id string
	splitURL, err := getSplitURL(url, 2)
	if err.Status != 0 {
		return err
	}

	// Validate vital data
	validID, errReg := regexp.MatchString("^ObjectIdHex\u0028\"[0-9a-f]{24}\"\u0029$", splitURL[1])
	if errReg != nil {
		return CustError{http.StatusInternalServerError, errorStr[10] + ": \"" + splitURL[1] + "\""}
	} else if !validID {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": Given webhook id is invalid!"}
	}

	// Try to delete webhook entry with given id
	errDel := db.db.DelWebhook(splitURL[1])
	if errDel.Status != 0 {
		return errDel
	}

	fmt.Fprintf(w, "Webhook deleted!")

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

/* ---------- Latest Handler functions ---------- */

func procLatest(r *http.Request, w http.ResponseWriter) CustError {
	http.Header.Add(w.Header(), "content-type", "text/plain")
	currReq := CurrencyReq{}
	currInfo := CurrencyInfo{}

	// Decode request
	errDec := json.NewDecoder(r.Body).Decode(&currReq)
	if errDec != nil {
		fmt.Println(errDec.Error())
		return CustError{http.StatusInternalServerError, "Latest currency exchange request decoding failed!"}
	}

	// Validate vital data
	if len(currReq.BaseCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"baseCurrency\" has to be 3 long!"}
	}
	if len(currReq.TargetCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	// Update all currencies, so their fresh for immediate viewing
	currUpdate()

	// Retrieve currency info
	errRetCurr := db.db.GetCurrByTarget(currReq.TargetCurrency, &currInfo)
	if errRetCurr.Status != 0 {
		return errRetCurr
	}

	fmt.Fprintf(w, strconv.FormatFloat(currInfo.Rate, 'f', -1, 64))

	// Nothing bad happened
	return CustError{0, errorStr[0]}

}

/* ---------- Average Handler functions ---------- */

func fetchAverage(baseCurr string, targetCurr string) (float64, CustError) {
	curr := LatestCurrResp{}
	avg := 0.0

	for i := 0; i > -7; i-- {
		date := time.Now().AddDate(0, 0, i).Format("yyyy-mm-dd")
		err := fetchDecodedJSON(fixerURL+date+"?symbols="+baseCurr+","+targetCurr, curr)
		if err.Status != 0 {
			return 0, CustError{http.StatusInternalServerError, errorStr[1] + " | Error: Couldn't locate at part of the average! | url"}
		}

		avg += curr.Rates[targetCurr]
	}

	return avg / 7, CustError{0, errorStr[0]}
}

func procAverage(r *http.Request, w http.ResponseWriter) CustError {
	http.Header.Add(w.Header(), "content-type", "text/plain")
	curr := CurrencyReq{}

	// Decode request
	errDec := json.NewDecoder(r.Body).Decode(&curr)
	if errDec != nil {
		return CustError{http.StatusInternalServerError, "Latest currency exchange request decoding failed!"}
	}

	// Validate vital data
	if len(curr.BaseCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"baseCurrency\" has to be 3 long!"}
	}
	if len(curr.TargetCurrency) != 3 {
		return CustError{http.StatusNotAcceptable, errorStr[3] + ": \"targetCurrency\" has to be 3 long!"}
	}

	avg, err := fetchAverage(curr.BaseCurrency, curr.TargetCurrency)
	if err.Status != 0 {
		return err
	}

	fmt.Fprintf(w, strconv.FormatFloat(avg, 'f', -1, 64))

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// HandlerRoot - Root path handler, handles registration, viewing and
// deleting of webhooks
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if checkPrintErr(procGetWebHook(r.URL.Path, w), w) {
			return
		}
	case "POST":
		if checkPrintErr(procAddWebHook(r, w), w) {
			return
		}
	case "DELETE":
		if checkPrintErr(procDelWebHook(r.URL.Path, w), w) {
			return
		}
	default:
		checkPrintErr(CustError{http.StatusBadRequest, errorStr[9]}, w)
		return
	}
}

// HandlerLatest - Latest Currency Exchange handler
func handlerLatest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if checkPrintErr(procLatest(r, w), w) {
			return
		}
	default:
		checkPrintErr(CustError{http.StatusNotImplemented, errorStr[9]}, w)
		return
	}
}

// HandlerAverage - 7-day Average Currency Exchange handler
func handlerAverage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if checkPrintErr(procAverage(r, w), w) {
			return
		}
	default:
		checkPrintErr(CustError{http.StatusNotImplemented, errorStr[9]}, w)
		return
	}
}

// HandlerEvalTrigger - Evaluation Trigger handler, used when service is evaluated
func handlerEvalTrigger(w http.ResponseWriter, r *http.Request) {
	switch RunningTest {
	case false:
		// Get general repo data
		/*if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}*/
	case true:
		// Load general testrepo data
		/*if checkPrintErr(openMainJSON(&project), w) {
			return
		}*/
	}
}

func main() {
	err := db.db.Init()
	if err != nil {
		panic(err)
	}
	// Registration of webhhooks goes over "/root"
	// While viewing and deletion of webhooks goes over "/root/"
	http.HandleFunc("/root", handlerRoot)
	http.HandleFunc("/root/", handlerRoot)
	http.HandleFunc("/root/latest", handlerLatest)
	http.HandleFunc("/root/average", handlerAverage)
	http.HandleFunc("/root/evaluationtrigger", handlerEvalTrigger)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
