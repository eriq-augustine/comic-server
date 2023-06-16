package database

import (
    "database/sql"
    _ "embed"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    _ "github.com/mattn/go-sqlite3"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
)

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

    var dbPath = config.GetString("db.path");
    os.MkdirAll(filepath.Dir(dbPath), 0755);

    var err error;
	db, err = sql.Open("sqlite3", dbPath);
	if err != nil {
        return fmt.Errorf("Failed to open database %v: %w.", dbPath, err);
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
