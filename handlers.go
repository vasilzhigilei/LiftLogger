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

func logliftsHandler(w http.ResponseWriter, r *http.Request){
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	response, err := cache.Do("GET", c.Value)
	checkErr(err)
	if response != nil {
		r.ParseForm()
		//fmt.Println(r.Form)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := generateStateOauthCookie(w)
	url := authconf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r * http.Request) {
	c, err := r.Cookie("oauthstate")
	checkErr(err)
	_, err = cache.Do("DEL", c.Value)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	// insert user into postgresql, auto does check if already exists
	err = db.InsertUser(user.Email)
	checkErr(err)

	err = db.SetDemoData(user.Email)
	checkErr(err)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func verificationHandler(w http.ResponseWriter, r *http.Request) {
	verificationTemplate.Execute(w, nil)
}