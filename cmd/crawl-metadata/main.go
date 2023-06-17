package main

import (
    "github.com/rs/zerolog/log"

    _ "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
)

func main() {
    err := database.Open();
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
