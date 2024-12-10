/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)


type AppDatabase interface {
    Ping() error
    
    // User management
    CheckUserExists(username string) (bool, error)
    CreateUser(username string, password string) (string, error) // Returns identifier
    GetUserByCredentials(username string, password string) (string, error) // Returns identifier
}

type appdbimpl struct {
    c *sql.DB
}

func New(db *sql.DB) (AppDatabase, error) {
    if db == nil {
        return nil, errors.New("database is required when building a AppDatabase")
    }

    // Create users table if it doesn't exist
    sqlStmt := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        identifier TEXT UNIQUE NOT NULL
    );`
    
    _, err := db.Exec(sqlStmt)
    if err != nil {
        return nil, fmt.Errorf("error creating database structure: %w", err)
    }

    return &appdbimpl{
        c: db,
    }, nil
}

func (db *appdbimpl) Ping() error {
    return db.c.Ping()
}