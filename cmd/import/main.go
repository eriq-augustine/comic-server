package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/rs/zerolog/log"

    _ "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
    "github.com/eriq-augustine/comic-server/model"
)

func main() {
    targets := parseArgs();

    err := database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    archives := make([]*model.Archive, 0);

    for _, target := range targets {
        newArchives, err := metadata.ImportPath(target);
        if (err != nil) {
            log.Fatal().Err(err).Str("path", target).Msg("Failed to import path.");
        }

        archives = append(archives, newArchives...);
    }

    fmt.Printf("Found %d archives.\n", len(archives));
    for _, archive := range archives {
        fmt.Println("    ", archive);
    }
}

func parseArgs() []string {
    flag.Parse();
    targets := flag.Args();

    if (len(targets) == 0) {
        fmt.Println("No targets specified.");
        flag.Usage();
        os.Exit(1);
    }

    return targets;
}
