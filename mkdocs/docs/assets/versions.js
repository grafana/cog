window.addEventListener("DOMContentLoaded", function() {
    const topic = document.getElementsByClassName("md-header__topic");
    if (topic.length === 0) {
        return;
    }

    const urlParts = window.location.pathname.split('/').filter(elm => elm !== "");
    const currentVersionFromURL = urlParts.length !== 0 ? urlParts[0] : null;

    const elementFromHTML = (htmlString) =>  {
        const div = document.createElement('div');
        div.innerHTML = htmlString;
        return div.firstChild;
    };

    const urlForVersion = (url, version) => {
        const currentVersionPrefix = '/'+currentVersionFromURL;

        if (currentVersionFromURL === null || !url.startsWith(currentVersionPrefix)) {
            return '/'+version.version;
        }

        return '/'+version.version+url.substring(currentVersionPrefix.length);
    };

    const renderVersionSelector = (versions, currentVersion) => {
        let selector = '<div class="md-version">';

        selector += '<button class="md-version__current">' + currentVersion.title + '</button>';
        selector += '<ul class="md-version__list">';

        versions.forEach((version) => {
            let item = '<li class="md-version__item">';
            item += '<a href="'+ urlForVersion(window.location.pathname, version) +'" class="md-version__link">' + version.title + '</a>';
            item += '</li>';

            selector += item;
        });

        selector += '</ul>';
        selector += '</div>';

        return elementFromHTML(selector);
    };

    fetch('/versions.json').then((response) => {
        return response.json();
    }).then((versions) => {
        const currentVersion = currentVersionFromURL || versions[0].version;
        const currentVersionObj = versions.find((version) => {
            return version.version === currentVersion;
        })
        topic[0].appendChild(renderVersionSelector(versions, currentVersionObj));
    });
});
