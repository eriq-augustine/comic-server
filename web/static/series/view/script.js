'use strict';

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

function renderArchives(archives) {
    archives.sort(function(a, b) {
        let result = a.Volume.localeCompare(b.Volume);
        if (result != 0) {
            return result;
        }

        result = a.Chapter.localeCompare(b.Chapter);
        if (result != 0) {
            return result;
        }

        return a.Path.localeCompare(b.Path);
    });

    let entries = document.createElement('div');
    entries.className = 'archive-list';

    for (const archive of archives) {
        let title = createArchiveTitle(archive);

        entries.insertAdjacentHTML('beforeend', `
            <div class='archive-entry preview' data-id='${archive.ID}'>
                <a href='/static/reader/index.html?archive=${archive.ID}'>
                    <img src='${getArchiveThumbnailPath(archive)}' alt='${title}' />
                    <div class='title'>
                        ${title}
                    </div>
                </a>
            </div>
        `);
    }

    document.querySelector('.series-info .archives').appendChild(entries);
}

function render(series) {
    fetchAPI(`/api/archive/series/${series.ID}`, '.series-info .archives', 'Fetching archives.')
        .then(archives => renderArchives(archives))
        .catch(error => {
            console.error(error);
        });

    let title = series.Name;
    if (series.URL) {
        title = `<a href='${series.URL}'>${series.Name}</a>`;
    }

    let html = `
        <div class='series-info'>
            <h2 class='title'>
                ${title}
            </h2>
            <div class='cover'>
                <img src='${getSeriesCoverPath(series)}' alt='${series.Name}' />
            </div>
            <div class='author'>
                by ${series.Author}
            </div>
            <div class='year'>
                ${series.Year}
            </div>
            <div class='description'>
                ${series.Description}
            </div>
            <div class='archives'>
            </div>
        </div>
    `;

    document.querySelector('.page-contents').innerHTML = html;
}

function main() {
    let params = new URLSearchParams(window.location.search);
    let seriesID = params.get('series');

    fetchAPI(`/api/series/${seriesID}`, '.page-contents', 'Fetching series data.')
        .then(series => render(series))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
