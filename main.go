package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main(){
	conn := db_connect("postgres://postgres:password@localhost:5433/liftlogger")
	_, err := conn.Exec(context.Background(), "INSERT INTO userdata values($1, $2, $3, $4)",
		"example@example.com", false, 200, 19)
	if err != nil{
		panic(err)
	}
	stuff, err := conn.Exec(context.Background(), "SELECT * FROM userdata")
	if err != nil{
		panic(err)
	}
	fmt.Println(stuff)
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
