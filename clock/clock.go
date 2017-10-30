package main

import "gopkg.in/mgo.v2/bson"

func currUpdate() CustError {
	currResp := LatestCurrResp{}
	keys := []string{}

	// Fetch latest currencies and decode them
	errCurr := fetchDecodedJSON(fixerURL+"/latest?base=EUR", currResp)
	if errCurr.Status != 0 {
		return errCurr
	}

	// Extract the keys so we can loop through them
	// map[string]float64 to []string convertion
	for k := range currResp.Rates {
		keys = append(keys, k)
	}

	for i := 0; i < len(currResp.Rates); i++ {
		// Update the target currency with it corresponding value
		errUpdate := db.UpdateCurr(keys[0], bson.M{"$set": bson.M{"rate": currResp.Rates[keys[0]]}})
		// Return imidietly any arrers occur
		if errUpdate.Status != 0 {
			return errUpdate
		}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

func main() {

}
