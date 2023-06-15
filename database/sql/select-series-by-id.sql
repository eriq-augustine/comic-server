SELECT
    id,
    name,
    author,
    year,
    url,
    description,
    cover_image_path,
    metadata_source,
    metadata_source_id
FROM Series
WHERE id = ?
;
