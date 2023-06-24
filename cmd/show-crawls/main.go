package main

import (
    "fmt"

    "github.com/alecthomas/kong"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
)

var args struct {
    config.ConfigArgs
    ID int `help:"The integer id of the series to show the crawls for" arg:""`
}

func main() {
    kong.Parse(&args);
    err := config.HandleConfigArgs(args.ConfigArgs);
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not load config options.");
    }

    err = database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    series, err := database.FetchSeriesByID(args.ID);
    if (err != nil) {
        log.Fatal().Err(err).Int("id", args.ID).Msg("Could not fetch series.");
    }

    crawls, err := database.FetchCrawlsBySourceSeries(args.ID);
    if (err != nil) {
        log.Fatal().Err(err).Int("id", args.ID).Msg("Could not fetch crawls for series.");
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
