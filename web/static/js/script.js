'use strict';

const DEFAULT_SERIES_COVER_IMAGE = '/static/images/default_series_cover.png';
const DEFAULT_ARCHIVE_THUMBNAIL_IMAGE = '/static/images/default_series_cover.png';

function getArchiveThumbnailPath(archive) {
    if (archive.CoverImageRelPath) {
        return '/api/image/blob/' + archive.CoverImageRelPath;
    }

    return DEFAULT_ARCHIVE_THUMBNAIL_IMAGE;
}

function getSeriesCoverPath(series) {
    if (series.CoverImageRelPath) {
        return '/api/image/blob/' + series.CoverImageRelPath;
    }

    return DEFAULT_SERIES_COVER_IMAGE;
}

function createArchiveTitle(archive) {
    let titleParts = [];

    if (archive.Volume) {
        titleParts.push(`Volume ${archive.Volume}`);
    }

    if (archive.Chapter) {
        titleParts.push(`Chapter ${archive.Chapter}`);
    }

    if (titleParts.length == 0) {
        return `ID ${archive.ID}`;
    }

    return titleParts.join(' ');
}

async function fetchAPI(url, loadingQuery, loadingMessage) {
    // TODO(eriq): Have a loading element and message.
    if (loadingMessage) {
        console.log(loadingMessage);
    }

    const response = await fetch(url);
    if (response.status != 200) {
        throw new Error(`Fetch returned error status: "${response.status}".`);
    }

    return response.json();
}
