package util

import (
    "os"
    "path/filepath"

    "github.com/eriq-augustine/comic-server/config"
)

// Returns (abspath, relpath, error).
func FetchImage(url string) (string, string, error) {
    var relpath = filepath.Clean(url);

    path, err := filepath.Abs(GetImagePath(relpath));
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

func WriteImage(data []byte, relpath string) (error) {
    path := GetImagePath(relpath);
    os.MkdirAll(filepath.Dir(path), 0755);
    return os.WriteFile(path, data, 0644);
}

func CopyImage(relSource string, relDest string) (string, error) {
    source := GetImagePath(relSource);
    dest := GetImagePath(relDest);

    err := CopyFile(source, dest);
    return dest, err;
}

func GetImagePath(relpath string) string {
    return filepath.Join(config.GetString("image.dir"), relpath);
}
