SELECT
    A.id,
    A.path,
    A.volume,
    A.chapter,
    A.page_count,
    S.id,
    S.name,
    S.author,
    S.year,
    S.url,
    S.description,
    S.cover_image_relpath,
    S.metadata_source,
    S.metadata_source_id
FROM
    Archives A
    JOIN Series S ON S.id = A.series_id
;
