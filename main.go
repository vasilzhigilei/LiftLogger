package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
)

var db *Database // user database

var indexTemplate *template.Template
var aboutTemplate *template.Template
var loginbtnHTML, logoutbtnHTML template.HTML // log in & out buttons

// Store redis connection as a package level variable
var cache redis.Conn

func initCache() {
	//conn, err := redis.DialURL("redis://localhost:6379")
	conn, err := redis.Dial("tcp", "localhost:6379")
	checkErr(err) // check error

	// assign connection to package level 'cache' variable
	cache = conn
}

func main(){
	var err error // declare error variable err to avoid :=
	initCache() // initialize redis cache

	content, err := ioutil.ReadFile("templates/loginbtn.html")
	checkErr(err)
	loginbtnHTML = template.HTML(string(content))
	content, err = ioutil.ReadFile("templates/logoutbtn.html")
	checkErr(err)
	logoutbtnHTML = template.HTML(string(content))

	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	aboutTemplate = template.Must(template.ParseFiles("templates/about.html"))

	// Connect to database
	db = NewDatabase("postgres://postgres:password@localhost:5433/liftlogger")

	// Declare a new router
	r := mux.NewRouter()

	// Login/logout management
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/callback", callbackHandler).Methods("GET")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	// Main page
	r.HandleFunc("/", indexHandler).Methods("GET")
	// About page
	r.HandleFunc("/about", aboutHandler).Methods("GET")

	// API AJAX calls to log lifts or fetch lifting history
	r.HandleFunc("/loglifts", logliftsHandler).Methods("POST")
	r.HandleFunc("/getlifts", getliftsHandler).Methods("POST")

	// Get latest user data (more specifically, other users than oneself)
	//r.HandleFunc("/user", userHandler).Methods("POST")

	// file directory for file serving
	staticFileDirectory := http.Dir("./static/")
	// the prefix is the routing address, the address the user goes to
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

	// keep PathPrefix empty
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")
	http.ListenAndServe(":8000", r)
}