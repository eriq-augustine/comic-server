package metadata

import (
    "errors"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "github.com/rs/zerolog/log"

    "github.com/eriq-augustine/comic-server/config"
    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/model"
    "github.com/eriq-augustine/comic-server/util"
)

const IMPORTED_ARCHIVES_DIR = "__imported_archives__";

type ImportedArchive struct {
    Archive *model.Archive
    New bool
}

func ImportPath(path string) ([]*ImportedArchive, error) {
    if (util.IsDir(path)) {
        return ImportDir(path);
    }

    archive, err := ImportFile(path);
    if (err != nil) {
        return nil, err;
    }

    result := make([]*ImportedArchive, 0, 1);

    if (archive == nil) {
        return result, nil;
    }

    return append(result, archive), nil;
}

// May return (nil, nil) if the path is not an archive.
func ImportFile(path string) (*ImportedArchive, error) {
    log.Trace().Str("path", path).Msg("Importing file.");

    rawArchive, err := fromPath(path);
    if (err != nil) {
        return nil, fmt.Errorf("Failed to import file (%s): %w.", path, err);
    }

    if (rawArchive == nil) {
        return nil, nil;
    }

    exists, err := database.PersistArchive(rawArchive);
    if (err != nil) {
        return nil, fmt.Errorf("Failed to persist imported file (%s): %w.", path, err);
    }

    log.Debug().Str("path", path).Str("name", rawArchive.Series.Name).Bool("exists", exists).Msg("Imported archive.");

    var archive = ImportedArchive{Archive: rawArchive, New: !exists};
    return &archive, nil;
}

// Recursively import archive from a directory.
// Archives will be added as they are collected.
// Any errors that occur at the archive level will be joined in the return error and a nil will be placed in the retrned slice.
func ImportDir(rootPath string) ([]*ImportedArchive, error) {
    log.Debug().Str("path", rootPath).Msg("Importing Dir");

    var archives = make([]*ImportedArchive, 0);
    var allErrors error = nil;

    err := filepath.WalkDir(rootPath, func(path string, dirent fs.DirEntry, err error) error {
        if (err != nil) {
            return fmt.Errorf("Failed to walk '%v': %w.", path, err);
        }

        if (dirent.IsDir()) {
            return nil;
        }

        archive, err := ImportFile(path);

        archives = append(archives, archive);
        allErrors = errors.Join(allErrors, err);

        return nil;
    });

    allErrors = errors.Join(allErrors, err);

    return archives, allErrors;
}

// Try to re-create metadata using only path information.
func fromPath(path string) (*model.Archive, error) {
    abspath, err := filepath.Abs(path);
    if (err != nil) {
        return nil, fmt.Errorf("Could not form abs path from '%s': %w.", path, err);
    }

    relpath, requiresCopy, err := resolveArchivePath(abspath);
    if (err != nil) {
        return nil, err;
    }

    var filename = filepath.Base(path);
    var archive = model.EmptyArchive(relpath);

    var pattern = regexp.MustCompile(`^(.*)\s+v(\d+[a-z]?)\s*c(\d+[a-z]?)\.(?i:cbz|zip)$`);
    match := pattern.FindStringSubmatch(filename);
    if (match != nil) {
        archive.Series.Name = match[1];
        archive.Volume = &match[2];
        archive.Chapter = &match[3];
    } else {
        archive.Series.Name = filename;
    }

    ext := filepath.Ext(strings.ToLower(filename));
    switch ext {
    case ".zip":
        pageCount, err := util.ZipImageCount(abspath);
        if (err != nil) {
            log.Warn().Err(err).Str("path", abspath).Msg("Failed to get zip page count.");
        } else {
            archive.PageCount = &pageCount;
        }
    case ".cbz":
        pageCount, coverImagePath, err := util.CBZInfo(abspath);
        if (err != nil) {
            log.Warn().Err(err).Str("path", abspath).Msg("Failed to get CBZ image info.");
        } else {
            archive.PageCount = &pageCount;
            archive.CoverImageRelPath = &coverImagePath;
        }
    default:
        log.Debug().Str("extension", ext).Msg("Unrecognized archive extension, can only use cbz and zip. Skipping");
        return nil, nil;
    }

    if (requiresCopy) {
        destPath := filepath.Join(config.GetString("paths.archives"), relpath);
        err = util.CopyFile(abspath, destPath);
        if (err != nil) {
            return nil, err;
        }
    }

    return archive, nil;
}

// Create a relative path that represents |path| relative to the archives directory.
// If the path is outside of the archives directory, then the boolean will be true.
// Returns: (relpath, requires copy, error).
func resolveArchivePath(path string) (string, bool, error) {
    archivesPath, err := util.AbsWithSlash(config.GetString("paths.archives"));
    if (err != nil) {
        return "", false, err;
    }

    path, err = util.AbsWithSlash(path);
    if (err != nil) {
        return "", false, err;
    }

    isPrefix, err := util.IsPrefixPath(path, archivesPath);
    if (err != nil) {
        return "", false, err;
    }

    if (isPrefix) {
        relpath, err := filepath.Rel(archivesPath, path);
        if (err != nil) {
            return "", false, err;
        }

        return relpath, false, nil;
    }

    err = os.MkdirAll(filepath.Join(archivesPath, IMPORTED_ARCHIVES_DIR), 0775);
    if (err != nil) {
        return "", false, err;
    }

    relpath := filepath.Join(IMPORTED_ARCHIVES_DIR, filepath.Base(path));
    return relpath, true, nil;
}
