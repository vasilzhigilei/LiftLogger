package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/", handler)

	http.ListenAndServe(":8000", nil)
}

func handler(w http.ResponseWriter, r *http.Request){
	// simply pipe "Hello World" into response writer
	fmt.Fprintf(w, "Hello World!")
}
