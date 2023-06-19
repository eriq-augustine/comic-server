SELECT
    A.id,
    A.series_id,
    A.path,
    A.volume,
    A.chapter,
    A.page_count,
    A.cover_image_relpath
FROM Archives A
WHERE A.series_id = ?
;
