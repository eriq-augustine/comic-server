package util

import (
    "archive/zip"
    "fmt"
    "mime"
    "path/filepath"
    "strings"
)

func ZipListFilenames(path string) ([]string, error) {
    reader, err := zip.OpenReader(path);
    if (err != nil) {
        return nil, fmt.Errorf("Could not open zip file (%s): %w.", path, err);
    }
    defer reader.Close()

    var filenames = make([]string, 0);
    for _, file := range reader.File {
        filenames = append(filenames, file.Name);
    }

    return filenames, nil;
}

// Count the number of image files in a zip archive.
func ZipImageCount(path string) (int, error) {
    filenames, err := ZipListFilenames(path);
    if (err != nil) {
        return 0, err;
    }

    count := 0;

    for _, filename := range filenames {
        mimeStr := mime.TypeByExtension(filepath.Ext(filename));
        if (strings.HasPrefix(mimeStr, "image/")) {
            count++;
        }
    }

    return count, nil;
}
