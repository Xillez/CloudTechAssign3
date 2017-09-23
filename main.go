// Program is made in collaboration with / got help from
//          - Jonas J. Solsvik
// 			- Zohaib Butt
// 			- Eldar Hauge Torkelsen

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	//"strconv"
)

var baseURL = "https://api.github.com/repos/"

// RunningTest tells the program wether we're testing or not
var RunningTest = false

var errorStr = []string{
	"NoError",
	"InvalidServer",
	"StatusCodeOutOfBounds",
	"InvalidVitalData",
	"ReadFileFaliure",
	"UnmarshalFail",
	"InvalidURL",
	"EncodeFail",
	"DecodeFail",
}

// CustError - A custum error type
type CustError struct {
	status int
	msg    string
}

// ProjectInfo is the entire inport stucture
type ProjectInfo struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Langs      map[string]interface{}
	TopContrib []struct {
		Name    string `json:"login"`
		Commits int    `json:"contributions"`
	}
}

// ExportInfo is the entire export structure
type ExportInfo struct {
	Name    string   `json:"project"`
	Owner   string   `json:"owner"`
	Langs   []string `json:"language"`
	Contrib string   `json:"committer"`
	Commits int      `json:"commits"`
}

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

func checkVitalData(data *ProjectInfo) CustError {
	// Checks wether vital repo data exists
	if (*data).Name == "" || (*data).Owner.Login == "" {
		return CustError{http.StatusPartialContent, errorStr[3]}
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

// openContribJSON reads "main.json" and decodes it into "updated", writes errors to browser if any errors
func openMainJSON(updated interface{}) CustError {
	file, err := ioutil.ReadFile("main.json")

	// Failed reading the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to load local file: main.json"}
	}

	json.Unmarshal(file, &updated)

	// Failed unmarshaling the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to unmarshal local file: main.json"}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// openLangJSON reads "languages.json" and decodes it into "updated", writes errors to browser if any errors
func openLangJSON(updated interface{}) CustError {
	file, err := ioutil.ReadFile("languages.json")

	// Failed reading the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to load local file: languages.json"}
	}

	err = json.Unmarshal(file, &updated)

	// Failed unmarshaling the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to unmarshal local file: languages.json"}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// openContribJSON reads "contributors.json" and decodes it into "updated", writes errors to browser if any errors
func openContribJSON(updated interface{}) CustError {
	file, err := ioutil.ReadFile("contributors.json")

	// Failed reading the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to load local file: contributors.json"}
	}

	err = json.Unmarshal(file, &updated)

	// Failed unmarshaling the file
	if err != nil {
		return CustError{http.StatusInternalServerError, "Failed to unmarshal local file: contributors.json"}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
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

// Converts and exports to http.ResponseWriter
func export(w http.ResponseWriter, export *ProjectInfo) CustError {
	http.Header.Add(w.Header(), "content-type", "application/json")

	var output = ExportInfo{}

	output.Langs = make([]string, 0, len((*export).Langs))

	// Check vital data
	err := checkVitalData(export)
	if err.status != 0 {
		return CustError{http.StatusPartialContent, ": Some vital values are empty or nonexistent"}
	}

	// map[string]interface{} to []string convertion
	for k := range (*export).Langs {
		output.Langs = append(output.Langs, k)
	}

	sort.Strings(output.Langs)

	output.Name = (*export).Name
	output.Owner = (*export).Owner.Login
	output.Contrib = (*export).TopContrib[0].Name
	output.Commits = (*export).TopContrib[0].Commits

	errEncode := json.NewEncoder(w).Encode(output)

	// Encoding failed
	if errEncode != nil {
		return CustError{http.StatusInternalServerError, errorStr[7]}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

// Fetches and decodes json into given variable
func fetchDecodedJSON(url string, updated interface{}) CustError {
	resp, err := http.Get(url)
	if err != nil {
		return CustError{resp.StatusCode, errorStr[6]}
	}

	// Check if the response is valid
	errResp := checkValidResponse(resp)
	if errResp.status != 0 {
		return errResp
	}

	defer resp.Body.Close()

	// Decode
	err = json.NewDecoder(resp.Body).Decode(updated)
	if err != nil {
		return CustError{http.StatusInternalServerError, errorStr[8]}
	}

	// Nothing bad happened
	return CustError{0, errorStr[0]}
}

func handlerProjectinfo(w http.ResponseWriter, r *http.Request) {
	project := ProjectInfo{}

	// Check if splitting of URL, gives expected nr parts
	parts, errorSplit := getSplitURL(r, 6)
	if checkPrintErr(errorSplit, w) {
		return
	}

	switch RunningTest {
	case false:
		// Get general repo data
		if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5], &project), w) {
			return
		}
		// Get repo language data
		if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5]+"/languages", &project.Langs), w) {
			return
		}
		// Get repo top contributor data
		if checkPrintErr(fetchDecodedJSON(baseURL+parts[4]+"/"+parts[5]+"/contributors", &project.TopContrib), w) {
			return
		}
	case true:
		// Load general testrepo data
		if checkPrintErr(openMainJSON(&project), w) {
			return
		}
		// Load testrepo language data
		if checkPrintErr(openLangJSON(&(project.Langs)), w) {
			return
		}
		// Load testrepo top contributor data
		if checkPrintErr(openContribJSON(&(project.TopContrib)), w) {
			return
		}
	}

	// Check if exported correctly
	if checkPrintErr(export(w, &project), w) {
		return
	}
}

func main() {
	//fmt.Println(os.Getenv("PORT"))
	http.HandleFunc("/projectinfo/v1/", handlerProjectinfo)
	/*log.Println(*/ http.ListenAndServe(":"+os.Getenv("PORT"), nil) /*)*/
}
