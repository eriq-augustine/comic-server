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

    document.querySelector('.archives').appendChild(entries);
}

function render(series) {
    let title = series.Name;
    if (series.URL) {
        title = `<a href='${series.URL}'>${series.Name}</a>`;
    }
    series['Title'] = title;

    if (series.Description) {
        series.Description = series.Description.trim();
    }

    let infoPanels = [
        {'key': 'Title', 'classname': 'title', 'label': 'Title'},
        {'key': 'Year', 'classname': 'year', 'label': 'Year'},
        {'key': 'Author', 'classname': 'author', 'label': 'Author'},
        {'key': 'Description', 'classname': 'description', 'label': 'Description'},
    ];
    let infoPanelParts = [];

    for (const infoPanel of infoPanels) {
        infoPanelParts.push(`
            <div class='series-info-entry'>
                <label>${infoPanel['label']}</label>
                <div class='${infoPanel['classname']}'>${series[infoPanel['key']]}</div>
            </div>
        `);
    }

    let html = `
        <div class='left-panel'>
            <div class='cover'>
                <img src='${getSeriesCoverPath(series)}' alt='${series.Name}' />
            </div>
        </div>
        <div class='right-panel'>
            ${infoPanelParts.join("\n")}
        </div>
    `;

    let seriesInfo = document.createElement('div');
    seriesInfo.className = 'series-info';
    seriesInfo.innerHTML = html;
    document.querySelector('.page-contents').appendChild(seriesInfo);

    // Create a stub for the archives.
    let archives = document.createElement('div');
    archives.className = 'archives';
    document.querySelector('.page-contents').appendChild(archives);

    // Get the archives.
    fetchAPI(`/api/archive/series/${series.ID}`, '.archives', 'Fetching archives.')
        .then(archives => renderArchives(archives))
        .catch(error => {
            console.error(error);
        });
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
