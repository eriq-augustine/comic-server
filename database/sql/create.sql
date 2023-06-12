DROP TABLE IF EXISTS MetadataSources;
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
    path TEXT NOT NULL,
    volume TEXT,
    chater TEXT,
    num_pages INTEGER
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

CREATE TABLE MetadataSources (
    id INTEGER PRIMARY KEY,
    series_id INTEGER REFERENCES Series(id),
    name TEXT,
    author TEXT,
    year INTEGER,
    source TEXT,
    timestamp INTEGER
);
