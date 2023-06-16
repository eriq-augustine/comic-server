package util

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    neturl "net/url"
    "os"
    "path/filepath"
    "time"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
)

// Do not query a single site too quickly.
const RATE_LIMIT_DELAY_SEC = 1
var rateLimit = make(map[string]time.Time);

func GetWithCache(url string) ([]byte, error) {
    data, err := checkCache(url);
    if (err != nil) {
        return nil, err;
    }

    if (data != nil) {
        return data, nil;
    }

    data, err = Get(url);
    if (err != nil) {
        return nil, err;
    }

    err = saveCache(url, data);
    if (err != nil) {
        return nil, err;
    }

    return data, nil;
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

func Get(url string) ([]byte, error) {
    ensureRateLimit(url);

    log.Debug().Str("url", url).Msg("GET");

    response, err := http.Get(url);
    if (err != nil) {
        return nil, err;
    }
    defer response.Body.Close()

    if (response.StatusCode != 200) {
        return nil, fmt.Errorf("Got non-200 status code (%d) for '%s'.", response.StatusCode, url);
    }

    bytes := new(bytes.Buffer);
	_, err = io.Copy(bytes, response.Body);
	if (err != nil) {
        return nil, err;
	}

	return bytes.Bytes(), nil;
}

func checkCache(url string) ([]byte, error) {
    var path = getCachePath(url);
    if (!PathExists(path)) {
        return nil, nil;
    }

    bytes, err := os.ReadFile(path);
    if (err != nil) {
        return nil, err;
    }

    log.Debug().Str("url", url).Msg("Cache hit.");
    return bytes, nil;
}

func saveCache(url string, data []byte) error {
    var path = getCachePath(url);
    os.MkdirAll(filepath.Dir(path), 0755);

    if (!PathExists(path)) {
        os.Remove(path);
    }

    return os.WriteFile(path, data, 0644);
}

func getCachePath(url string) string {
    return filepath.Join(config.GetString("cache.dir"), url);
}
