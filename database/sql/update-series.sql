UPDATE SERIES
SET
    name = ?,
    author = ?,
    year = ?,
    url = ?,
    description = ?,
    cover_image_path = ?,
    metadata_source = ?,
    metadata_source_id = ?
WHERE id = ?
;
