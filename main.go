package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/metadata"
    "github.com/eriq-augustine/comic-server/types"
)

// TEST
const PORT = 8080;
const DATA_DIR = "test-data";
const CLIENT_DIR = "client";

type Server struct {
    archives []*types.Archive
}

func (this *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    var path = request.URL.Path;

    log.Debug().Str("URL",  path).Msg("");

    if (strings.HasPrefix(path, "/client/")) {
        var targetPath = filepath.Join(CLIENT_DIR, strings.TrimPrefix(path, "/client/"));
        http.ServeFile(response, request, targetPath);
        return;
    }

    // TEST: Make actual router.

    if (strings.HasPrefix(path, "/api/list")) {
        data, err := json.Marshal(this.archives);
        if (err != nil) {
            log.Fatal().Err(err);
        }

        response.Header().Add("Content-Type", "application/json");
        response.Write(data);
        return;
    }

    if (strings.HasPrefix(path, "/blob/archive/")) {
        var pattern = regexp.MustCompile(`^/blob/archive/(\d+)$`);
        match := pattern.FindStringSubmatch(path);
        if (match == nil) {
            log.Fatal().Msg("Could not match archive regex.");
        }

        index, _ := strconv.Atoi(match[1]);
        http.ServeFile(response, request, this.archives[index].Path);
        return;
    }
}

func (this *Server) Run() {
    log.Info().Msgf("Serving on %d.", PORT);
    log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", PORT), this));
}

func main() {
    // TODO(eriq): Config
    zerolog.SetGlobalLevel(zerolog.DebugLevel);

    err := database.Open();
    if (err != nil) {
        log.Fatal().Err(err).Msg("Could not open database.");
    }
    defer database.Close();

    archives, err := metadata.ImportDir(DATA_DIR);
    if (err != nil) {
        log.Fatal().Err(err).Msg("Failed to import dir.");
    }

    // TEST: TODO: Get archives from DB.
    var server = Server{
        archives: archives,
    };

    server.Run();
}
