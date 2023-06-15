package model

import (
    "fmt"
    "path/filepath"
    "time"

    "github.com/eriq-augustine/comic-server/util"
)

const SERIES_IMAGE_BASEDIR = "series";
const COVER_IMAGE_FILENAME = "cover";

type Series struct {
    ID int
    Name string
    Author *string
    Year *int
    URL *string
    Description *string
    CoverImageRelPath *string
    MetadataSource *string
    MetadataSourceID *string
}

func EmptySeries() *Series {
    return &Series{ID: -1};
}

func (this *Series) AssumeCrawl(crawl *MetadataCrawl) error {
    this.MetadataSource = &crawl.MetadataSource;
    this.MetadataSourceID = &crawl.MetadataSourceID;

    if (this.Author == nil) {
        this.Author = crawl.Author;
    }

    if (this.Year == nil) {
        this.Year = crawl.Year;
    }

    if (this.URL == nil) {
        this.URL = crawl.URL;
    }

    if (this.Description == nil) {
        this.Description = crawl.Description;
    }

    if (this.CoverImageRelPath == nil) {
        var filename = COVER_IMAGE_FILENAME + filepath.Ext(*crawl.CoverImageRelPath);
        var coverRelPath = filepath.Join(SERIES_IMAGE_BASEDIR, fmt.Sprintf("%06d", this.ID), filename);

        _, err := util.CopyImage(*crawl.CoverImageRelPath, coverRelPath);
        if (err != nil) {
            return fmt.Errorf("Failed to copy image from '%s' to '%s': %w.", *crawl.CoverImageRelPath, coverRelPath, err);
        }

        this.CoverImageRelPath = &coverRelPath;
    }

    return nil;
}

func (this *Series) String() string {
    text, _ := util.ToJSON(this);
    return text;
}

// Archives are things the relate to physical files/directories on disk.
// The are (or can be) packaged up in a single CBZ file.
type Archive struct {
    ID int
    Path string
    Series *Series
    Volume *string
    Chapter *string
    PageCount *int
}

func EmptyArchive(path string) *Archive {
    return &Archive{
        ID: -1,
        Path: path,
        Series: EmptySeries()}
    ;
}

// Assume all the attributes of other.
func (this *Archive) Assume(other *Archive) {
    this.ID = other.ID;
    this.Path = other.Path;
    this.Series = other.Series;
    this.Volume = other.Volume;
    this.Chapter = other.Chapter;
    this.PageCount = other.PageCount;
}

func (this *Archive) String() string {
    text, _ := util.ToJSON(this);
    return text;
}

type MetadataCrawlRequest struct {
    ID int
    Series *Series
    Query string
    Timestamp time.Time
}

func EmptyCrawlRequest() *MetadataCrawlRequest {
    return &MetadataCrawlRequest{
        ID: -1,
        Series: EmptySeries(),
    };
}

func (this *MetadataCrawlRequest) String() string {
    text, _ := util.ToJSON(this);
    return text;
}

type MetadataCrawl struct {
    ID int
    MetadataSource string
    MetadataSourceID string
    SourceSeries *Series
    Name string
    Author *string
    Year *int
    URL *string
    Description *string
    CoverImageRelPath *string
    Timestamp time.Time
}

func EmptyCrawl(source string, sourceID string) *MetadataCrawl {
    return &MetadataCrawl{
        ID: -1,
        MetadataSource: source,
        MetadataSourceID: sourceID,
        SourceSeries: EmptySeries(),
        Timestamp: time.Now(),
    };
}

func (this *MetadataCrawl) String() string {
    text, _ := util.ToJSON(this);
    return text;
}
