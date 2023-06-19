package web

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/eriq-augustine/comic-server/database"
)

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

func handleArchivesBySeries(matches []string, response http.ResponseWriter, request *http.Request) error {
    id, _ := strconv.Atoi(matches[1]);

    archives, err := database.FetchArchivesBySeries(id);
    if (err != nil) {
        return fmt.Errorf("Failed to list archives by series (%d): %w.", id, err);
    }

    data, err := json.Marshal(archives);
    if (err != nil) {
        return fmt.Errorf("Failed to serialize archives: %w.", err);
    }

    response.Header().Add("Content-Type", "application/json");
    response.Write(data);
    return nil;
}
