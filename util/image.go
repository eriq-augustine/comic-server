package util

import (
    "os"
    "path/filepath"

    "github.com/eriq-augustine/comic-server/config"
)

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
    return filepath.Join(config.GetString("image.dir"), relpath);
}
