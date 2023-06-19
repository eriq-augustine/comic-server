'use strict';

function render(allSeries) {
    allSeries.sort(function(a, b) {
        return a.Name.localeCompare(b.Name);
    });

    let entries = document.createElement('div');
    entries.className = 'series-list';

    for (const series of allSeries) {
        let title = series.Name;
        if (series.Year) {
            title += ` (${series.Year})`;
        }

        entries.insertAdjacentHTML('beforeend', `
            <div class='series-entry preview' data-id='${series.ID}'>
                <a href='/static/series/view/index.html?series=${series.ID}'>
                    <img src='${getSeriesCoverPath(series)}' loading='lazy' alt='${series.Name}' />
                    <div class='title'>
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
        .then(allSeries => render(allSeries))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
