package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/rs/zerolog/log"

    _ "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
    "github.com/eriq-augustine/comic-server/util"
)

func main() {
    target := parseArgs();

    err := database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    archives, err := metadata.ImportPath(target);
    if (err != nil) {
        log.Fatal().Err(err).Str("path", target).Msg("Failed to import path.");
    }

    fmt.Printf("Found %d archives.\n", len(archives));
    for _, archive := range archives {
        fmt.Println("    ", archive);
    }
}

func parseArgs() string {
    var target = flag.String("target", "", "the file or directory to import");

    flag.Parse();

    if (target == nil || *target == "") {
        fmt.Println("No target specified.");
        flag.Usage();
        os.Exit(1);
    }

    if (!util.PathExists(*target)) {
        fmt.Printf("Target path does does exist: '%s'.\n", *target);
        flag.Usage();
        os.Exit(1);
    }

    return *target;
}
