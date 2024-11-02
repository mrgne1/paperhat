
    function resetPage() {
        const url = window.location.href;
        const newUrl = url.split("/v1/")[0] + "/v1/";
        window.location.assign(newUrl);
    }

    function copyElementValue(id) {
        const link = document.getElementById(id).value;
        navigator.clipboard.writeText(link);
    }

