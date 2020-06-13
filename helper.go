package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

/**
Page data for index.html templating
 */
type PageData struct {
	Username string
	Loginoutbtn template.HTML
	Sex bool
	Age int
	Weight float64
	DLWeight int
	SWeight int
	BPWeight int
	OHPWeight int
}

/**
User struct, data as would be found in postgres row
 */
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

/**
Atoi that ignores errors (allows for inline struct initialization)
 */
func myatoi(str string) int {
	result, _ := strconv.Atoi(str)
	return result
}

/**
ParseFloat that ignores errors (allows for inline struct initialization)
 */
func myparsefloat(str string) float64 {
	result, _ := strconv.ParseFloat(str, 64)
	return result
}

/**
Generates new session with 1 year expiration time
 */
func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

/**
Unused function for the purpose of debugging cache
 */
func cachePrintAll(){
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

/**
Check error func
 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
