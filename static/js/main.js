document.body.addEventListener('htmx:beforeSwap', function (evt) {
    if (evt.detail.xhr.status === 400) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = false;
    }
});