package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"time"
)

/**
Database struct to encase pgx.Conn with SQL command methods
 */
type Database struct{
	conn *pgx.Conn
}

func NewDatabase(connString string) *Database{
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	d := new(Database)
	d.conn = conn
	// happy to know I can return a pointer without worrying about the object deallocating :)
	return d
}

/**
Insert user, if conflict (user already exists), do nothing
 */
func (d *Database) InsertUser(email string) error {
	_, err := d.conn.Exec(context.Background(), "INSERT INTO userdata values($1) ON CONFLICT DO NOTHING", email)
	return err
}

type LiftData struct {
	email string
	day time.Time
	logs map[string]float32 // allow floats, but lifting data will be stored as ints (weight will be float
}

func (d *Database) LogLifts(ld LiftData) error {
	// construct json to append to json[] in the database
	liftstr := fmt.Sprintf("{\"day\": %s", ld.day.String())
	for key, value := range ld.logs {
		liftstr += fmt.Sprintf(", \"%s\": %f", key, value)
	}
	liftstr += "}"
	_, err := d.conn.Exec(context.Background(), "UPDATE userdata SET lifts = lifts || " + liftstr +
		" WHERE email IS " + ld.email)
	return err
}

type LatestData struct {
	Sex bool
	Age int
	Weight float32
	DLWeight int
	DLReps int
	SWeight int
	SReps int
	BPWeight int
	BPReps int
	OHPWeight int
	OHPReps int
}

func (d *Database) GetUser(email string) *LatestData{
	rows, err := d.conn.Query(context.Background(), "SELECT sex, age, latest FROM userdata WHERE email = " + email)
	checkErr(err)
	var sex bool
	var age int
	var latestlifts []byte
	for rows.Next() {
		err = rows.Scan(&sex, &age, &latestlifts)
		checkErr(err)
	}
	var latestdata *LatestData
	latestdata.Sex = sex
	latestdata.Age = age
	err = json.Unmarshal(latestlifts, &latestdata)
	checkErr(err)
	return latestdata
}

func (d *Database) SelectAllUsers() pgx.Rows{
	rows, _ := d.conn.Query(context.Background(), "SELECT email, sex, weight, age FROM userdata")
	return rows
}

func (d *Database) PrintAllUsers() {
	rows := d.SelectAllUsers()
	for rows.Next() {
		var email string
		var sex bool
		var weight float32
		var age float32
		err := rows.Scan(&email, &sex, &weight, &age)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s %t %f %f\n", email, sex, weight, age)
	}
}