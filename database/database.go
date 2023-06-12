package database

import (
    "database/sql"
    _ "embed"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
    "github.com/rs/zerolog/log"
)

// TODO(eriq): config
const DB_PATH = "comic-server.db"

var db *sql.DB = nil;

//go:embed sql/create.sql
var SQL_CREATE_TABLES string;

func Open() error {
    var err error;
	db, err = sql.Open("sqlite3", DB_PATH);
	if err != nil {
        return fmt.Errorf("Failed to open database %v: %w.", DB_PATH, err);
	}

    return ensureTables();
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

func Close() {
    if (db == nil) {
        return;
    }

    err := db.Close();
    if (err != nil) {
        log.Error().Err(err);
    }

    db = nil;
}
