package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"math/rand"
	"time"
)

/**
Database struct to encase pgx.Conn with SQL command methods
 */
type Database struct{
	conn *pgx.Conn
}

/**
Creates new postgres connection and returns as a Database struct
Current program architecture only takes advantage of a single postgres connection,
in the real world, with more users visiting, this system would have to be rewritten
to increase postgres connections and pool queries
 */
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
Generates table if it doesn't yet exist
Useful and convenient when deploying on a new system
 */
func (d *Database) GenerateTable() error {
	execstring := `
CREATE TABLE IF NOT EXISTS userdata (
email text NOT NULL UNIQUE,
sex boolean NOT NULL DEFAULT false,
age integer NOT NULL DEFAULT 0,
weight double precision[] NOT NULL DEFAULT ARRAY[]::double precision[],
deadlift integer[] NOT NULL DEFAULT ARRAY[]::integer[],
squat integer[] NOT NULL DEFAULT ARRAY[]::integer[],
bench integer[] NOT NULL DEFAULT ARRAY[]::integer[],
overhead integer[] NOT NULL DEFAULT ARRAY[]::integer[],
date text[] NOT NULL DEFAULT ARRAY[]::text[]
);
`
	_, err := d.conn.Exec(context.Background(), execstring)
	return err
}

/**
Insert user, if conflict (user already exists), do nothing
 */
func (d *Database) InsertUser(email string) error {
	_, err := d.conn.Exec(context.Background(), "INSERT INTO userdata values($1) ON CONFLICT DO NOTHING", email)
	return err
}

/**
Log lifts into user row. Appends data and log date to corresponding arrays in the row
If today's date already exists, replace data for this date
 */
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

/**
Get latest user data. Useful for index.html templating of the weight/reps input fields
 */
func (d *Database) GetUserLatest(email string) *PageData{
	querystring := "SELECT sex, age, weight[array_upper(weight, 1)], deadlift[array_upper(deadlift, 1)], " +
		"squat[array_upper(squat, 1)], bench[array_upper(bench, 1)], overhead[array_upper(overhead, 1)] " +
		"FROM userdata WHERE email = '" + email + "';"
	rows, err := d.conn.Query(context.Background(), querystring)
	checkErr(err)
	pagedata := PageData{}
	for rows.Next() {
		err = rows.Scan(&pagedata.Sex, &pagedata.Age, &pagedata.Weight, &pagedata.DLWeight, &pagedata.SWeight,
			&pagedata.BPWeight, &pagedata.OHPWeight)
		//checkErr(err) okay soooo... it does have an error, but if you don't check it code works great haha :)
		// reason for err is if there is an empty array (all users who haven't logged a certain lift)
	}
	return &pagedata
}

/**
Get ALL of a user's data. Useful for chart AJAX call
 */
func (d *Database) GetUserAll(email string) *User{
	querystring := "SELECT sex, age, weight, deadlift, squat, bench, overhead, date " +
		"FROM userdata WHERE email = '" + email + "';"
	rows, err := d.conn.Query(context.Background(), querystring)
	checkErr(err)
	user := User{
		Email:    email,
	}
	for rows.Next() {
		err = rows.Scan(&user.Sex, &user.Age, &user.Weight, &user.Deadlift, &user.Squat, &user.Bench,
			&user.Overhead, &user.Date)
		checkErr(err)
	}
	return &user
}

/**
Unused debugging function to select all users in the table
Selects email, sex, weight, and age from rows
 */
func (d *Database) SelectAllUsers() pgx.Rows{
	rows, _ := d.conn.Query(context.Background(), "SELECT email, sex, weight, age FROM userdata")
	return rows
}

/**
Unused debugging function that prints all users in the table
Prints all rows, only email, sex, weight, and age columns
 */
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

/**
Auto-generate random lifting data and log in database
Used every time user logs in when website set to "Demo mode"
 */
func (d *Database) SetDemoData(email string) error {
	// number of days to generate
	days := 100
	user := User{
		Email:    email,
		Sex:      false,
		Age:      21,
		Weight:   make([]float64, days),
		Deadlift: make([]int, days),
		Squat:    make([]int, days),
		Bench:    make([]int, days),
		Overhead: make([]int, days),
		Date:     make([]string, days),
	}
	// initial values to base generation from
	user.Weight[0] = 185.6
	user.Deadlift[0] = 225
	user.Squat[0] = 205
	user.Bench[0] = 155
	user.Overhead[0] = 95
	user.Date[0] = fmt.Sprint((time.Now().AddDate(0, 0, -days*2)).Date())

	for i := 1; i < days; i++ {
		user.Weight[i] = myparsefloat(fmt.Sprintf("%.1f", (user.Weight[0] + user.Weight[i - 1] +
			float64(rand.Intn(15)-7))/2))
		user.Deadlift[i] = user.Deadlift[i - 1] + (rand.Intn(8) - 2)*rand.Intn(2)
		user.Squat[i] = user.Squat[i - 1] + (rand.Intn(7) - 2)*rand.Intn(2)
		user.Bench[i] = user.Bench[i - 1] + (rand.Intn(5) - 1)*rand.Intn(2)
		user.Overhead[i] = user.Overhead[i - 1] + (rand.Intn(5) - 1)*rand.Intn(2)
		user.Date[i] = fmt.Sprint((time.Now().AddDate(0, 0, (-days + i)*2)).Date())
	}

	execstring := `
UPDATE userdata
SET age = $1, weight = $2, deadlift = $3, squat = $4, bench = $5, overhead = $6, date = $7
WHERE email = $8;`

	_, err := d.conn.Exec(context.Background(), execstring, user.Age, user.Weight, user.Deadlift, user.Squat,
		user.Bench, user.Overhead, user.Date, user.Email)
	return err
}