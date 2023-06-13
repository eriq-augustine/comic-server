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
    year INTEGER
);

CREATE TABLE Archives (
    id INTEGER PRIMARY KEY,
    series_id INTEGER REFERENCES Series(id),
    path TEXT NOT NULL UNIQUE,
    volume TEXT,
    chapter TEXT,
    page_count INTEGER
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
    source_series INTEGER REFERENCES Series(id),
    search_name TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_series, search_name)
);

CREATE TABLE MetadataCrawls (
    id INTEGER PRIMARY KEY,
    source TEXT NOT NULL,
    source_id TEXT NOT NULL,
    series_id INTEGER REFERENCES Series(id),
    name TEXT NOT NULL,
    author TEXT,
    year INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source, source_id)
);
