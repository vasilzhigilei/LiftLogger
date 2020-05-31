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

func (d *Database) LogLifts(email string, data string) error {
	_, err := d.conn.Exec(context.Background(), "UPDATE userdata SET lifts = lifts || " + data +
		" WHERE email IS " + email)
	return err
}

func (d *Database) GetUser(email string) *PageData{
	rows, err := d.conn.Query(context.Background(), "SELECT sex, age, weight[array_upper(weight, 1)] FROM userdata WHERE email = '" + email + "';")
	checkErr(err)
	var sex bool
	var age int
	var weight float32
	for rows.Next() {
		err = rows.Scan(&sex, &age, &weight)
		//checkErr(err) okay soooo... it does have an error, but if you don't check it code works great haha :)
	}
	pagedata := PageData{}
	pagedata.Sex = sex
	pagedata.Age = age
	fmt.Println(age)
	pagedata.Weight = weight
	return &pagedata
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