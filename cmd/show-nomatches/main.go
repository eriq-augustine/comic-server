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

    allSeries, err := database.FetchUnmatchedSeries();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not fetch unmatched series.");
    }

    fmt.Printf("Found %d unmatched series:\n", len(allSeries));
    for _, series := range allSeries {
        fmt.Printf("    %d - '%s'\n", series.ID, series.Name);
    }
}
