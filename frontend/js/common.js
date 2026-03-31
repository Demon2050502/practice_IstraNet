async function loadLayoutPart(targetSelector, path) {
    const target = document.querySelector(targetSelector);
    if (!target) {
        return;
    }

    try {
        const response = await fetch(path, {
            cache: "no-store",
        });

        if (!response.ok) {
            return;
        }

        target.innerHTML = await response.text();
    } catch (error) {
        console.error("Не удалось загрузить часть макета", error);
    }
}

function getStoredUser() {
    try {
        const raw = sessionStorage.getItem("istranet_user");
        return raw ? JSON.parse(raw) : null;
    } catch (error) {
        return null;
    }
}

function getHomePathByRole(role) {
    switch (role) {
        case "user":
            return "/account";
        case "operator":
            return "/operator";
        case "admin":
            return "/admin";
        default:
            return "/auth";
    }
}

function syncPublicActions() {
    const user = getStoredUser();
    const authLink = document.querySelector("[data-public-auth-link]");
    const cabinetLink = document.querySelector("[data-public-cabinet-link]");

    if (!authLink && !cabinetLink) {
        return;
    }

    if (!user || !user.role) {
        if (authLink) {
            authLink.textContent = "Войти";
            authLink.setAttribute("href", "/auth");
        }

        if (cabinetLink) {
            cabinetLink.textContent = "Открыть кабинет";
            cabinetLink.setAttribute("href", "/auth");
        }
        return;
    }

    const homePath = getHomePathByRole(user.role);

    if (authLink) {
        authLink.textContent = "Личный кабинет";
        authLink.setAttribute("href", homePath);
    }

    if (cabinetLink) {
        cabinetLink.textContent = "Открыть кабинет";
        cabinetLink.setAttribute("href", homePath);
    }
}

async function initLayout() {
    await Promise.all([
        loadLayoutPart("[data-layout-header]", "/components/header.html"),
        loadLayoutPart("[data-layout-footer]", "/components/footer.html"),
    ]);

    syncPublicActions();

    if (typeof window.initSiteHeader === "function") {
        window.initSiteHeader();
    }
}

initLayout();
