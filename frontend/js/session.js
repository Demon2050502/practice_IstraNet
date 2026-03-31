const TOKEN_KEY = "istranet_token";
const USER_KEY = "istranet_user";

function getStoredToken() {
    return sessionStorage.getItem(TOKEN_KEY) || "";
}

function getStoredUser() {
    try {
        const raw = sessionStorage.getItem(USER_KEY);
        return raw ? JSON.parse(raw) : null;
    } catch (error) {
        return null;
    }
}

function decodeBase64Url(value) {
    const normalized = value.replace(/-/g, "+").replace(/_/g, "/");
    const padded = normalized.padEnd(normalized.length + ((4 - normalized.length % 4) % 4), "=");
    const binary = atob(padded);
    const bytes = Uint8Array.from(binary, (symbol) => symbol.charCodeAt(0));

    return new TextDecoder().decode(bytes);
}

function getTokenPayload() {
    const token = getStoredToken();
    if (!token) {
        return null;
    }

    const parts = token.split(".");
    if (parts.length < 2) {
        return null;
    }

    try {
        return JSON.parse(decodeBase64Url(parts[1]));
    } catch (error) {
        return null;
    }
}

function getCurrentUserID() {
    const payload = getTokenPayload();
    const userID = Number(payload ? payload.sub : 0);

    return Number.isFinite(userID) && userID > 0 ? userID : 0;
}

function saveAuthData(token, user) {
    sessionStorage.setItem(TOKEN_KEY, token);
    sessionStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearStoredAuth() {
    sessionStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(USER_KEY);
}

function getRoleLabel(role) {
    switch (role) {
        case "user":
            return "Клиент";
        case "operator":
            return "Специалист";
        case "admin":
            return "Администратор";
        default:
            return role || "-";
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
            return "/";
    }
}

function redirectToRoleHome(role) {
    window.location.href = getHomePathByRole(role);
}

function requireAuth(allowedRoles) {
    const token = getStoredToken();
    const user = getStoredUser();

    if (!token || !user) {
        window.location.replace("/auth");
        return null;
    }

    if (Array.isArray(allowedRoles) && allowedRoles.length > 0 && !allowedRoles.includes(user.role)) {
        redirectToRoleHome(user.role);
        return null;
    }

    return {
        ...user,
        id: getCurrentUserID(),
    };
}

function getRequestUrl(path) {
    const apiBaseUrl = window.__DESKNET_API_BASE_URL || "";
    return `${apiBaseUrl}${path}`;
}

async function authorizedFetch(path, options = {}) {
    const token = getStoredToken();
    const headers = new Headers(options.headers || {});

    if (!headers.has("Accept")) {
        headers.set("Accept", "application/json");
    }

    if (options.body && !headers.has("Content-Type")) {
        headers.set("Content-Type", "application/json");
    }

    if (token) {
        headers.set("Authorization", `Bearer ${token}`);
    }

    const response = await fetch(getRequestUrl(path), {
        ...options,
        headers,
    });

    if (response.status === 401) {
        clearStoredAuth();
        window.location.replace("/auth");
        throw new Error("Требуется повторная авторизация");
    }

    return response;
}

function formatDate(value) {
    if (!value) {
        return "-";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return "-";
    }

    return new Intl.DateTimeFormat("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
    }).format(date);
}

function formatDateTime(value) {
    if (!value) {
        return "-";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
        return "-";
    }

    return new Intl.DateTimeFormat("ru-RU", {
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
    }).format(date);
}

window.appSession = {
    TOKEN_KEY,
    USER_KEY,
    getStoredToken,
    getStoredUser,
    getTokenPayload,
    getCurrentUserID,
    saveAuthData,
    clearStoredAuth,
    getRoleLabel,
    getHomePathByRole,
    redirectToRoleHome,
    requireAuth,
    getRequestUrl,
    authorizedFetch,
    formatDate,
    formatDateTime,
};
