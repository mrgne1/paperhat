
    function resetPage() {
        const url = window.location.href;
        const newUrl = url.split("/v1/")[0] + "/v1/";
        window.location.assign(newUrl);
    }
