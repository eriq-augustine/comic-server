/*
TEST

Metadata is layered:
 - User
 - Archive (ComicInfo.xml)
 - Path
 - Web

Users can choose which source to prefer for a piece of infomation.
*/
package metadata

import (
    "path/filepath"
    "regexp"
)

// Archives are things the relate to physical files/directories on disk.
// The are (or can be) packaged up in a single CBZ file.
type Archive struct {
    ID int
    Path string
    Filename string
    Series string
    Volume string
    Chapter string
}

// Try to re-create metadata using only path information.
func FromPath(path string) Archive {
    var filename = filepath.Base(path);
    var archive = Archive{Path: path, Filename: filename};

    var pattern = regexp.MustCompile(`^(.*)\s+v(\d+[a-z]?)\s+c(\d+[a-z]?)\.cbz$`);
    match := pattern.FindStringSubmatch(filename);
    if (match != nil) {
        archive.Series = match[1];
        archive.Volume = match[2];
        archive.Chapter = match[3];
    }

    return archive;
}
