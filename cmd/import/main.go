package main

import (
    "fmt"

    "github.com/alecthomas/kong"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
)

var args struct {
    config.ConfigArgs
    Path []string `help:"Target paths to import." arg:"" type:"path"`
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

    numNewArchives := 0;
    numExistingArchives := 0;

    for _, target := range args.Path {
        newArchives, err := metadata.ImportPath(target);
        if (err != nil) {
            log.Fatal().Err(err).Str("path", target).Msg("Failed to import path.");
        }

        for _, archive := range newArchives {
            if (archive.New) {
                numNewArchives++;
            } else {
                numExistingArchives++;
            }
        }
    }

    fmt.Printf("Found %d new and %d existing archives.\n", numNewArchives, numExistingArchives);
}
