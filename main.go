package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){
	// Declare a new router
	r := mux.NewRouter()

	r.HandleFunc("/hello", handler).Methods("GET")

	// file directory for assets
	staticFileDirectory := http.Dir("./frontend/")
	staticFileHandler := http.StripPrefix("/frontend/", http.FileServer(staticFileDirectory))

	r.PathPrefix("/frontend/").Handler(staticFileHandler).Methods("GET")
	http.ListenAndServe(":8000", r)
}

func handler(w http.ResponseWriter, r *http.Request){
	// simply pipe "Hello World" into response writer
	fmt.Fprintf(w, "Hello World!")
}
