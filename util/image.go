package util

import (
    "os"
    "path/filepath"
)

// TODO(eriq): Config
const IMAGE_DIR = "__images__"

// Returns (abspath, relpath, error).
func FetchImage(url string) (string, string, error) {
    var relpath = filepath.Clean(url);

    path, err := filepath.Abs(getImagePath(relpath));
    if (err != nil) {
        return "", "", err;
    }

    if (PathExists(path)) {
        return path, relpath, nil;
    }

    data, err := GetWithCache(url);
    if (err != nil) {
        return "", "", err;
    }

    os.MkdirAll(filepath.Dir(path), 0755);
    err = os.WriteFile(path, data, 0644);

    return path, relpath, err;
}

func CopyImage(relSource string, relDest string) (string, error) {
    source := getImagePath(relSource);
    dest := getImagePath(relDest);

    err := CopyFile(source, dest);
    return dest, err;
}

func getImagePath(relpath string) string {
    return filepath.Join(IMAGE_DIR, relpath);
}
