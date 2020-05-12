package main

import (
	"context"
	"github.com/jackc/pgx"
)

// happy to know I can return a pointer without worrying about the object deallocating :)
func db_connect(connString string) *pgx.Conn{
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	return conn
}