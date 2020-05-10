package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// plain text db details, for now
const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "password"
	dbname   = "liftlogger"
)

func main(){
	// prints postgresql info to string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// opens sql connection to db
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defers the closing of the db connection to end of program
	defer db.Close()
	// check if connected by pinging db
	err = db.Ping()
	if err != nil {
		panic(err)
	}
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
