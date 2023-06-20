'use strict';

function render(archive) {
    let html = `
        <div class='info'>
            <div class='series'>
                <a href="/static/series/view/index.html?series=${archive.Series.ID}">
                    ${archive.Series.Name}
                </a>
            </div>
            <div class='title'>
                ${createArchiveTitle(archive)}
            </div>
        </div>
        <div class='comic-reader'>
        </div>
    `;

    let readerArea = document.createElement('div');
    readerArea.classname = 'reader-area';
    readerArea.innerHTML = html;

    document.querySelector('.page-contents').appendChild(readerArea);

    let reader = new ComicReader('.comic-reader');
    fetchZip(`/api/archive/blob/${archive.ID}`)
        .then(files => reader.load(files))
        .catch(error => {
            console.error(error);
        });
}

function main() {
    let params = new URLSearchParams(window.location.search);
    let archiveID = params.get('archive');

    fetchAPI(`/api/archive/${archiveID}`, '.page-contents', 'Fetching archive information.')
        .then(archive => render(archive))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
