SELECT
    id,
    author,
    year,
    metadata_source,
    metadata_source_id
FROM Series
WHERE id = ?
;
