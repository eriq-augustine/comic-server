SELECT
    M.id,
    M.query,
    M.timestamp,
    S.id,
    S.name,
    S.author,
    S.year,
    S.metadata_source,
    S.metadata_source_id
FROM
    MetadataCrawlRequests M
    JOIN Series S ON S.id = M.source_series
ORDER BY timestamp
;
