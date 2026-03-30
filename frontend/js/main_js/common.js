async function loadLayoutPart(targetSelector, path) {
    const target = document.querySelector(targetSelector);
    if (!target) {
        return;
    }

    try {
        const response = await fetch(path);
        if (!response.ok) {
            return;
        }

        target.innerHTML = await response.text();
    } catch (error) {
        console.error("Не удалось загрузить часть макета", error);
    }
}

async function initLayout() {
    await Promise.all([
        loadLayoutPart("[data-layout-header]", "/components/header.html"),
        loadLayoutPart("[data-layout-footer]", "/components/footer.html"),
    ]);

    if (typeof window.initSiteHeader === "function") {
        window.initSiteHeader();
    }
}

initLayout();
