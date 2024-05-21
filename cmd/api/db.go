package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// *sql.DB - represents a database connection, or connection pool
// pointer to database connection
// *sql.DB - SAFE FOR GO CONCURRENT USE BY MULT GOROUTINES

// This 'layout' of code allows - in the future - to add another `openDB` method for another database as needed
// Of course, we would need another dbrepo/ file, e.g. `mongo_dbrepo.go` for a connection to MongoDb for example (with another driver haha)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	conn, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to postgres!")
	return conn, nil
}
