package web

import (
    "fmt"
    "net/http"
    "regexp"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
)

var routes = []route{
    newRedirect("GET", `/static/index.html`, `/static/series/list/index.html`),
    newRedirect("GET", `/static/series/index.html`, `/static/series/list/index.html`),

    newRoute("GET", `/api/archive/(\d+)`, handleArchive),
    newRoute("GET", `/api/archive/blob/(\d+)`, handleArchiveBlob),
    newRoute("GET", `/api/archive/list`, handleArchiveListAll),
    newRoute("GET", `/api/archive/series/(\d+)`, handleArchivesBySeries),
    newRoute("GET", `/api/image/blob/(.*)`, handleImageBlob),
    newRoute("GET", `/api/series/(\d+)`, handleSeries),
    newRoute("GET", `/api/series/list`, handleSeriesListAll),
    newRoute("GET", `/static`, handleStatic),
    newRoute("GET", `/static/.*`, handleStatic),
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

func newRedirect(method string, pattern string, target string) route {
    redirectFunc := func(matches []string, response http.ResponseWriter, request *http.Request) error {
        return handleRedirect(target, matches, response, request);
    };
    return route{method, regexp.MustCompile("^" + pattern + "$"), redirectFunc};
}

func handleRedirect(target string, matches []string, response http.ResponseWriter, request *http.Request) error {
    http.Redirect(response, request, target, 301);
    return nil;
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
