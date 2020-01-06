package util

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDatabaseConn() *sqlx.DB {
	conn, err := sqlx.Connect("postgres", "host=localhost user=postgres password=mysecretpassword dbname=prodapp sslmode=disable")
	if err != nil {
		panic(err)
	}
	return conn
}
