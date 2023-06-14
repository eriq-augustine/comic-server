package model

import (
    "time"
)

type Series struct {
    ID int
    Name string
    Author *string
    Year *int
    MetadataSource *string
    MetadataSourceID *string
}

func EmptySeries() *Series {
    return &Series{ID: -1};
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

func EmptyArchive() *Archive {
    return &Archive{ID: -1, Series: EmptySeries()};
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

type MetadataCrawlRequest struct {
    ID int
    Series *Series
    Query string
    Timestamp time.Time
}

func EmptyCrawlRequest() *MetadataCrawlRequest {
    return &MetadataCrawlRequest{ID: -1, Series: EmptySeries()};
}

type MetadataCrawl struct {
    ID int
    MetadataSource *string
    MetadataSourceID *string
    SourceSeries *Series
    Name string
    Author *string
    Year *int
    Query string
    Timestamp time.Time
}

func EmptyCrawl() *MetadataCrawl {
    return &MetadataCrawl{ID: -1, SourceSeries: EmptySeries()};
}
