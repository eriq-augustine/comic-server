package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
)

// TODO(eriq): embed/templates.
const CLIENT_DIR = "client";

var routes = []route{
    newRoute("GET", `/client/.*`, handleClient),
    newRoute("GET", `/api/archive/list`, handleArchiveListAll),
    newRoute("GET", `/api/archive/blob/(\d+)`, handleArchiveBlob),
}

type routeHandler func(matches []string, response http.ResponseWriter, request *http.Request) error;

// Inspired by https://benhoyt.com/writings/go-routing/
type route struct {
    method string
    regex *regexp.Regexp
    handler routeHandler
}

func newRoute(method string, pattern string, handler routeHandler) route {
    return route{method, regexp.MustCompile("^" + pattern + "$"), handler};
}

func Serve(response http.ResponseWriter, request *http.Request) {
    log.Debug().
        Str("method", request.Method).
        Str("url", request.URL.Path).
        Msg("");

    for _, route := range routes {
        if (route.method != request.Method) {
            continue;
        }

        matches := route.regex.FindStringSubmatch(request.URL.Path);
        if (matches == nil) {
            continue;
        }

        err := route.handler(matches, response, request);
        if (err != nil) {
            log.Error().Err(err).Str("path", request.URL.Path).Msg("Handler had an error.");
            http.Error(response, "Server Error", http.StatusInternalServerError);
        }

        return;
    }

    http.NotFound(response, request);
}


func StartServer() {
    var port = config.GetInt("server.port");

    log.Info().Msgf("Serving on %d.", port);

    err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(Serve));
    if (err != nil) {
        log.Fatal().Err(err).Msg("Server stopped.");
    }
}

func handleClient(matches []string, response http.ResponseWriter, request *http.Request) error {
    var targetPath = filepath.Join(CLIENT_DIR, strings.TrimPrefix(request.URL.Path, "/client/"));
    http.ServeFile(response, request, targetPath);
    return nil;
}

func handleArchiveListAll(matches []string, response http.ResponseWriter, request *http.Request) error {
    archives, err := database.FetchArchives();
    if (err != nil) {
        return fmt.Errorf("Failed to list archives: %w.", err);
    }

    data, err := json.Marshal(archives);
    if (err != nil) {
        return fmt.Errorf("Failed to serialize archives: %w.", err);
    }

    response.Header().Add("Content-Type", "application/json");
    response.Write(data);
    return nil;
}

func handleArchiveBlob(matches []string, response http.ResponseWriter, request *http.Request) error {
    id, _ := strconv.Atoi(matches[1]);

    archive, err := database.FetchArchiveByID(id);
    if (err != nil) {
        return fmt.Errorf("Failed to fetch archive (%d): %w.", id, err);
    }

    http.ServeFile(response, request, archive.Path);
    return nil;
}
