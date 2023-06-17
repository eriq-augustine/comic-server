package metadata

import (
    "fmt"
    "io/fs"
    "path/filepath"
    "regexp"

    "github.com/eriq-augustine/comic-server/database"
    "github.com/eriq-augustine/comic-server/model"
    "github.com/eriq-augustine/comic-server/util"
)

func ImportPath(path string) ([]*model.Archive, error) {
    if (util.IsDir(path)) {
        return ImportDir(path);
    }

    archive, err := ImportFile(path);
    if (err != nil) {
        return nil, err;
    }

    return []*model.Archive{archive}, nil;
}

func ImportFile(path string) (*model.Archive, error) {
    return fromPath(path);
}

// Recursively import archive from a directory.
// First the directory will be walked and all the archives collected.
// Then, they will be added to the database (if no error has occured).
// On error, no archives will be added to the database.
func ImportDir(rootPath string) ([]*model.Archive, error) {
    var archives = make([]*model.Archive, 0);

    err := filepath.WalkDir(rootPath, func(path string, dirent fs.DirEntry, err error) error {
        if (err != nil) {
            return fmt.Errorf("Failed to walk '%v': %w.", path, err);
        }

        if (dirent.IsDir()) {
            return nil;
        }

        archive, err := ImportFile(path);
        if (err != nil) {
            return err;
        }

        archives = append(archives, archive);
        return nil;
    });

    if (err != nil) {
        return nil, err;
    }

    err = database.PersistArchives(archives);
    if (err != nil) {
        return nil, err;
    }

    return archives, nil;
}

// Try to re-create metadata using only path information.
func fromPath(path string) (*model.Archive, error) {
    var filename = filepath.Base(path);

    var archive = model.EmptyArchive(path);

    var pattern = regexp.MustCompile(`^(.*)\s+v(\d+[a-z]?)\s+c(\d+[a-z]?)\.cbz$`);
    match := pattern.FindStringSubmatch(filename);
    if (match != nil) {
        archive.Series.Name = match[1];
        archive.Volume = &match[2];
        archive.Chapter = &match[3];
    }

    // TODO(eriq): Page Count

    return archive, nil;
}
