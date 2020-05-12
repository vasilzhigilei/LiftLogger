package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main(){
	var err error // declare error variable err to avoid :=

	// Connect to database
	db := NewDatabase("postgres://postgres:password@localhost:5433/liftlogger")

	// test insert new user
	err = db.InsertUser("hi3@hi.com", true)
	if err != nil{
		panic(err)
	}
	rows := db.SelectAll()
	for rows.Next() {
		var email string
		var sex bool
		var weight float32
		var age float32
		rows.Scan(&email, &sex, &weight, &age)
		fmt.Printf("%s %t %f %f", email, sex, weight, age)
	}
	fmt.Println(rows)

	// Declare a new router
	r := mux.NewRouter()

	r.HandleFunc("/hello", handler).Methods("GET")

	// file directory for file serving
	staticFileDirectory := http.Dir("./assets/")
	// the prefix is the routing address, the address the user goes to
	staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))

	// keep PathPrefix empty
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")
	http.ListenAndServe(":8000", r)
}

func handler(w http.ResponseWriter, r *http.Request){
	// simply pipe "Hello World" into response writer
	fmt.Fprintf(w, "Hello World!")
}
