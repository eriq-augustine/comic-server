package util

import (
    "os"
    "path/filepath"
)

// TODO(eriq): Config
const IMAGE_DIR = "__images__"

func FetchImage(url string) (string, error) {
    var path = getImagePath(url);
    if (PathExists(path)) {
        return path, nil;
    }

    data, err := GetWithCache(url);
    if (err != nil) {
        return "", err;
    }

    os.MkdirAll(filepath.Dir(path), 0755);
    err = os.WriteFile(path, data, 0644);
    return path, err;
}

func getImagePath(url string) string {
    return filepath.Join(IMAGE_DIR, url);
}
