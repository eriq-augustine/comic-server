package metadata

import (
    "fmt"
    "strings"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/model"
    "github.com/eriq-augustine/comic-server/util"
)

func attemptCrawlMatch(query string, year int, series *model.Series, crawls []*model.MetadataCrawl) error {
    matches := make([]*model.MetadataCrawl, 0);

    for _, crawl := range crawls {
        // If both have yer available, use that to disqualify right away.
        if ((year > 0) && (crawl.Year != nil) && (*crawl.Year != year)) {
            continue;
        }

        // Look for a name match.
        if (!checkNameMatch(query, series, crawl)) {
            continue;
        }

        matches = append(matches, crawl);
    }

    if (len(matches) == 0) {
        return nil;
    }

    if (len(matches) > 1) {
        log.Debug().Str("query", query).Int("year", year).Int("series", series.ID).Int("count", len(matches)).Msg("Found multiple possible metadata matches for a series.");
        for _, match := range matches {
            log.Debug().Interface("match", match).Msg("Possible match.");
        }
        return nil;
    }

    crawl := matches[0];
    err := series.AssumeCrawl(crawl);
    if (err != nil) {
        return fmt.Errorf("Failed to assume crawl (%d) for series (%d): %w.", crawl.ID, series.ID, err);
    }

    return database.UpdateSeries(series);
}

func checkNameMatch(query string, series *model.Series, crawl *model.MetadataCrawl) bool {
    names := make([]string, 0);
    names = append(names, crawl.Name);

    if (crawl.AltNames != nil) {
        names = append(names, util.UnsafeSplit(*crawl.AltNames)...);
    }

    for _, name := range names {
        if (strings.EqualFold(query, name)) {
            return true;
        }
    }

    return false;
}
