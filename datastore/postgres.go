// Package postgres implements connection, table creation, and writing rows
// to a postgres database

package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rknizzle/livetest/datastore"
	"time"
)

// Postgres struct implements Datastore interface
type Postgres struct {
	db *sql.DB
}

// Connect to a postgres database and create
func (p *Postgres) Connect(config *datastore.Connection) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return errors.New("Failed to connect to postgres database")
	}

	err = db.Ping()
	if err != nil {
		return errors.New("Failed to ping database")
	}
	// create the table to store request execution results
	// to track results over time
	createStatement := `
    CREATE TABLE IF NOT EXISTS request_data (
      time timestamp,
      title TEXT,
      success BOOlEAN,
      status_code INT,
      duration double precision,
      PRIMARY KEY (time, title)
    );`

	_, err = db.Exec(createStatement)
	if err != nil {
		return errors.New("Failed to execute create statement")
	}
	// set the postgres database object
	p.db = db
	return nil
}

// Write a requests execution data to a new row
func (p *Postgres) Write(r *datastore.Record) {
	// get duration of request in milliseconds
	duration := float64(r.Duration) / float64(time.Millisecond)

	sqlStatement := fmt.Sprintf(`
    INSERT INTO request_data (time, title, success, status_code, duration)
    VALUES (NOW(), '%s', '%t', %d, %f)`, r.Title, r.Success, r.StatusCode, duration)

	_, err := p.db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
