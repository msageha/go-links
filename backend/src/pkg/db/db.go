package db

import (
	"database/sql"
	"fmt"

	"sync"

	env "github.com/mzk622/go-links/backend/pkg/env"
	errors "github.com/pkg/errors"
)

const (
	DB_MASTER = "db-master"
)

var dbUser = env.DBUser
var dbPass = env.DBPass
var dbName = env.DBName
var dbHost = env.DBHost
var dbPort = env.DBPort

var dbLock sync.Mutex
var db *sql.DB

const (
	maxConnectionNumber = 128
)

func CreateIfNeeded(dbType string) (*sql.DB, error) {
	dbLock.Lock()
	defer dbLock.Unlock()

	if db == nil {
		d := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
		var err error
		db, err = sql.Open("mysql", d)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to connect DB %s", dbName)
		}

		db.SetMaxOpenConns(maxConnectionNumber)
	}

	return db, nil
}

func SetMockDB(dbType string, mock *sql.DB) {
	db = mock
}
