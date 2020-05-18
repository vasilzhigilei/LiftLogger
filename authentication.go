package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

conf := &oauth2.Config{
	ClientID:     c.Cid,
	ClientSecret: c.Csecret,
	RedirectURL:  "http://localhost:8000/auth",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	},
	Endpoint: google.Endpoint,
}
type Credentials struct {
	Cid string `json:"cid"`
	Csecret string `json:"csecret"`
}

func init() {
	var c Credentials
	file, err := ioutil.ReadFile("./creds/creds.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &c)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthStateString := uniuri.New()
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}