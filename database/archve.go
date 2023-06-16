package database

import (
    _ "embed"
    "fmt"

    "github.com/eriq-augustine/comic-server/model"
)

//go:embed sql/insert-archive.sql
var SQL_INSERT_ARCHIVE string;

//go:embed sql/select-archives.sql
var SQL_SELECT_ARCHIVES string;

//go:embed sql/select-archive-by-path.sql
var SQL_SELECT_ARCHIVE_BY_PATH string;

// Ensure that each archive exists in the database.
// Whether from the passed in archive or the database,
// each archive will have its latest metadata attatched when done.
func PersistArchives(archives []*model.Archive) error {
    for _, archive := range archives {
        err := persistArchive(archive);
        if (err != nil) {
            return err;
        }
    }

    return nil;
}

// TODO(eriq): This function does not consider updating an existing archive.
func persistArchive(archive *model.Archive) error {
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

func insertArchive(archive *model.Archive) error {
    transaction, err := db.Begin();
    if (err != nil) {
        return err
    }
    defer transaction.Rollback();

    err = ensureSeries(archive);
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

    transaction.Commit();
    return nil;
}

func FetchArchives() ([]*model.Archive, error) {
    rows, err := db.Query(SQL_SELECT_ARCHIVES);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    var archives = make([]*model.Archive, 0);

    for (rows.Next()) {
        archive, err := scanArchive(rows);
        if (err != nil) {
            return nil, err;
        }

        archives = append(archives, archive);
    }

    return archives, nil;
}

func FetchArchiveByPath(path string) (*model.Archive, error) {
    statement, err := db.Prepare(SQL_SELECT_ARCHIVE_BY_PATH);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

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

    archive, err := scanArchive(rows);
    if (err != nil) {
        return nil, err;
    }

    return archive, nil;
}

func scanArchive(scanner RowScanner) (*model.Archive, error) {
    var archive = model.EmptyArchive("");

    err := scanner.Scan(
            &archive.ID,
            &archive.Path,
            &archive.Volume,
            &archive.Chapter,
            &archive.PageCount,
            &archive.Series.ID,
            &archive.Series.Name,
            &archive.Series.Author,
            &archive.Series.Year,
            &archive.Series.URL,
            &archive.Series.Description,
            &archive.Series.CoverImageRelPath,
            &archive.Series.MetadataSource,
            &archive.Series.MetadataSourceID,
    );

    return archive, err;
}
