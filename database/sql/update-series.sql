UPDATE SERIES
SET
    name = ?,
    alt_names = ?,
    author = ?,
    year = ?,
    url = ?,
    description = ?,
    cover_image_relpath = ?,
    metadata_source = ?,
    metadata_source_id = ?
WHERE id = ?
;
