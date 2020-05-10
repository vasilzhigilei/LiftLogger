package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/jackc/pgx"
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
	//db := db_connect()
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

// happy to know I can return a pointer without worrying about the object deallocating :)
func db_connect(applicationName string) *pgx.Conn{
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5433/liftlogger")
	if(err != nil){
		panic(err)
	}
	return conn
}