package util

import (
    "fmt"
    "io"
    "net/http"
    neturl "net/url"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/rs/zerolog/log"
)

// TODO(eriq): Config.
const CACHE_DIR = "__cache__";

// Do not query a single site too quickly.
const RATE_LIMIT_DELAY_SEC = 1
var rateLimit = make(map[string]time.Time);

func GetWithCache(url string) (string, error) {
    text, err := checkCache(url);
    if (err != nil) {
        return "", err;
    }

    if (text != "") {
        return text, nil;
    }

    text, err = Get(url);
    if (err != nil) {
        return "", err;
    }

    err = saveCache(url, text);
    if (err != nil) {
        return "", err;
    }

    return text, nil;
}

func ensureRateLimit(rawURL string) {
    now := time.Now();

    url, err := neturl.Parse(rawURL);
    if (err != nil) {
        log.Warn().Err(err).Str("url", rawURL).Msg("Failed to parse url.");
        return;
    }

    hostname := url.Hostname();

    _, exists := rateLimit[hostname];
    if (!exists) {
        rateLimit[hostname] = now;
    }

    lastAccess := rateLimit[hostname];
    delta := now.Sub(lastAccess);

    // Zero or negative sleeps return right away.
    sleepTime := (time.Second * RATE_LIMIT_DELAY_SEC) - delta;
    time.Sleep(sleepTime);

    rateLimit[hostname] = time.Now();
}

func Get(url string) (string, error) {
    ensureRateLimit(url);

    log.Debug().Str("url", url).Msg("GET");

    response, err := http.Get(url);
    if (err != nil) {
        return "", err;
    }
    defer response.Body.Close()

    if (response.StatusCode != 200) {
        return "", fmt.Errorf("Got non-200 status code (%d) for '%s'.", response.StatusCode, url);
    }

    buffer := new(strings.Builder);
	_, err = io.Copy(buffer, response.Body);
	if (err != nil) {
        return "", err;
	}

	return buffer.String(), nil;
}

func checkCache(url string) (string, error) {
    var path = getCachePath(url);
    if (!PathExists(path)) {
        return "", nil;
    }

    bytes, err := os.ReadFile(path);
    if (err != nil) {
        return "", err;
    }

    log.Debug().Str("url", url).Msg("Cache hit.");
    return string(bytes), nil;
}

func saveCache(url string, text string) error {
    var path = getCachePath(url);
    os.MkdirAll(filepath.Dir(path), 0755);

    if (!PathExists(path)) {
        os.Remove(path);
    }

    return os.WriteFile(path, []byte(text), 0644);
}

func getCachePath(url string) string {
    return filepath.Join(CACHE_DIR, url);
}
