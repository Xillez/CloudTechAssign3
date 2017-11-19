package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// FixerURL - Root address to fixer api
var FixerURL = "https://api.fixer.io"

// ErrorStr - Premade error messages for easy use
var ErrorStr = []string{
	"No Error",
	"Invalid Server",
	"Status Code Out Of Bounds",
	"Invalid Vital Data",
	"Faliure Reading From Database",
	"Faliure Unmarshaling",
	"Invalid URL",
	"Faliure Encoding",
	"Faliure Decoding",
	"Method Not Implemented",
	"RegExp String Validation Faliure",
}

// CustError - A custom error type
type CustError struct {
	Status int
	Msg    string
}

const logWarn = "[WARNING]: "
const logError = "[ERROR]: "
const logInfo = "[INFO]: "

// CheckPrintErr - Checks for error, print it if any and returns true, otherwise returns false
func CheckPrintErr(err CustError, w http.ResponseWriter) bool {
	if err.Status != 0 {
		log.Println(logInfo + "Found error to user, displaying...")
		http.Error(w, http.StatusText(err.Status)+" | Program error: "+err.Msg, err.Status)
		return true
	}

	// Say that every thing went ok
	//log.Println(logInfo + "Found no error to report")
	return false
}

// GetSplitURL - Get the URL given to server and splits it for processing
func GetSplitURL(url string, expectedNrSplits int) ([]string, CustError) {
	log.Println(logInfo + "Splitting URL...")
	parts := strings.Split(url, "/")

	// Missing a field/part of URL
	if len(parts) != expectedNrSplits {
		log.Println(logError + "Splitting resulted in not expected nr components!")
		return nil, CustError{http.StatusBadRequest, ErrorStr[6]}
	}

	log.Println(logInfo + "Splitting URL finished successfully")
	// Nothing bad happened
	return parts, CustError{0, ErrorStr[0]}
}

// FetchDecodedJSON - Fetches and decodes json into given variable
func FetchDecodedJSON(url string, updated interface{}) CustError {
	log.Println(logInfo + "Getting from URL: " + url)
	resp, err := http.Get(url)
	if err != nil {
		// Somehting went wrong! Inform the user!
		log.Println(logError + "Couldn't find the URL user requested!")
		return CustError{http.StatusBadRequest, ErrorStr[6]}
	}

	// Decode
	log.Println(logInfo + "Decoding request into given interface")
	err = json.NewDecoder(resp.Body).Decode(&updated)
	if err != nil {
		// Somehting went wrong! Inform the user!
		log.Println(logInfo + "Decoding failed! Inform User!")
		return CustError{http.StatusInternalServerError, ErrorStr[8]}
	}

	log.Println(logInfo + "Fetching and decoding URL finished successfully")
	// Nothing bad happened
	return CustError{0, ErrorStr[0]}
}

// CheckValidResponse - Checks the response if StatusCode er 2XX and Server is "http://api.github.com"
/*func checkValidResponse(resp *http.Response) CustError {
	if resp.StatusCode < 200 || resp.StatusCode > 226 {
		if resp.StatusCode == 404 {
			return CustError{http.StatusNotFound, "Repository NOT FOUND - Check URL or repository details"}
		}
		return CustError{http.StatusBadRequest, "Check URL or repository details"}
	}

	// Treat 206 as error, we're missing some vital repo info
	if resp.StatusCode == 206 {
		return CustError{http.StatusPartialContent, "Repo missing either name or owner"}
	}

	// Nothing bad happened
	return CustError{0, ErrorStr[0]}
}*/
