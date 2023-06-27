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

//go:embed sql/select-archives-by-series.sql
var SQL_SELECT_ARCHIVES_BY_SERIES string;

//go:embed sql/select-archive-by-id.sql
var SQL_SELECT_ARCHIVE_BY_ID string;

//go:embed sql/select-archive-by-path.sql
var SQL_SELECT_ARCHIVE_BY_PATH string;

func PersistArchive(archive *model.Archive) (bool, error) {
    if (archive.RelPath == "") {
        return false, fmt.Errorf("Persisting archive requires a RelPath.");
    }

    if ((archive.Series == nil) || (archive.Series.Name == "")) {
        return false, fmt.Errorf("Persisting archive requires a series name.");
    }

    dbArchive, err := FetchArchiveByPath(archive.RelPath);
    if (err != nil) {
        return false, err;
    }

    if (dbArchive != nil) {
        // This archive already exists in the DB, just use the DB version.
        archive.Assume(dbArchive);
        return true, nil;
    }

    // The archive does not exist in the db, add it.
    return false, insertArchive(archive);
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

    result, err := statement.Exec(
        archive.Series.ID,
        archive.RelPath,
        archive.Volume,
        archive.Chapter,
        archive.PageCount,
        archive.CoverImageRelPath,
    );

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

func FetchArchiveByID(id int) (*model.Archive, error) {
    statement, err := db.Prepare(SQL_SELECT_ARCHIVE_BY_ID);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    archive, err := scanArchive(statement.QueryRow(id));
    if (err != nil) {
        return nil, err;
    }

    return archive, nil;
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

func FetchArchivesBySeries(seriesID int) ([]*model.Archive, error) {
    statement, err := db.Prepare(SQL_SELECT_ARCHIVES_BY_SERIES);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    rows, err := statement.Query(seriesID);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    var archives = make([]*model.Archive, 0);

    for (rows.Next()) {
        archive, _, err := scanArchiveNoSeries(rows);
        if (err != nil) {
            return nil, err;
        }

        archives = append(archives, archive);
    }

    return archives, nil;
}

func scanArchive(scanner RowScanner) (*model.Archive, error) {
    var archive = model.EmptyArchive("");

    err := scanner.Scan(
            &archive.ID,
            &archive.RelPath,
            &archive.Volume,
            &archive.Chapter,
            &archive.PageCount,
            &archive.CoverImageRelPath,
            &archive.Series.ID,
            &archive.Series.Name,
            &archive.Series.AltNames,
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

func scanArchiveNoSeries(scanner RowScanner) (*model.Archive, int, error) {
    var archive = model.Archive{};
    var seriesID int;

    err := scanner.Scan(
            &archive.ID,
            &seriesID,
            &archive.RelPath,
            &archive.Volume,
            &archive.Chapter,
            &archive.PageCount,
            &archive.CoverImageRelPath,
    );

    return &archive, seriesID, err;
}
