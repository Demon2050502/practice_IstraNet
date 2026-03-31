function initSiteHeader() {
    const header = document.querySelector(".site-header");
    const toggle = document.querySelector("[data-header-toggle]");
    const menu = document.querySelector("[data-header-menu]");
    const overlay = document.querySelector("[data-header-overlay]");
    const closeButton = document.querySelector("[data-header-close]");

    if (!header || !toggle || !menu || !overlay || !closeButton) {
        return;
    }

    const closeMenu = () => {
        header.classList.remove("site-header--open");
        document.body.classList.remove("site-menu-open");
        toggle.setAttribute("aria-expanded", "false");
    };

    const openMenu = () => {
        header.classList.add("site-header--open");
        document.body.classList.add("site-menu-open");
        toggle.setAttribute("aria-expanded", "true");
    };

    toggle.addEventListener("click", () => {
        if (header.classList.contains("site-header--open")) {
            closeMenu();
            return;
        }

        openMenu();
    });

    overlay.addEventListener("click", closeMenu);
    closeButton.addEventListener("click", closeMenu);

    menu.querySelectorAll("a").forEach((link) => {
        link.addEventListener("click", () => {
            closeMenu();
        });
    });

    document.addEventListener("keydown", (event) => {
        if (event.key === "Escape") {
            closeMenu();
        }
    });

    window.addEventListener("resize", () => {
        if (window.innerWidth > 860) {
            closeMenu();
        }
    });
}

window.initSiteHeader = initSiteHeader;
