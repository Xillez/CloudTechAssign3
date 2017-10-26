package main

import (
	"net/http"
	"strings"
)

// CheckValidResponse - Checks the response if StatusCode er 2XX and Server is "http://api.github.com"
func checkValidResponse(resp *http.Response) CustError {
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
}

// CheckPrintErr - Checks for error, print it if any and returns true, otherwise returns false
func checkPrintErr(err CustError, w http.ResponseWriter) bool {
	if err.status != 0 {
		http.Error(w, http.StatusText(err.status)+" | Program error: "+err.msg, err.status)
		return true
	}

	return false
}

// Get the URL given to server and splits it for processing
func getSplitURL(r *http.Request, expectedNrSplits int) ([]string, CustError) {
	parts := strings.Split(r.URL.Path, "/")

	// Missing a field/part of URL
	if len(parts) != expectedNrSplits {
		return nil, CustError{http.StatusBadRequest, errorStr[6]}
	}

	// Nothing bad happened
	return parts, CustError{0, errorStr[0]}
}

// Fetches and decodes json into given variable
/*func fetchDecodedJSON(url string, updated interface{}) CustError {
	// Decode
	err = json.NewDecoder(resp.Body).Decode(updated)
	if err != nil {
		return CustError{http.StatusInternalServerError, errorStr[8]}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}*/
