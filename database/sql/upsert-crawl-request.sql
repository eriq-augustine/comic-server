INSERT INTO MetadataCrawlRequests (
    source_series,
    query
) VALUES (
    ?,
    ?
)
ON CONFLICT DO UPDATE SET
    timestamp = CURRENT_TIMESTAMP
;
