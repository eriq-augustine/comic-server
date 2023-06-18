'use strict';

const DEFAULT_SERIES_COVER_IMAGE = '/static/images/default_series_cover.png';

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

function renderSeries(allSeries) {
    allSeries.sort(function(a, b) {
        return a.Name.localeCompare(b.Name);
    });

    let entries = document.createElement('div');
    entries.className = 'series-list';

    for (const series of allSeries) {
        let image = DEFAULT_SERIES_COVER_IMAGE;
        if (series.CoverImageRelPath) {
            image = '/api/image/blob/' + series.CoverImageRelPath;
        }

        let title = series.Name;
        if (series.Year) {
            title += ` (${series.Year})`;
        }

        entries.insertAdjacentHTML('beforeend', `
            <div class='series-entry' data-id='${series.ID}'>
                <a class='series-thumbnail' href='/static/series/view.html?series=${series.ID}'>
                    <img src='${image}' loading='lazy' alt='${series.Name}' />
                    <div class='series-thumbnail-title'>
                        ${title}
                    </div>
                </a>
            </div>
        `);
    }

    document.querySelector('.page-contents').appendChild(entries);
}

function main() {
    let url = '/api/series/list';

    fetchAPI(url, '.page-contents', 'Fetching series.')
        .then(series => renderSeries(series))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
