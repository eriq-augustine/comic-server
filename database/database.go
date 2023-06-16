package database

import (
    "database/sql"
    _ "embed"
    "fmt"
    "sync"

    _ "github.com/mattn/go-sqlite3"
    "github.com/rs/zerolog/log"
)

// TODO(eriq): config
const DB_PATH = "comic-server.db"

var db *sql.DB = nil;
var dbMutex sync.Mutex;

//go:embed sql/create.sql
var SQL_CREATE_TABLES string;

type RowScanner interface {
    Scan(dest ...interface{}) error
}

func Open() error {
    dbMutex.Lock();
    defer dbMutex.Unlock();

    if (db != nil) {
        return nil;
    }

    var err error;
	db, err = sql.Open("sqlite3", DB_PATH);
	if err != nil {
        return fmt.Errorf("Failed to open database %v: %w.", DB_PATH, err);
	}

    return ensureTables();
}

func Close() {
    dbMutex.Lock();
    defer dbMutex.Unlock();

    if (db == nil) {
        return;
    }

    err := db.Close();
    if (err != nil) {
        log.Error().Err(err);
    }

    db = nil;
}

func ensureTables() error {
    var count int;
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master").Scan(&count);
    if (err != nil) {
        return err;
    }

    if (count > 0) {
        // Tables exist, assume the database is fine.
        return nil;
    }

    // No tables exist, create them.
	_, err = db.Exec(SQL_CREATE_TABLES);
	if (err != nil) {
		return fmt.Errorf("Could not create tables: %w.", err);
	}

    return nil;
}
