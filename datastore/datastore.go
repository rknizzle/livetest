// Package datastore provides an abstraction for different flavors of data storage

package datastore

import (
	"time"
)

// Datastore abstraction
type Datastore interface {
	Connect(*Connection) error
	Write(*Record)
}

// Row of data to write to the database table
type Record struct {
	Success    bool
	Title      string
	StatusCode int
	Duration   time.Duration
}

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}
