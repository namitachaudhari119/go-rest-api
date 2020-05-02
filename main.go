package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Student struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Standard string `json:"standard"`
}

var students []Student

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range students {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Student{})
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var student Student
	_ = json.NewDecoder(r.Body).Decode(&student)
	student.ID = strconv.Itoa(rand.Intn(1000000))
	students = append(students, student)
	json.NewEncoder(w).Encode(&student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range students {
		if item.ID == params["id"] {
			students = append(students[:index], students[index+1:]...)
			var student Student
			_ = json.NewDecoder(r.Body).Decode(&student)
			student.ID = params["id"]
			students = append(students, student)
			json.NewEncoder(w).Encode(&student)
			return
		}
	}
	json.NewEncoder(w).Encode(students)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range students {
		if item.ID == params["id"] {
			students = append(students[:index], students[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(students)
}

func main() {
	router := mux.NewRouter()

	students = append(students, Student{ID: "1", Name: "Namita", Standard: "nine"})
	students = append(students, Student{ID: "2", Name: "Pratik", Standard: "six"})

	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
