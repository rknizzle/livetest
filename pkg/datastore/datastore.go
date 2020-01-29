// Package datastore provides an abstraction for different flavors of data storage

package datastore

import (
	"github.com/rknizzle/livetest/pkg/parser"
)

// Datastore abstraction
type Datastore interface {
	Connect(parser.DatastoreConfig)
	Write(*Record)
}

// Row of data to write to the database table
type Record struct {
	Status     string
	Title      string
	StatusCode int
}
