package api

import (
    "net/http"
    "path/filepath"
    "strings"
)

// TODO(eriq): embed/templates.
const CLIENT_DIR = "client";

func handleClient(matches []string, response http.ResponseWriter, request *http.Request) error {
    var targetPath = filepath.Join(CLIENT_DIR, strings.TrimPrefix(request.URL.Path, "/client/"));
    http.ServeFile(response, request, targetPath);
    return nil;
}
