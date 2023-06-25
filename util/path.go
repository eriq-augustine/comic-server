package util

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// Tell if a path exists.
func PathExists(path string) bool {
    _, err := os.Stat(path);
    if (err != nil) {
        if os.IsNotExist(err) {
            return false;
        }
    }

    return true;
}

func IsDir(path string) bool {
    if (!PathExists(path)) {
        return false;
    }

    stat, err := os.Stat(path);
    if (err != nil) {
        if os.IsNotExist(err) {
            return false;
        }

        return false;
    }

    return stat.IsDir();
}

// Is |prefix| a prefix of |target|.
// (Same paths counts as a prefix.)
func IsPrefixPath(target string, prefix string) (bool, error) {
    targetPath, err := AbsWithSlash(target);
    if (err != nil) {
        return false, err;
    }

    prefixPath, err := AbsWithSlash(prefix);
    if (err != nil) {
        return false, err;
    }

    return strings.HasPrefix(targetPath, prefixPath), nil;
}

// Call filepath.Abs(), but add a trailing slash for dirs.
// Checking for a dir requires a stat and implies the existance of the dirent.
func AbsWithSlash(path string) (string, error) {
    path, err := filepath.Abs(path);
    if (err != nil) {
        return "", fmt.Errorf("Could not get abs path for '%s': %w.", path, err);
    }

    if (IsDir(path)) {
        path += "/";
    }

    return path, nil;
}
