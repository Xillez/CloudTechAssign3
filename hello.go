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
    "strconv"
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
    fmt.Fprintln(w, "Hello, " +  parts[2] + " " + parts[3] + "!")
}

func replyWithAllStudents (w http.ResponseWriter, db StudentsDB) {
    json.NewEncoder(w).Encode(db.students)
}

func replyWithStudents (w http.ResponseWriter, db StudentsDB, i int) {
    // Make sure that "i" is valid
    if (db.Count() <= i) {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(db.Get(i))
}

func handlerStudent (w http.ResponseWriter, r *http.Request) {
    // --------------
    db := StudentsDB{}
    // --------------
    http.Header.Add(w.Header(), "content-type", "application/json")
    //fmt.Fprintln(w, os.Getenv("PORT"))

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
    i, err := strconv.Atoi(parts[2])
    if (err != nil) {
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    if (parts[2] == "") {
        replyWithAllStudents(&w, &db)
        // Handle "/student/1"
    } else if (parts[2] == "1") {
        replyWithStudent(&w, &db, i)
    }
}

func main() {
    //s := Student{ "Tom", 21 }
    fmt.Println(os.Getenv("PORT"))
    http.HandleFunc("/hello/", handlerHello)
    http.HandleFunc("/student/", handlerStudent)
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}
