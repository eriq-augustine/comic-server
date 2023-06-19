package database

import (
    _ "embed"
    "fmt"

    "github.com/eriq-augustine/comic-server/model"
)

//go:embed sql/delete-crawl-request.sql
var SQL_DELETE_CRAWL_REQUEST string;

//go:embed sql/insert-crawl.sql
var SQL_INSERT_CRAWL string;

//go:embed sql/select-crawls-by-series.sql
var SQL_SELECT_CRAWLS_BY_SERIES string;

//go:embed sql/select-crawl-requests.sql
var SQL_SELECT_CRAWL_REQUESTS string;

//go:embed sql/upsert-crawl-request.sql
var SQL_UPSERT_CRAWL_REQUEST string;

func FetchCrawlsBySourceSeries(seriesID int) ([]*model.MetadataCrawl, error) {
    statement, err := db.Prepare(SQL_SELECT_CRAWLS_BY_SERIES);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

    rows, err := statement.Query(seriesID);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    crawls := make([]*model.MetadataCrawl, 0);

    for (rows.Next()) {
        crawl, _, err := scanCrawlNoSeries(rows);
        if (err != nil) {
            return nil, err;
        }

        crawls = append(crawls, crawl);
    }

    return crawls, nil;
}

func FetchCrawlRequests() ([]*model.MetadataCrawlRequest, error) {
    rows, err := db.Query(SQL_SELECT_CRAWL_REQUESTS);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    var requests = make([]*model.MetadataCrawlRequest, 0);

    for (rows.Next()) {
        var request = model.EmptyCrawlRequest();

        err = rows.Scan(
                &request.ID,
                &request.Query,
                &request.Timestamp,
                &request.Series.ID,
                &request.Series.Name,
                &request.Series.Author,
                &request.Series.Year,
                &request.Series.URL,
                &request.Series.Description,
                &request.Series.CoverImageRelPath,
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

func RequestMetadataCrawl(series *model.Series) error {
    statement, err := db.Prepare(SQL_UPSERT_CRAWL_REQUEST);
    if (err != nil) {
        return err;
    }
    defer statement.Close();

    _, err = statement.Exec(series.ID, series.Name);
    return err;
}

func ResolveCrawlRequest(request *model.MetadataCrawlRequest, crawls []*model.MetadataCrawl) error {
    err := deleteCrawlRequest(request);
    if (err != nil) {
        return fmt.Errorf("Failed for delete crawl request: %w", err);
    }

    err = insertCrawls(crawls);
    if (err != nil) {
        return fmt.Errorf("Failed for insert crawl results: %w", err);
    }

    return nil;
}

func deleteCrawlRequest(request *model.MetadataCrawlRequest) error {
    statement, err := db.Prepare(SQL_DELETE_CRAWL_REQUEST);
    if (err != nil) {
        return fmt.Errorf("Failed to prepare crawl request delete (%s): %w.", SQL_DELETE_CRAWL_REQUEST, err);
    }
    defer statement.Close();

    _, err = statement.Exec(request.ID);
    return err;
}

func insertCrawls(crawls []*model.MetadataCrawl) error {
    transaction, err := db.Begin();
    if (err != nil) {
        return err;
    }
    defer transaction.Rollback();

    statement, err := db.Prepare(SQL_INSERT_CRAWL);
    if (err != nil) {
        return fmt.Errorf("Failed to prepare crawl insert (%s): %w.", SQL_INSERT_CRAWL, err);
    }
    defer statement.Close();

    for _, crawl := range crawls {
        _, err = statement.Exec(
            crawl.MetadataSource,
            crawl.MetadataSourceID,
            crawl.SourceSeries.ID,
            crawl.Name,
            crawl.Author,
            crawl.Year,
            crawl.URL,
            crawl.Description,
            crawl.CoverImageRelPath,
        );

        if (err != nil) {
            return err;
        }
    }

    transaction.Commit();
    return nil;
}

func scanCrawlNoSeries(scanner RowScanner) (*model.MetadataCrawl, int, error) {
    var crawl = model.MetadataCrawl{};
    var seriesID int;

    err := scanner.Scan(
            &crawl.ID,
            &crawl.MetadataSource,
            &crawl.MetadataSourceID,
            &seriesID,
            &crawl.Name,
            &crawl.Author,
            &crawl.Year,
            &crawl.URL,
            &crawl.Description,
            &crawl.CoverImageRelPath,
            &crawl.Timestamp,
    );

    return &crawl, seriesID, err;
}
