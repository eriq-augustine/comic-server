package types

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

type Series struct {
    ID int
    Name string
    Author *string
    Year *int
}

func EmptyArchive() *Archive {
    return &Archive{ID: -1, Series: EmptySeries()};
}

func EmptySeries() *Series {
    return &Series{ID: -1};
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
