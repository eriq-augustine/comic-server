package main

import (
    "github.com/alecthomas/kong"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
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

    // TODO(eriq): Setup as background job.
    err = metadata.ProcessCrawlRequests();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Failed to crawl.");
    }
}
