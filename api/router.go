package api

import (
    "fmt"
    "net/http"
    "regexp"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
)

var routes = []route{
    newRoute("GET", `/api/archive/list`, handleArchiveListAll),
    newRoute("GET", `/api/archive/blob/(\d+)`, handleArchiveBlob),
    newRoute("GET", `/api/image/blob/(.*)`, handleImageBlob),
    newRoute("GET", `/api/series/(\d+)`, handleSeries),
    newRoute("GET", `/api/series/list`, handleSeriesListAll),
    newRoute("GET", `/client`, handleClient),
    newRoute("GET", `/client/.*`, handleClient),
}

type routeHandler func(matches []string, response http.ResponseWriter, request *http.Request) error;

// Inspired by https://benhoyt.com/writings/go-routing/
type route struct {
    method string
    regex *regexp.Regexp
    handler routeHandler
}

func StartServer() {
    var port = config.GetInt("server.port");

    log.Info().Msgf("Serving on %d.", port);

    err := http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(serve));
    if (err != nil) {
        log.Fatal().Err(err).Msg("Server stopped.");
    }
}

func newRoute(method string, pattern string, handler routeHandler) route {
    return route{method, regexp.MustCompile("^" + pattern + "$"), handler};
}

func serve(response http.ResponseWriter, request *http.Request) {
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
