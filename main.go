package main

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var db *Database // user database

var indexTemplate *template.Template
var aboutTemplate *template.Template
var loginbtnHTML, logoutbtnHTML template.HTML // log in & out buttons

// Store redis connection as a package level variable
var cache redis.Conn

/**
Initialize redis cache for session/user pairs
 */
func initCache() {
	conn, err := redis.DialURL(os.Getenv("REDIS_URL"))
	checkErr(err) // check error

	// assign connection to package level 'cache' variable
	cache = conn
}

/**
Initialize database connection
 */
func initDB() *Database {
	db = NewDatabase(os.Getenv("DATABASE_URL"))
	err := db.GenerateTable()
	checkErr(err)
	return db
}

func main(){
	var err error // declare error variable err to avoid :=
	initCache() // initialize redis cache
	// defer cache.Close() // global var, should be okay without closing, may fix a bug I'm experiencing by not closing

	content, err := ioutil.ReadFile("templates/loginbtn.html")
	checkErr(err)
	loginbtnHTML = template.HTML(string(content))
	content, err = ioutil.ReadFile("templates/logoutbtn.html")
	checkErr(err)
	logoutbtnHTML = template.HTML(string(content))

	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	aboutTemplate = template.Must(template.ParseFiles("templates/about.html"))

	// Connect to database
	db = initDB()
	defer db.conn.Close(context.Background())

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
	r.HandleFunc("/getlifts", getliftsHandler).Methods("GET")

	// favicon.ico handler
	r.HandleFunc("/favicon.ico", faviconHandler)

	// Get latest user data (more specifically, other users than oneself)
	//r.HandleFunc("/user", userHandler).Methods("POST")

	// file directory for file serving
	staticFileDirectory := http.Dir("./static/")
	// the prefix is the routing address, the address the user goes to
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

	// keep PathPrefix empty
	r.PathPrefix("/").Handler(staticFileHandler).Methods("GET")

	port := os.Getenv("PORT")

	if port == "" {
		// if running locally
		port = "8000"
		authconf.RedirectURL = "http://localhost:8000/callback"
	}

	http.ListenAndServe(":" + port, r)
}