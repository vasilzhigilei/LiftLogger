package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
