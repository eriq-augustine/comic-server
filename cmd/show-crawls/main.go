package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/rs/zerolog/log"

    _ "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
)

func main() {
    id := parseArgs();

    err := database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    series, err := database.FetchSeriesByID(id);
    if (err != nil) {
        log.Fatal().Err(err).Int("id", id).Msg("Could not fetch series.");
    }

    crawls, err := database.FetchCrawlsBySourceSeries(id);
    if (err != nil) {
        log.Fatal().Err(err).Int("id", id).Msg("Could not fetch crawls for series.");
    }

    fmt.Printf("Series %d -- '%s' (%s)\n", series.ID, series.Name, getYear(series.Year));
    fmt.Printf("Found %d crawls for series:\n", len(crawls));
    for _, crawl := range crawls {
        fmt.Printf("    %d [%s::%s] - '%s' (%s)\n", crawl.ID, crawl.MetadataSource, crawl.MetadataSourceID, crawl.Name, getYear(crawl.Year));
    }
}

func getYear(year *int) string {
    if (year == nil) {
        return "????";
    }

    return fmt.Sprintf("%d", *year);
}

func parseArgs() int {
    var id = flag.Int("series", -1, "the integer id of the series to show the crawls for");

    flag.Parse();

    if (*id <= 0) {
        fmt.Println("Positive id must be specified.");
        flag.Usage();
        os.Exit(1);
    }

    return *id;
}
