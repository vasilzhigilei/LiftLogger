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

var indexTemplate *template.Template
var loginbtnHTML, logoutbtnHTML template.HTML // log in & out buttons

type PageData struct {
	Username string
	Loginoutbtn template.HTML
}

func main(){
	var err error // declare error variable err to avoid :=
	initCache() // initialize redis cache
	
	content, err := ioutil.ReadFile("templates/loginbtn.html")
	checkErr(err)
	loginbtnHTML = template.HTML(string(content))
	content, err = ioutil.ReadFile("templates/logoutbtn.html")
	logoutbtnHTML = template.HTML(string(content))
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))

	// Connect to database
	db := NewDatabase("postgres://postgres:password@localhost:5433/liftlogger")

	// test insert new user
	//err = db.InsertUser("example@example.com", true)
	checkErr(err)

	db.PrintAllUsers()

	// Declare a new router
	r := mux.NewRouter()

	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/callback", callbackHandler).Methods("GET")

	r.HandleFunc("/", indexHandler).Methods("GET")

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("oauthstate")
	if err != nil {
		// If the session token is not present in cache, set to not logged in
		// For any other type of error, return a bad request status
		if err == http.ErrNoCookie {
			// If the cookie is not set, set to not logged in
			data := PageData{Username: "Not logged in", Loginoutbtn: loginbtnHTML}
			indexTemplate.Execute(w, data)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response == nil {
		data := PageData{Username: "Not logged in", Loginoutbtn: loginbtnHTML}
		indexTemplate.Execute(w, data)
		return
	}else {
		data := PageData{Username: fmt.Sprintf("%s",response), Loginoutbtn: logoutbtnHTML}
		indexTemplate.Execute(w, data)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := generateStateOauthCookie(w)
	url := authconf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, _ := authconf.Exchange(oauth2.NoContext, code)

	if !token.Valid(){
		fmt.Fprintln(w, "Retrieved invalid token")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	checkErr(err)

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	var user *GoogleUser
	err = json.Unmarshal(contents, &user)
	checkErr(err)

	state, err := r.Cookie("oauthstate")
	checkErr(err)
	_, err = cache.Do("SETEX", state.Value, 365 * 24 * 60 * 60, user.Email)
	checkErr(err)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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