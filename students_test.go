package main

import "testing"

func Test_addStudent(t *testing.T) {
    db := StudentsDB{}
    db.AddStudent(Student{"Tom", 21})
    if (db.Count() != 1) {
        t.Error("Wrong student count")
    }

    s, err := db.Get(0)

    if (err != nil) {
        t.Error(err)
    }
    if (s.Name != "Tom") {
        t.Error("Student 'Tom' was not added.")
    }
}

func Test_multipleStudents(t *testing.T) {
    testData := []Students {
        Student{"Bob", 21},
        Student{"Alice", 20},
    }

    for _, s := range testData {
        db.Add(s)
    }

    if (db.Count() != lne(testData)) {
        t.Error("Wrong number of students")
    }

    for i, s := range db.students {
        if (db.Get(i).Name != s.Name) {
            t.Error("Wrong name")
        }
    }
}
