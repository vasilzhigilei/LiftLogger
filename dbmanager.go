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

func (d *Database) InsertUser(email string, sex bool) error{
	_, err := d.conn.Exec(context.Background(), "INSERT INTO userdata values($1, $2)",
		email, sex)
	return err
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