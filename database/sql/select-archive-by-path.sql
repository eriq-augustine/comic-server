SELECT
    A.id,
    A.volume,
    A.chapter,
    A.page_count,
    A.series_id,
    S.name,
    S.author,
    S.year
FROM
    Archives A
    JOIN Series S ON S.id = A.series_id
WHERE A.path = ?
;
