// Package datastore provides an abstraction for different flavors of data storage

package datastore

import (
	"github.com/rknizzle/livetest/pkg/config"
)

// Datastore abstraction
type Datastore interface {
	Connect(config.DatastoreConfig)
	Write(*Record)
}

// Row of data to write to the database table
type Record struct {
	Success    bool
	Title      string
	StatusCode int
}
