SELECT
    id,
    source,
    source_id,
    origin_series_id,
    name,
    alt_names,
    author,
    year,
    url,
    description,
    cover_image_relpath,
    timestamp
FROM MetadataCrawls
WHERE origin_series_id = ?
;
