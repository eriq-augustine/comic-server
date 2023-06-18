'use strict';

function render(series) {
    // TEST
    console.log(series);
}

function main() {
    let params = new URLSearchParams(window.location.search);
    let seriesID = params.get('series');

    let url = `/api/series/${seriesID}`;

    fetchAPI(url, '.page-contents', 'Fetching series data.')
        .then(series => render(series))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
