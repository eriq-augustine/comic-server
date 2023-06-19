SELECT
    id,
    name,
    author,
    year,
    url,
    description,
    cover_image_relpath,
    metadata_source,
    metadata_source_id
FROM Series
WHERE metadata_source IS NULL
;
