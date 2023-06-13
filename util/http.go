package util

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

// TODO(eriq): Config.
const CACHE_DIR = "__cache__";

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

func Get(url string) (string, error) {
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
