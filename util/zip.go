package util

import (
    "archive/zip"
    "fmt"
    "io"
    "mime"
    "path/filepath"
    "sort"
    "strings"

    "github.com/google/uuid"
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

// Assume that the archive is a CBZ, and get the number of imaeges and the first image (as a thumbnail).
func CBZInfo(path string) (int, string, error) {
    filenames, err := ZipListFilenames(path);
    if (err != nil) {
        return 0, "", err;
    }

    sort.Strings(filenames);

    var count int = 0;
    var imagePath string;

    for _, filename := range filenames {
        mimeStr := mime.TypeByExtension(filepath.Ext(filename));
        if (strings.HasPrefix(mimeStr, "image/")) {
            if (count == 0) {
                imagePath, err = extractImage(path, filename);
                if (err != nil) {
                    return 0, "", fmt.Errorf("Could not extract image file (%s) from archive (%s): %w.", filename, path, err);
                }
            }

            count++;
        }
    }

    return count, imagePath, nil;
}

func ZipExtractFile(zipPath string, filename string) ([]byte, error) {
    reader, err := zip.OpenReader(zipPath);
    if (err != nil) {
        return nil, fmt.Errorf("Could not open zip file (%s): %w.", zipPath, err);
    }
    defer reader.Close()

    file, err := reader.Open(filename);
    if (err != nil) {
        return nil, fmt.Errorf("Could not open internal zip file (%s) from archive (%s): %w.", filename, zipPath, err);
    }
    defer file.Close();

    data, err := io.ReadAll(file);
    if (err != nil) {
        return nil, fmt.Errorf("Could not read internal zip file (%s) from archive (%s): %w.", filename, zipPath, err);
    }

    return data, nil;
}

func extractImage(zipPath string, filename string) (string, error) {
    var relpath string;
    for {
        id := uuid.New().String();
        relpath = filepath.Join("archive", id, "image" + filepath.Ext(filename));
        // HACK(eriq): There is a race condition here.
        if (!PathExists(GetImagePath(relpath))) {
            break;
        }
    }

    data, err := ZipExtractFile(zipPath, filename);
    if (err != nil) {
        return "", err;
    }

    err = WriteImage(data, relpath);
    if (err != nil) {
        return "", fmt.Errorf("Could not write image data (%s): %w.", relpath, err);
    }

    return relpath, nil;
}
