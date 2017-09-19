package main

/*
{
    "name": <value>, e.g. "Tom"
    "age": <value>   e.g. 21
}
*/

import
(
    "strings"
    "net/http"
    "fmt"
    "encoding/json"
    "os"
    //"strconv"
)

var baseUrl = "https://api.github.com/repos/"

// Entire inport stucture
type ProjectInfo struct
{
    Name string       `json:"name"`
    Owner struct {
        Login string  `json:"login"`
    } `json:"owner"`
    Langs []string
    TopContrib []struct {
        Name string   `json:"login"`
        Commits int   `json:"contributions"`
    }
}

// Export structure
type ExportInfo struct {
    Name string         `json:"name"`
    Owner string        `json:"owner"`
    Langs []string      `json:"languages"`
    Contrib string      `json:"contributor"`
    Commits int         `json:"commits"`
}

// Check for error and print it if any
func checkPrintErr(err error/*, customErrCode, customErrMsg string*/) {
    if (err != nil) { panic(err) }
}

func export(w *http.ResponseWriter, export *ProjectInfo) {
    http.Header.Add((*w).Header(), "content-type", "application/json")
    checkPrintErr(json.NewEncoder(*w).Encode(ExportInfo{Name: export.Name,
                                        Owner: export.Owner.Login,
                                        Langs: export.Langs,
                                        Contrib: export.TopContrib[0].Name,
                                        Commits: export.TopContrib[0].Commits}))
}

// Converts a maps keys into []string
func deMapToStringArray(langMap map[string]interface{}) []string {
    array := make([]string, 0, len(langMap))
    for k := range langMap {
        array = append(array, k)
    }
    return array
}

func handlerProjectinfo (w http.ResponseWriter, r *http.Request) {
    project := ProjectInfo{}
    var langMap map[string]interface{}
    parts := strings.Split(r.URL.Path, "/")
    if (len(parts) != 6) {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    resp, err := http.Get(baseUrl + parts[4] + "/" + parts[5])
    checkPrintErr(err)
    defer resp.Body.Close()
    checkPrintErr(json.NewDecoder(resp.Body).Decode(&project))
    resp, err = http.Get(baseUrl + parts[4] + "/" + parts[5] + "/languages")
    checkPrintErr(err)
    checkPrintErr(json.NewDecoder(resp.Body).Decode(&langMap))
    project.Langs = deMapToStringArray(langMap)
    resp, err = http.Get(baseUrl + parts[4] + "/" + parts[5] + "/contributors")
    checkPrintErr(err)
    checkPrintErr(json.NewDecoder(resp.Body).Decode(&project.TopContrib))

    export(&w, &project)
}

func main() {
    //s := Student{ "Tom", 21 }
    fmt.Println(os.Getenv("PORT"))
    http.HandleFunc("/projectinfo/v1/", handlerProjectinfo)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
