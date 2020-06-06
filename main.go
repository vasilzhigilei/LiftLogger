package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// global authentication variable
var authconf = &oauth2.Config{
	RedirectURL: "http://localhost:8000/callback",
	ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

type GoogleUser struct {
	ID string `json:"id"`
	Email string `json:"email"`
	VerifiedEmail bool `json:"verified_email"`
	Name string `json:"name"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Link string `json:"link"`
	Picture string `json:"picture"`
	Gender string `json:"gender"`
	Locale string `json:"locale"`
}

var db *Database // user database

var indexTemplate *template.Template
var aboutTemplate *template.Template
var loginbtnHTML, logoutbtnHTML template.HTML // log in & out buttons

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

// Store redis connection as a package level variable
var cache redis.Conn

func initCache() {
	//conn, err := redis.DialURL("redis://localhost:6379")
	conn, err := redis.Dial("tcp", "localhost:6379")
	checkErr(err) // check error

	// assign connection to package level 'cache' variable
	cache = conn
}

type PageData struct {
	Username string
	Loginoutbtn template.HTML
	Sex bool
	Age int
	Weight float64
	DLWeight int
	DLReps int
	SWeight int
	SReps int
	BPWeight int
	BPReps int
	OHPWeight int
	OHPReps int
}

type User struct {
	Email string
	Sex bool
	Age int
	Weight []float64
	Deadlift []int
	Squat []int
	Bench []int
	Overhead []int
	Date []string
}

func myatoi(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

func myparsefloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 64)
	return result
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func cachePrintAll(){
	// unused function, primarily for debugging
	keys, err := redis.Strings(cache.Do("KEYS", "*"))
	checkErr(err)
	for _, key := range keys {
		fmt.Println(key)
		value, err := cache.Do("GET", key)
		checkErr(err)
		fmt.Println(fmt.Sprintf("%s", value))
	}
	fmt.Println()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}