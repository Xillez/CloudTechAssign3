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
)

type Student struct
{
    Name string `json:"name"`
    Age int     `json:"age"`
}

func handlerHello (w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if (len(parts) != 4) {
        status := 400
        http.Error(w, http.StatusText(status), status)
        return
    }
    name := parts[2]
    fmt.Fprintln(w, parts)
    fmt.Fprintln(w, name)
}

func handlerStudent (w http.ResponseWriter, r *http.Request) {
    http.Header.Add(w.Header(), "content-type", "application/json")

    parts := strings.Split(r.URL.Path, "/")

    // 0
    s0 := Student{"Tom", 21}
    // 1
    s1 := Student{"Alice", 24}
    students := []Student { s0, s1, }

    // Error handling
    if (len(parts) != 3) {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    // Handle "/student/"
    if (len(parts) == 3) {
        if (parts[2] == "") {
            json.NewEncoder(w).Encode(students)
        } else {
            // Handle "/student/1"
            if (parts[2] == "0") {
                json.NewEncoder(w).Encode(s0)
                // Handle "/student/1"
            } else if (parts[2] == "1") {
                json.NewEncoder(w).Encode(s1)
            }
        }
    }
}

func main() {
    //s := Student{ "Tom", 21 }
    http.HandleFunc("/hello/", handlerHello)
    http.HandleFunc("/student/", handlerStudent)
    http.ListenAndServe(os.Getenv("PORT"), nil)
}
