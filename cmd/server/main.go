package main

import (
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
    "github.com/eriq-augustine/comic-server/api"
)

// TODO(eriq): Handle import in a specific bin.
const DATA_DIR = "test-data";

func main() {
    err := database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    _, err = metadata.ImportDir(DATA_DIR);
    if (err != nil) {
        log.Fatal().Err(err).Msg("Failed to import dir.");
    }

    // TODO(eriq): Setup as background job.
    err = metadata.ProcessCrawlRequests();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Failed to crawl.");
    }

    api.StartServer();
}
