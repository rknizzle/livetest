// Package postgres implements connection, table creation, and writing rows
// to a postgres database

package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rknizzle/livetest/pkg/config"
	"github.com/rknizzle/livetest/pkg/datastore"
)

// Postgres struct implements Datastore interface
type Postgres struct {
	db *sql.DB
}

// Connect to a postgres database and create
func (p *Postgres) Connect(config config.DatastoreConfig) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// create the table to store request execution results
	// to track results over time
	createStatement := `
    CREATE TABLE IF NOT EXISTS executions (
      time timestamp,
      title TEXT,
      success BOOlEAN,
      status_code INT,
      PRIMARY KEY (time, title)
    );`

	_, err = db.Exec(createStatement)
	if err != nil {
		panic(err)
	}
	// set the postgres database object
	p.db = db
}

// Write a requests execution data to a new row
func (p *Postgres) Write(r *datastore.Record) {

	sqlStatement := fmt.Sprintf(`
    INSERT INTO executions (time, title, success, status_code)
    VALUES (NOW(), '%s', '%t', %d)`, r.Title, r.Success, r.StatusCode)

	_, err := p.db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}
