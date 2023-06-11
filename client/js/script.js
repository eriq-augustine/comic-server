'use strict';

async function fetchAPI(url) {
    const response = await fetch(url);
    if (response.status != 200) {
        throw new Error(`Fetch returned error status: "${response.status}".`);
    }

    return response.json();
}

function loadList(archives) {
    // TEST
    console.log(archives);

    let entries = document.createElement('ul');

    for (const archive of archives) {
        entries.insertAdjacentHTML('beforeend', `
            <p><a href='/client/reader.html?archive=${archive.ID}'>${archive.Filename}</a></p>
        `);
    }

    document.querySelector('.page').appendChild(entries);
}

function main() {
    let url = '/api/list';

    fetchAPI(url)
        .then(archives => loadList(archives))
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", main);
