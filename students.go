package main

type StudentsDB struct {
    students []Student
}

func (db *StudentsDB) Add(s Student) {
    db.students = append(sb.students, s)
}

func (db *StudentsDB) Count int {
    return len(db.students)
}

func (db *StudentsDB) Get (i int) Student {
    if (i < 0 || i >= len(db.students)) {
        // Log error and notify
        return Student{}//, errors.New("Index out of range")
    }
    return db.students[i], nil
}
