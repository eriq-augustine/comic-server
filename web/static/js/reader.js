'use strict';

// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
// {extension: mime, ...}
const SUPPORTED_IMAGE_TYPES = {
    'apng': 'image/apng',
    'avif': 'image/avif',
    'gif': 'image/gif',
    'jpg': 'image/jpeg',
    'jpeg': 'image/jpeg',
    'jfif': 'image/jpeg',
    'pjpeg': 'image/jpeg',
    'pjp': 'image/jpeg',
    'png': 'image/png',
    'svg': 'image/svg+xml',
    'webp': 'image/webp',
};

const NEXT_PAGE = 'NEXT_PAGE';
const PREV_PAGE = 'PREV_PAGE';

// TODO(eriq): Some of these are reading direction dependent.
// Keep a reverse map of shortcuts (since it is shorter), then invert it.
const KEYBOARD_SHORTCUTS_REVERSE = {
    NEXT_PAGE: ['ArrowDown', 'ArrowRight', 'PageDown', 'j', 'J', 'l', 'L'],
    PREV_PAGE: ['ArrowUp', 'ArrowLeft', 'PageUp', 'k', 'K', 'h', 'H'],
};

const KEYBOARD_SHORTCUTS = {}
for (const action of Object.keys(KEYBOARD_SHORTCUTS_REVERSE)) {
    for (const key of KEYBOARD_SHORTCUTS_REVERSE[action]) {
        KEYBOARD_SHORTCUTS[key] = action;
    }
}

function getImageMime(filename) {
    let ext = filename.split('.').pop().toLowerCase();
    if (!SUPPORTED_IMAGE_TYPES.hasOwnProperty(ext)) {
        return undefined;
    }

    return SUPPORTED_IMAGE_TYPES[ext];
}

async function fetchZip(url) {
    console.debug(`Fetchig "${url}".`);

    let files = [];

    const response = await fetch(url);
    if (response.status != 200) {
        throw new Error(`Fetch returned error status: "${response.status}".`);
    }

    const blob = await response.blob();
    const zipReader = new zip.ZipReader(new zip.BlobReader(blob));

    const entries = await zipReader.getEntries();
    for (const entry of entries) {
        let mime = getImageMime(entry.filename);

        const uri = await entry.getData(new zip.Data64URIWriter(mime));

        files.push({
            'filename': entry.filename,
            'data': uri,
        });
    }

    await zipReader.close();

    return files;
}

class ComicReader {
    static nextID = 0;

    id;
    container;

    currentPage;
    pageCount;

    fit;
    seamless;

    constructor(containerQuery) {
        this.id = ComicReader.nextID++;
        this.container = document.querySelector(containerQuery);

        this.currentPage = null;
        this.pageCount = null;

        this.#init();
    }

    load(files) {
        let imagesContainer = this.container.querySelector('.images');

        // TODO(eriq): Sort.

        this.pageCount = 0;
        for (const file of files) {
            let filename = file['filename'];
            if (!getImageMime(filename)) {
                continue;
            }

            this.pageCount++;

            if (this.currentPage === null) {
                this.currentPage = this.pageCount;
            }

            let div = document.createElement('div');
            div.className = 'image';
            div.dataset.filename = filename;
            div.dataset.page = this.pageCount;

            let img = document.createElement('img');
            img.src = file['data'];

            div.appendChild(img);
            imagesContainer.appendChild(div);
        }

        this.container.dataset.pageCount = this.pageCount;
        for (const element of this.container.querySelectorAll('.page-count')) {
            element.innerText = this.pageCount;
        }

        this.navigate(this.currentPage);
    }

    #init() {
        this.#initContainer();

        this.#initKeyboard();

        // Controls.
        this.#initForm('control-fit', 'setFit');
        this.#initForm('control-seamless', 'setSeamless');
        this.#initForm('control-page', 'navigate');
    }

    #initContainer() {
        this.container.innerHTML = TEMPLATE_HTML;
    }

    #initKeyboard() {
        const reader = this;

        document.addEventListener('keydown', (event) => {
            if (!KEYBOARD_SHORTCUTS.hasOwnProperty(event.key)) {
                return;
            }

            // TODO(eriq): Config option for preventing default key actions.
            event.preventDefault();
            reader.navigate(KEYBOARD_SHORTCUTS[event.key]);
        });
    }

    #initForm(formName, callbackMethodName) {
        const form = this.container.querySelector(`form.${formName}-form`);
        const reader = this;

        // These forms are not for submitting.
        form.addEventListener('submit', (event) => { event.preventDefault() }, false);
        form.addEventListener('onsubmit', (event) => {event.preventDefault() }, false);

        reader[callbackMethodName]((new FormData(form)).get(formName));

        form.addEventListener('change', function(event) {
            reader[callbackMethodName]((new FormData(form)).get(formName));
        }, false);
    }

    setSeamless(seamless) {
        seamless = (seamless === 'on');

        if (this.seamless === seamless) {
            return;
        }

        this.seamless = seamless;
        this.container.dataset.seamless = this.seamless;
    }

    setFit(fit) {
        if (this.fit === fit) {
            return;
        }

        this.fit = fit;
        this.container.dataset.fit = this.fit;
    }

    navigate(target) {
        if (target === NEXT_PAGE) {
            target = this.currentPage + 1;
        }

        if (target === PREV_PAGE) {
            target = this.currentPage - 1;
        }

        target = Number.parseInt(target);
        if (!Number.isInteger(target)) {
            return;
        }

        target = Math.max(1, Math.min(target, this.pageCount));

        let image = this.container.querySelector(`.image[data-page='${target}']`);
        if (!image) {
            return;
        }

        this.currentPage = target;
        this.container.querySelector('input#control-page').value = this.currentPage;

        image.scrollIntoView();
    }
}

function main() {
    let reader = new ComicReader('.comic-reader');

    let params = new URLSearchParams(window.location.search);
    let archiveID = params.get('archive');

    let url = `/api/archive/blob/${archiveID}`;

    fetchZip(url)
        .then(files => reader.load(files))
        .catch(error => {
            console.error(error);
        });
}

const TEMPLATE_HTML = `
    <div class='container'>
        <div class='controls'>
            <div class='control control-fit'>
                <form class='control-fit-form'>
                    <legend>Fit:</legend>

                    <label for='control-fit-page'>Page</label>
                    <input type='radio' name='control-fit' id='control-fit-page' value='page' checked />

                    <label for='control-fit-width'>Width</label>
                    <input type='radio' name='control-fit' id='control-fit-width' value='width' />

                    <label for='control-fit-none'>None</label>
                    <input type='radio' name='control-fit' id='control-fit-none' value='none' />
                </form>
            </div>

            <div class='control control-seamless'>
                <form class='control-seamless-form'>
                    <legend>Seamless:</legend>
                    <input type='checkbox' name='control-seamless' id='control-seamless' checked />
                </form>
            </div>

            <div class='control control-page'>
                <form class='control-page-form'>
                    <legend>Page:</legend>
                    <input type='number' name='control-page' id='control-page' min='1' value='1' />
                    <span> / </span><span class='page-count'>?</span>
                </form>
            </div>
        </div>
        <div class='images'>
        </div>
    </div>
`;

document.addEventListener("DOMContentLoaded", main);
