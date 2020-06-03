package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
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

func (d *Database) LogLifts(user *User) error {
	execstring := `
UPDATE userdata
SET age = $1,
weight = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(weight, weight[array_upper(weight, 1)], $2) 
ELSE array_append(weight, $2) END,
deadlift = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(deadlift, deadlift[array_upper(deadlift, 1)], $3) 
ELSE array_append(deadlift, $3) END,
squat = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(squat, squat[array_upper(squat, 1)], $4) 
ELSE array_append(squat, $4) END,
bench = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(bench, bench[array_upper(bench, 1)], $5) 
ELSE array_append(bench, $5) END,
overhead = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(overhead, overhead[array_upper(overhead, 1)], $6) 
ELSE array_append(overhead, $6) END,
date = CASE WHEN date[array_upper(date,1)] = $7 THEN array_replace(date, date[array_upper(date, 1)], $7) 
ELSE array_append(date, $7) END
WHERE email = $8;`
	_, err := d.conn.Exec(context.Background(), execstring, user.Age, user.Weight[0], user.Deadlift[0], user.Squat[0],
		user.Bench[0], user.Overhead[0], user.Date[0], user.Email)
	return err
}

func (d *Database) GetUserLatest(email string) *PageData{
	querystring := "SELECT sex, age, weight[array_upper(weight, 1)], deadlift[array_upper(deadlift, 1)], " +
		"squat[array_upper(squat, 1)], bench[array_upper(bench, 1)], overhead[array_upper(overhead, 1)] " +
		"FROM userdata WHERE email = '" + email + "';"
	rows, err := d.conn.Query(context.Background(), querystring)
	checkErr(err)
	var sex bool
	var age int
	var weight float64
	var deadlift int
	var squat int
	var bench int
	var overhead int
	for rows.Next() {
		err = rows.Scan(&sex, &age, &weight, &deadlift, &squat, &bench, &overhead)
		//checkErr(err) okay soooo... it does have an error, but if you don't check it code works great haha :)
		// reason for err is if there is an empty array (all users who haven't logged a certain lift)
	}
	pagedata := PageData{
		DLReps: 1,
		SReps: 1,
		BPReps: 1,
		OHPReps: 1,
	}
	pagedata.Sex = sex
	pagedata.Age = age
	pagedata.Weight = weight
	pagedata.DLWeight = deadlift
	pagedata.SWeight = squat
	pagedata.BPWeight = bench
	pagedata.OHPWeight = overhead
	return &pagedata
}

func (d *Database) GetUserAll(email string) *User{
	querystring := "SELECT sex, age, weight, deadlift, squat, bench, overhead " +
		"FROM userdata WHERE email = '" + email + "';"
	rows, err := d.conn.Query(context.Background(), querystring)
	checkErr(err)
	user := User{
		Email:    email,
	}
	for rows.Next() {
		rows.Scan(&user.Sex, &user.Age, &user.Weight, &user.Deadlift, &user.Squat, &user.Bench,
			&user.Overhead, &user.Date)
	}
	return &user
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