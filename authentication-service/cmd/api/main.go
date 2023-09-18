package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Unable to connect to DB")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN") // DSN is the data source name and is defined in the docker-compose.yml file

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to DB", err)
			counts++
		} else {
			log.Println("Connected to DB")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Retrying connection to DB")
		time.Sleep(2 * time.Second)
		continue
	}
}
