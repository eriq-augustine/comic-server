package main

import (
    "fmt"
    "os"

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
    numErrorArchives := 0;

    for _, target := range args.Path {
        newArchives, err := metadata.ImportPath(target);
        if (err != nil) {
            log.Error().Err(err).Str("path", target).Msg("Failed to import all archives on path.");
            continue;
        }

        for _, archive := range newArchives {
            if (archive == nil) {
                numErrorArchives++;
            } else if (archive.New) {
                numNewArchives++;
            } else {
                numExistingArchives++;
            }
        }
    }

    fmt.Printf("Encountered %d new, %d existing, and %d erroneous archives.\n", numNewArchives, numExistingArchives, numErrorArchives);

    if (numErrorArchives > 0) {
        os.Exit(1);
    }

    os.Exit(0);
}
