package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/**
Index page handler
This is where all the cool things happen. User can log in, log lifts, modify reps display (saved as cookie),
view various rep maximums for the four big compound lifts!
 */
func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Username:    "Not Logged In",
		Loginoutbtn: loginbtnHTML,
	}
	c, err := r.Cookie("oauthstate")
	if err != nil {
		// If the session token is not present in cache, set to not logged in
		// For any other type of error, return a bad request status
		if err == http.ErrNoCookie {
			// If the cookie is not set, set to not logged in
			indexTemplate.Execute(w, data)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response == nil {
		indexTemplate.Execute(w, data)
		return
	}else {
		fmt.Println(fmt.Sprintf("%s", response), "has loaded index.html")
		data := db.GetUserLatest(fmt.Sprintf("%s", response))
		data.Username = fmt.Sprintf("%s", response)
		data.Loginoutbtn = logoutbtnHTML
		indexTemplate.Execute(w, data)
	}
}

/**
About page handler
Serves the about page, with the user's email templated into the page
 */
func aboutHandler(w http.ResponseWriter, r *http.Request){
	data := PageData{
		Username:    "Not Logged In",
		Loginoutbtn: loginbtnHTML,
	}
	c, err := r.Cookie("oauthstate")
	if err != nil {
		// If the session token is not present in cache, set to not logged in
		// For any other type of error, return a bad request status
		if err == http.ErrNoCookie {
			// If the cookie is not set, set to not logged in
			aboutTemplate.Execute(w, data)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response == nil {
		aboutTemplate.Execute(w, data)
		return
	}else {
		data.Username = fmt.Sprintf("%s", response)
		data.Loginoutbtn = logoutbtnHTML
		aboutTemplate.Execute(w, data)
	}
}
// global authentication variable
var authconf = &oauth2.Config {
	RedirectURL: "http://liftloggertesting.herokuapp.com/callback",
	ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: google.Endpoint,
}

/**
Loglifts API call handler
Takes submitted form data and logs user data into postgres database
 */
func logliftsHandler(w http.ResponseWriter, r *http.Request){
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response != nil {
		r.ParseForm()
		user := User{
			Email:    fmt.Sprintf("%s", response),
			Sex:      false,
			Age:      myatoi(r.FormValue("Age")),
			Weight:   []float64{myparsefloat(r.FormValue("Weight"))},
			Deadlift: []int{myatoi(r.FormValue("Deadlift"))},
			Squat:    []int{myatoi(r.FormValue("Squat"))},
			Bench:    []int{myatoi(r.FormValue("Bench"))},
			Overhead: []int{myatoi(r.FormValue("Overhead"))},
			Date:     []string{fmt.Sprint(time.Now().Date())},
		}
		err = db.LogLifts(&user)
		checkErr(err)
	}
}

/**
Get lifts API call handler
Returns all lifting data for the user who matches the request session
 */
func getliftsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response != nil {
		user := db.GetUserAll(fmt.Sprintf("%s", response))
		b, err := json.Marshal(user)
		checkErr(err)
		w.Write(b)
	}
}

/**
Struct to accept unmarshaling of Google user data
Can be expanded to accept a large variety of additional user information on Google login
Currently only need email address
 */
type GoogleUser struct {
	Email string `json:"email"`
}

/**
Login handler
Generates random session id, and then redirects client to Google's authentication service
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := generateStateOauthCookie(w)
	url := authconf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

/**
Logout handler
Handles user logout; Deletes session from Redis cache
 */
func logoutHandler(w http.ResponseWriter, r * http.Request) {
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	_, err = cache.Do("DEL", c.Value)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/**
Callback handler for login
Redirected to by Google's authentication service
Receives session ID and email address, sets session/email pair in cache,
and adds user to Postgres user DB if user doesn't already exist
Redirects to index.html
 */
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, _ := authconf.Exchange(oauth2.NoContext, code)

	if !token.Valid(){
		fmt.Fprintln(w, "Retrieved git invalid token")
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

	// insert user into postgresql, auto does check if already exists
	err = db.InsertUser(user.Email)
	checkErr(err)

	err = db.SetDemoData(user.Email)
	checkErr(err)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/**
Serves favicon.ico
 */
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}