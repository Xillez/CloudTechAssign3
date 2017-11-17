package utils

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_Pos_CheckPrintErr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test := CheckPrintErr(CustError{0, ErrorStr[0]}, w)
		if test {
			t.Error("\"checkPrintErr\" gave true at no error!")
		}
	}))
	defer server.Close()

	_, err := http.Get(server.URL)
	if err != nil {
		t.Error("\"GET\" request to testserver failed! Error: " + err.Error())
	}
}

func Test_Neg_checkPrintErr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		test := CheckPrintErr(CustError{http.StatusInternalServerError, ErrorStr[6]}, w)
		if !test {
			t.Error("\"checkPrintErr\" gave false at error!")
		}
	}))
	defer server.Close()

	_, err := http.Get(server.URL)
	if err != nil {
		t.Error("\"GET\" request to testserver failed! Error: " + err.Error())
	}
}

// Positive test, getSplitURL
func Test_Pos_getSplitURL(t *testing.T) {
	testURL := "localhost:8080/item1/item2/item3"
	testFrag := []string{"localhost:8080", "item1", "item2", "item3"}
	frag, err := GetSplitURL(testURL, 4)

	if err.Status != 0 {
		t.Error("getSplitURL returned an error: " + strconv.Itoa(err.Status) + " | " + err.Msg)
	}

	if len(frag) != 4 {
		t.Error("Unexpected nr of fragments! Nr: " + strconv.Itoa(len(frag)))
	}

	for i := 0; i < 4; i++ {
		if frag[i] != testFrag[i] {
			t.Error("Fragment " + strconv.Itoa(i) + "aren't similar!")
		}
	}
}

// Negative test, getSplitURL
func Test_Neg_getSplitURL(t *testing.T) {
	testURL := "localhost:8080/item1/item2/item3"
	_, err := GetSplitURL(testURL, 3)

	if err.Status == 0 {
		t.Error("getSplitURL didn't return an error! URL: " + testURL + " | Expected 3 Splits!")
	}

	if err.Status != http.StatusBadRequest {
		t.Error("Status code" + strconv.Itoa(err.Status) + " isn't " + strconv.Itoa(http.StatusBadRequest) + "! URL: " + testURL)
	}
}

// Positive test, fetchDecodedJSON
/*func Test_Pos_fetchDecodedJSON(t *testing.T) {
	testBaseCurr := "USD"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encode := CurrencyReq{testBaseCurr, testBaseCurr}
		fmt.Println("-------------------server.encode.BaseCurrency -------------\n" + encode.BaseCurrency)
		err := json.NewEncoder(w).Encode(&encode)

		if err != nil {
			t.Error("Server faild to encode! Error: " + err.Error())
		}
	}))
	defer server.Close()

	resp := CurrencyReq{}

	err := fetchDecodedJSON(server.URL, resp)

	fmt.Println(resp.BaseCurrency + " | " + resp.TargetCurrency)

	if err.Status != 0 {
		t.Error("fetchDecodedJSON returned an error: " + strconv.Itoa(err.Status) + " | \"" + err.Msg + "\"")
	}

	if resp.BaseCurrency != testBaseCurr {
		t.Error("Fetched \"baseCurrency\": " + resp.BaseCurrency + " isn't equal to \"USD\"!")
	}
}*/

// Negative test, fetchDecodedJSON
func Test_Neg_fetchDecodedJSON(t *testing.T) {
	testURL := "http://www..com"

	err := FetchDecodedJSON(testURL, nil)

	if err.Status == 0 {
		t.Error("fetchDecodedJSON didn't return an error!")
	}

	if err.Status != http.StatusBadRequest {
		t.Error("Status code" + strconv.Itoa(err.Status) + " isn't " + strconv.Itoa(http.StatusInternalServerError))
	}
}
