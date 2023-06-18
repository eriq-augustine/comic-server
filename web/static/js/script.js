'use strict';

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
