package web

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/eriq-augustine/comic-server/database"
)

func handleSeriesListAll(matches []string, response http.ResponseWriter, request *http.Request) error {
    allSeries, err := database.FetchSeries();
    if (err != nil) {
        return fmt.Errorf("Failed to list series: %w.", err);
    }

    data, err := json.Marshal(allSeries);
    if (err != nil) {
        return fmt.Errorf("Failed to serialize series: %w.", err);
    }

    response.Header().Add("Content-Type", "application/json");
    response.Write(data);
    return nil;
}

func handleSeries(matches []string, response http.ResponseWriter, request *http.Request) error {
    id, _ := strconv.Atoi(matches[1]);

    series, err := database.FetchSeriesByID(id);
    if (err != nil) {
        return fmt.Errorf("Failed to fetch series (%d): %w.", id, err);
    }

    data, err := json.Marshal(series);
    if (err != nil) {
        return fmt.Errorf("Failed to serialize series (%d): %w.", id, err);
    }

    response.Header().Add("Content-Type", "application/json");
    response.Write(data);
    return nil;
}
