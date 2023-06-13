INSERT INTO MetadataCrawlRequests (
    source_series,
    search_name
) VALUES (
    ?,
    ?
)
ON CONFLICT DO UPDATE SET
    timestamp = CURRENT_TIMESTAMP
;
