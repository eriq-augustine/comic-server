package metadata

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/model"
)

type crawler func(string, int, *model.Series) ([]*model.MetadataCrawl, error);

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
    var year = -1;

    match := regexp.MustCompile(`^(.*)\s+\((\d{4})\)\s*$`).FindStringSubmatch(query);
    if (match != nil) {
        query = match[1];
        year, _ = strconv.Atoi(match[2]);
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
            errorCount++;
            log.Error().Err(err).Msg("Failed to resolve crawl requests.");
        }

        err = attemptCrawlMatch(query, year, updatedSeries, results);
        if (err != nil) {
            errorCount++;
            log.Error().Err(err).Msg("Failed to attemp crawl match.");
        }
    }

    if (errorCount > 0) {
        return fmt.Errorf("Encountered %d error while processing crawl request for '%s'.", errorCount, request.Query);
    }

    return nil;
}

// TODO(eriq): Matching is super complex, keep it simple for now.
func attemptCrawlMatch(query string, year int, series *model.Series, crawls []*model.MetadataCrawl) error {
    for _, crawl := range crawls {
        if (crawl.Name != query) {
            continue;
        }

        if ((year > 0) && (crawl.Year != nil) && (*crawl.Year != year)) {
            continue;
        }

        series.AssumeCrawl(crawl);
        
        return database.UpdateSeries(series);
    }

    return nil;
}
