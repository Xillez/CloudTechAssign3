package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var errorStr = []string{
	"No Error",
	"Invalid Server",
	"Status Code Out Of Bounds",
	"Invalid Vital Data",
	"Faliure Reading File",
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

// CheckPrintErr - Checks for error, print it if any and returns true, otherwise returns false
func checkPrintErr(err CustError, w http.ResponseWriter) bool {
	if err.Status != 0 {
		http.Error(w, http.StatusText(err.Status)+" | Program error: "+err.Msg, err.Status)
		return true
	}

	// Say that every thing went ok
	w.WriteHeader(http.StatusOK)
	return false
}

// Get the URL given to server and splits it for processing
func getSplitURL(url string, expectedNrSplits int) ([]string, CustError) {
	parts := strings.Split(url, "/")

	// Missing a field/part of URL
	if len(parts) != expectedNrSplits {
		return nil, CustError{http.StatusBadRequest, errorStr[6]}
	}

	// Nothing bad happened
	return parts, CustError{0, errorStr[0]}
}

// Fetches and decodes json into given variable
func fetchDecodedJSON(url string, updated interface{}) CustError {
	resp, err := http.Get(url)
	if err != nil {
		return CustError{http.StatusBadRequest, errorStr[6]}
	}

	// Decode
	err = json.NewDecoder(resp.Body).Decode(&updated)
	if err != nil {
		return CustError{http.StatusInternalServerError, errorStr[8]}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
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
	return CustError{0, errorStr[0]}
}*/
