package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Student struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Standard string `json:"standard"`
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var students []Student
	result, err := db.Query("SELECT id, name, standard from student;")
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	for result.Next() {
		var student Student
		err := result.Scan(&student.ID, &student.Name, &student.Standard)
		if err != nil {
			panic(err.Error())
		}
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, name, standard FROM student WHERE id = ?;", params["id"])
	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var student Student
	for result.Next() {
		err := result.Scan(&student.ID, &student.Name, &student.Standard)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(student)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO student(id, name, standard) VALUES(?,?,?);")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	/* decoding json data */
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	standard := keyVal["standard"]
	id := keyVal["id"]

	_, err = stmt.Exec(id, name, standard)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(keyVal)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE student SET name = ?, standard = ? WHERE id = ?;")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	newStd := keyVal["standard"]
	_, err = stmt.Exec(newName, newStd, params["id"])
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(keyVal)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM student WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	keyVal := make(map[string]string)
	keyVal["message"] = "Student with ID = " + params["id"] + " was deleted"
	json.NewEncoder(w).Encode(keyVal)
}

func main() {
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/students")
	if err != nil {
		panic(err.Error())
	}

	/*In Go language, defer statements delay the execution of the function or method or an anonymous
	  method until the nearby functions returns.*/
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
