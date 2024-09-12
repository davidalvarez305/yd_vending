const params = new URLSearchParams(window.location.search);
const form = document.getElementById("filters");
const clear = document.getElementById("clearButton");

clear.addEventListener('click', () => window.location.href = "/crm/leads");

document.querySelectorAll('.pagination-link').forEach(item => {
    item.addEventListener('click', e => {
        const currentPage = params.get('page_num') ?? '0';
        const page = parseInt(currentPage);
        const maxPages = parseInt(document.getElementById('maxPages').textContent);
        const chevronValue = item.getAttribute('name');

        if (chevronValue === "left" && page > 1) {
            params.set('page_num', page - 1);
        }

        if (chevronValue === "right" && page < maxPages && maxPages > 1) {
            params.set('page_num', page + 1);
        }

        if (!params.toString()) return;
        updateURL();
    });
});

form.addEventListener("submit", e => {
    e.preventDefault();

    console.log(e.target);
    const formData = new FormData(e.target);
    for (const [key, value] of formData.entries()) {

        // Only append truthy values
        if (!value) continue;

        params.set(key, value);
    }

    // Do not submit if params are empty
    if (!params.toString()) return;

    updateURL();
});

function updateURL() {
    const { origin, pathname } = window.location;
    const url = new URL(origin + pathname);
    url.search = params.toString();

    window.location.replace(url.href);
}