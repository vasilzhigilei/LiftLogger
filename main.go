package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

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
	db := db_connect()
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
	var runtimeParams map[string]string
	runtimeParams = make(map[string]string)
	runtimeParams["application_name"] = applicationName
	connConfig := pgx.ConnConfig{
		User: "postgres",
		Password: "password",
		Host: "localhost",
		Port: 5433,
		Database: "liftlogger",
		TLSConfig: nil,
		UseFallbackTLS: false,
		FallbackTLSConfig: nil,
		RuntimeParams: runtimeParams,
	}
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v\n", err)
		os.Exit(1)
	}
	return conn
}