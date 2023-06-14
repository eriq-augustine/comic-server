package metadata

import (
    "fmt"
    "regexp"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/model"
)

type crawler func(string, string, *model.Series) ([]*model.MetadataCrawl, error);

var metadataSources = make(map[string]crawler);

func ProcessCrawlRequests() error {
    requests, err := database.FetchCrawlRequests();
    if (err != nil) {
        return err;
    }

    var errorCount = 0;

    for _, request := range requests {
        err = ProcessCrawlRequest(request);
        if (err != nil) {
            errorCount++;
            log.Error().Err(err).Str("Query", request.Query).Msg("Error crawling.");
        }
    }

    if (errorCount > 0) {
        return fmt.Errorf("Encountered %d error while processing crawl requests.", errorCount);
    }

    return nil;
}

func ProcessCrawlRequest(request *model.MetadataCrawlRequest) error {
    // Get the most up-to-date series.
    updatedSeries, err := database.FetchSeriesByID(request.Series.ID);

    if (updatedSeries.MetadataSource != nil) {
        // The series already has metadata, don't bother with this crawl.
        return nil;
    }

    var query = request.Query;
    var year = "";

    match := regexp.MustCompile(`^(.*)\s+\((\d{4})\)\s*$`).FindStringSubmatch(query);
    if (match != nil) {
        query = match[1];
        year = match[2];
    }

    var errorCount = 0;
    var results = make([]*model.MetadataCrawl, 0);
    
    for source, crawlFunction := range metadataSources {
        result, err := crawlFunction(query, year, updatedSeries);
        if (err != nil) {
            errorCount++;
            log.Warn().Err(err).Str("source", source).Msg("");
        }

        if (result != nil) {
            results = append(results, result...);
        }
    }

    if (len(results) > 0) {
        err = database.ResolveCrawlRequest(request, results);
        if (err != nil) {
            return err;
        }
    }

    if (errorCount > 0) {
        return fmt.Errorf("Encountered %d error while processing crawl request for '%s'.", errorCount, request.Query);
    }

    return nil;
}
