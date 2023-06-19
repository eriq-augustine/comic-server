DROP TABLE IF EXISTS MetadataCrawls;
DROP TABLE IF EXISTS MetadataCrawlRequests;
DROP TABLE IF EXISTS Tags;
DROP TABLE IF EXISTS RelatedSeries;
DROP TABLE IF EXISTS Archives;
DROP TABLE IF EXISTS Series;

CREATE TABLE Series (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    author TEXT,
    year INTEGER,
    url TEXT,
    description TEXT,
    cover_image_relpath TEXT,
    metadata_source TEXT,
    metadata_source_id TEXT
);

CREATE TABLE Archives (
    id INTEGER PRIMARY KEY,
    series_id INTEGER REFERENCES Series(id),
    path TEXT NOT NULL UNIQUE,
    volume TEXT,
    chapter TEXT,
    page_count INTEGER,
    cover_image_relpath TEXT
);

CREATE TABLE RelatedSeries (
    id INTEGER PRIMARY KEY,
    series_a_id INTEGER REFERENCES Series(id),
    series_b_id INTEGER REFERENCES Series(id),
    relation TEXT,
    metadata_source TEXT
);

CREATE TABLE Tags (
    id INTEGER PRIMARY KEY,
    series_id INTEGER REFERENCES Series(id),
    tag TEXT,
    metadata_source TEXT
);

CREATE TABLE MetadataCrawlRequests (
    id INTEGER PRIMARY KEY,
    source_series INTEGER REFERENCES Series(id) NOT NULL,
    query TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE MetadataCrawls (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    source_id TEXT NOT NULL,
    source_series_id INTEGER REFERENCES Series(id) NOT NULL,
    name TEXT NOT NULL,
    author TEXT,
    year INTEGER,
    url TEXT,
    description TEXT,
    cover_image_relpath TEXT,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
