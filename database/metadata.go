package database

import (
    _ "embed"
    "fmt"

    "github.com/eriq-augustine/comic-server/model"
)

//go:embed sql/select-archive-by-path.sql
var SQL_SELECT_ARCHIVE_BY_PATH string;

//go:embed sql/select-crawl-requests.sql
var SQL_SELECT_CRAWL_REQUESTS string;

//go:embed sql/select-series-by-id.sql
var SQL_SELECT_SERIES_BY_ID string;

//go:embed sql/select-series-by-name.sql
var SQL_SELECT_SERIES_BY_NAME string;

//go:embed sql/insert-series.sql
var SQL_INSERT_SERIES string;

//go:embed sql/insert-archive.sql
var SQL_INSERT_ARCHIVE string;

//go:embed sql/upsert-crawl-request.sql
var SQL_UPSERT_CRAWL_REQUEST string;

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

func FetchCrawlRequests() ([]*model.MetadataCrawlRequest, error) {
    rows, err := db.Query(SQL_SELECT_CRAWL_REQUESTS);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    var requests = make([]*model.MetadataCrawlRequest, 0);

    for (rows.Next()) {
        var request *model.MetadataCrawlRequest = model.EmptyCrawlRequest();

        err = rows.Scan(
                &request.ID,
                &request.Query,
                &request.Timestamp,
                &request.Series.ID,
                &request.Series.Name,
                &request.Series.Author,
                &request.Series.Year,
                &request.Series.MetadataSource,
                &request.Series.MetadataSourceID,
        );
        if (err != nil) {
            return nil, err;
        }

        requests = append(requests, request);
    }

    return requests, nil;
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

// If a series already exists with a 100% name match, then use that.
func ensureSeries(archive *model.Archive) error {
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

    err = insertSeries(archive.Series);
    if (err != nil) {
        return err;
    }

    return RequestMetadataCrawl(archive.Series);
}

func RequestMetadataCrawl(series *model.Series) error {
    statement, err := db.Prepare(SQL_UPSERT_CRAWL_REQUEST);
    if (err != nil) {
        return err;
    }
    defer statement.Close();

    _, err = statement.Exec(series.ID, series.Name);
    return err;
}

func FetchArchiveByPath(path string) (*model.Archive, error) {
    statement, err := db.Prepare(SQL_SELECT_ARCHIVE_BY_PATH);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    var archive *model.Archive = model.EmptyArchive();
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
            &archive.Series.Year,
            &archive.Series.MetadataSource,
            &archive.Series.MetadataSourceID,
    );
    if (err != nil) {
        return nil, err;
    }

    return archive, nil;
}

func FetchSeriesByName(name string) (*model.Series, error) {
    statement, err := db.Prepare(SQL_SELECT_SERIES_BY_NAME);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    var series *model.Series = model.EmptySeries();
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
            &series.Year,
            &series.MetadataSource,
            &series.MetadataSourceID,
    );
    if (err != nil) {
        return nil, err;
    }

    return series, nil;
}

// If the seies does not exist, an error will be returned.
func FetchSeriesByID(id int) (*model.Series, error) {
    statement, err := db.Prepare(SQL_SELECT_SERIES_BY_ID);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    var series *model.Series = model.EmptySeries();

    err = statement.QueryRow(id).Scan(
            &series.ID,
            &series.Author,
            &series.Year,
            &series.MetadataSource,
            &series.MetadataSourceID,
    );
    if (err != nil) {
        return nil, err;
    }

    return series, nil;
}

func insertSeries(series *model.Series) error {
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

func ResolveCrawlRequest(request *model.MetadataCrawlRequest, results []*model.MetadataCrawl) error {
    // TEST
    return nil;
}
