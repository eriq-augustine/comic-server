package database

import (
    _ "embed"
    "fmt"

    "github.com/eriq-augustine/comic-server/types"
)

//go:embed sql/select-archive-by-path.sql
var SQL_SELECT_ARCHIVE_BY_PATH string;

//go:embed sql/select-series-by-name.sql
var SQL_SELECT_SERIES_BY_NAME string;

//go:embed sql/insert-series.sql
var SQL_INSERT_SERIES string;

//go:embed sql/insert-archive.sql
var SQL_INSERT_ARCHIVE string;

// Ensure that each archive exists in the database.
// Whether from the passed in archive or the database,
// each archive will have its latest metadata attatched when done.
func PersistArchives(archives []*types.Archive) error {
    transaction, err := db.Begin();
    if (err != nil) {
        return err
    }
    defer transaction.Rollback();

    for _, archive := range archives {
        err = persistArchive(archive);
        if (err != nil) {
            return err;
        }
    }

    transaction.Commit();
    return nil;
}

// TODO(eriq): This function does not consider updating an existing archive.
func persistArchive(archive *types.Archive) error {
    if (archive.Path == "") {
        return fmt.Errorf("Persisting archive requires a Path.");
    }

    if ((archive.Series == nil) || (archive.Series.Name == "")) {
        return fmt.Errorf("Persisting archive requires a series name.");
    }

    dbArchive, err := FetchArchiveByPath(archive.Path);
    if (err != nil) {
        return err;
    }

    if (dbArchive != nil) {
        // This archive already exists in the DB, just use the DB version.
        archive.Assume(dbArchive);
        return nil;
    }

    // The archive does not exist in the db, add it.
    return insertArchive(archive);
}

func insertArchive(archive *types.Archive) error {
    err := ensureSeries(archive);
    if (err != nil) {
        return err;
    }

    statement, err := db.Prepare(SQL_INSERT_ARCHIVE);
    if (err != nil) {
        return err;
    }
    defer statement.Close();

    result, err := statement.Exec(archive.Series.ID, archive.Path, archive.Volume, archive.Chapter, archive.PageCount);
    if (err != nil) {
        return err;
    }

    id, err := result.LastInsertId();
    if (err != nil) {
        return err;
    }

    archive.ID = int(id);

    return nil;
}

// If a series already exists with a 100% name match, then use that.
func ensureSeries(archive *types.Archive) error {
    if (archive.Series == nil) {
        return fmt.Errorf("Cannot ensure a nil series.");
    }

    if (archive.Series.Name == "") {
        return fmt.Errorf("Cannot ensure a series without a name.");
    }

    if (archive.Series.ID != -1) {
        return nil;
    }

    series, err := FetchSeriesByName(archive.Series.Name);
    if (err != nil) {
        return err;
    }

    if (series != nil) {
        archive.Series = series;
        return nil;
    }

    return insertSeries(archive.Series);
}

func FetchArchiveByPath(path string) (*types.Archive, error) {
    statement, err := db.Prepare(SQL_SELECT_ARCHIVE_BY_PATH);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    var archive *types.Archive = types.EmptyArchive();
    archive.Path = path;

    rows, err := statement.Query(path);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    if (!rows.Next()) {
        if (rows.Err() != nil) {
            return nil, rows.Err();
        }

        return nil, nil;
    }

    err = rows.Scan(
            &archive.ID,
            &archive.Volume,
            &archive.Chapter,
            &archive.PageCount,
            &archive.Series.ID,
            &archive.Series.Name,
            &archive.Series.Author,
            &archive.Series.Year);
    if (err != nil) {
        return nil, err;
    }

    return archive, nil;
}

func FetchSeriesByName(name string) (*types.Series, error) {
    statement, err := db.Prepare(SQL_SELECT_SERIES_BY_NAME);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    var series *types.Series = types.EmptySeries();
    series.Name = name;

    rows, err := statement.Query(name);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    if (!rows.Next()) {
        if (rows.Err() != nil) {
            return nil, rows.Err();
        }

        return nil, nil;
    }

    err = rows.Scan(
            &series.ID,
            &series.Author,
            &series.Year);
    if (err != nil) {
        return nil, err;
    }

    return series, nil;
}

func insertSeries(series *types.Series) error {
    statement, err := db.Prepare(SQL_INSERT_SERIES);
    if (err != nil) {
        return err;
    }
    defer statement.Close();

    result, err := statement.Exec(series.Name, series.Author, series.Year);
    if (err != nil) {
        return err;
    }

    id, err := result.LastInsertId();
    if (err != nil) {
        return err;
    }

    series.ID = int(id);

    return nil;
}
