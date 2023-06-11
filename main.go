package main

import (
    "encoding/json"
    "fmt"
    "io/fs"
    "net/http"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/metadata"
)

// TEST
const PORT = 8080;
const DATA_DIR = "test-data";
const CLIENT_DIR = "client";

type Server struct {
    archives []metadata.Archive
}

func (this *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    var path = request.URL.Path;

    // TEST
    fmt.Println("URL: ", path);

    if (strings.HasPrefix(path, "/client/")) {
        var targetPath = filepath.Join(CLIENT_DIR, strings.TrimPrefix(path, "/client/"));

        // TEST
        fmt.Println("TEST", targetPath);

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
    var server = Server{
        archives: importDir(DATA_DIR),
    };

    server.Run();
}

func importDir(rootPath string) []metadata.Archive {
    var archives = make([]metadata.Archive, 0);

    filepath.WalkDir(rootPath, func(path string, dirent fs.DirEntry, err error) error {
        if err != nil {
            log.Fatal().Err(err);
        }

        if (dirent.IsDir()) {
            return nil;
        }

        var archive = importPath(path);
        archive.ID = len(archives);

        archives = append(archives, archive);

        return nil;
    });

    return archives;
}

func importPath(path string) metadata.Archive {
    var archive = metadata.FromPath(path);
    return archive;
}
