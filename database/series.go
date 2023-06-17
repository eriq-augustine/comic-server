package database

import (
    _ "embed"
    "fmt"

    "github.com/eriq-augustine/comic-server/model"
)

//go:embed sql/insert-series.sql
var SQL_INSERT_SERIES string;

//go:embed sql/select-series.sql
var SQL_SELECT_SERIES string;

//go:embed sql/select-series-by-id.sql
var SQL_SELECT_SERIES_BY_ID string;

//go:embed sql/select-series-by-name.sql
var SQL_SELECT_SERIES_BY_NAME string;

//go:embed sql/update-series.sql
var SQL_UPDATE_SERIES string;

func FetchSeries() ([]*model.Series, error) {
    rows, err := db.Query(SQL_SELECT_SERIES);
    if (err != nil) {
        return nil, err;
    }
    defer rows.Close();

    var allSeries = make([]*model.Series, 0);

    for (rows.Next()) {
        series, err := scanSeries(rows);
        if (err != nil) {
            return nil, err;
        }

        allSeries = append(allSeries, series);
    }

    return allSeries, nil;
}

func FetchSeriesByName(name string) (*model.Series, error) {
    statement, err := db.Prepare(SQL_SELECT_SERIES_BY_NAME);
    if (err != nil) {
        return nil, err;
    }
    defer statement.Close();

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

    series, err := scanSeries(rows);
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

    series, err := scanSeries(statement.QueryRow(id));
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

func UpdateSeries(series *model.Series) error {
    statement, err := db.Prepare(SQL_UPDATE_SERIES);
    if (err != nil) {
        return fmt.Errorf("Failed to prepare series update (%s): %w.", SQL_UPDATE_SERIES, err);
    }
    defer statement.Close();

    _, err = statement.Exec(
        series.Name,
        series.Author,
        series.Year,
        series.URL,
        series.Description,
        series.CoverImageRelPath,
        series.MetadataSource,
        series.MetadataSourceID,
        series.ID,
    );

    return err;
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

func scanSeries(scanner RowScanner) (*model.Series, error) {
    var series = model.Series{};

    err := scanner.Scan(
            &series.ID,
            &series.Name,
            &series.Author,
            &series.Year,
            &series.URL,
            &series.Description,
            &series.CoverImageRelPath,
            &series.MetadataSource,
            &series.MetadataSourceID,
    );

    return &series, err;
}
