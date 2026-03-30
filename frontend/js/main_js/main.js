const API_BASE_URL = "";
const TOKEN_KEY = "istranet_token";
const USER_KEY = "istranet_user";

const authForm = document.getElementById("auth-form");
const nameField = document.getElementById("name-field");
const nameInput = document.getElementById("full-name");
const emailInput = document.getElementById("email");
const passwordInput = document.getElementById("password");
const submitButton = document.getElementById("submit-button");
const messageBlock = document.getElementById("form-message");
const accountState = document.getElementById("account-state");
const accountName = document.getElementById("account-name");
const accountRole = document.getElementById("account-role");
const logoutButton = document.getElementById("logout-button");
const authTabs = document.querySelector(".auth-tabs");
const modeButtons = document.querySelectorAll("[data-mode-button]");

let currentMode = "sign-in";
let isLoading = false;

function init() {
    bindEvents();
    restoreAuthState();
    setMode(currentMode);
}

function bindEvents() {
    authForm.addEventListener("submit", handleSubmit);
    logoutButton.addEventListener("click", handleLogout);

    modeButtons.forEach((button) => {
        button.addEventListener("click", () => {
            setMode(button.dataset.modeButton);
        });
    });
}

function getRoleLabel(role) {
    switch (role) {
        case "user":
            return "Клиент";
        case "operator":
            return "Специалист";
        case "admin":
            return "Руководитель";
        default:
            return role || "-";
    }
}

function setMode(mode) {
    currentMode = mode;

    modeButtons.forEach((button) => {
        const isActive = button.dataset.modeButton === mode;
        button.classList.toggle("auth-tab--active", isActive);
    });

    const isSignUp = mode === "sign-up";
    nameField.classList.toggle("is-hidden", !isSignUp);
    nameInput.required = isSignUp;
    nameInput.autocomplete = isSignUp ? "name" : "off";
    passwordInput.autocomplete = isSignUp ? "new-password" : "current-password";
    submitButton.textContent = isSignUp ? "Зарегистрироваться" : "Войти";

    clearMessage();
    authForm.reset();
}

function restoreAuthState() {
    const token = sessionStorage.getItem(TOKEN_KEY);
    const userRaw = sessionStorage.getItem(USER_KEY);

    if (!token || !userRaw) {
        renderGuestState();
        return;
    }

    try {
        const user = JSON.parse(userRaw);
        renderAuthState(user);
    } catch (error) {
        clearStoredAuth();
        renderGuestState();
    }
}

function renderGuestState() {
    authForm.classList.remove("is-hidden");
    authTabs.classList.remove("is-hidden");
    accountState.classList.add("is-hidden");
}

function renderAuthState(user) {
    authForm.classList.add("is-hidden");
    authTabs.classList.add("is-hidden");
    accountState.classList.remove("is-hidden");
    accountName.textContent = user.name || "-";
    accountRole.textContent = getRoleLabel(user.role);
}

function showMessage(type, text) {
    messageBlock.textContent = text;
    messageBlock.className = "alert";

    if (type === "success") {
        messageBlock.classList.add("alert--success");
    } else if (type === "error") {
        messageBlock.classList.add("alert--error");
    } else {
        messageBlock.classList.add("alert--info");
    }
}

function clearMessage() {
    messageBlock.textContent = "";
    messageBlock.className = "alert is-hidden";
}

function setLoadingState(state) {
    isLoading = state;
    submitButton.disabled = state;

    if (state) {
        submitButton.textContent = currentMode === "sign-up" ? "Регистрация..." : "Вход...";
        return;
    }

    submitButton.textContent = currentMode === "sign-up" ? "Зарегистрироваться" : "Войти";
}

function validateForm() {
    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    const fullName = nameInput.value.trim();

    if (currentMode === "sign-up" && fullName === "") {
        return "Укажите полное имя";
    }

    if (email === "") {
        return "Укажите email";
    }

    if (!email.includes("@")) {
        return "Укажите корректный email";
    }

    if (password === "") {
        return "Укажите пароль";
    }

    if (currentMode === "sign-up" && password.length < 6) {
        return "Пароль должен быть не короче 6 символов";
    }

    return "";
}

function buildPayload() {
    const payload = {
        email: emailInput.value.trim(),
        password: passwordInput.value.trim(),
    };

    if (currentMode === "sign-up") {
        payload.full_name = nameInput.value.trim();
        payload.role = "user";
    }

    return payload;
}

function getRequestPath() {
    if (currentMode === "sign-up") {
        return "/auth/sign-up";
    }

    return "/auth/sign-in";
}

function getRequestUrl(path) {
    return API_BASE_URL ? `${API_BASE_URL}${path}` : path;
}

async function sendAuthRequest() {
    const response = await fetch(getRequestUrl(getRequestPath()), {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(buildPayload()),
    });

    const data = await response.json().catch(() => null);

    if (!response.ok) {
        const message = data && data.message ? data.message : "Не удалось выполнить запрос";
        throw new Error(message);
    }

    if (!data || !data.token || !data.user) {
        throw new Error("Сервис временно недоступен");
    }

    return data;
}

function saveAuthData(token, user) {
    sessionStorage.setItem(TOKEN_KEY, token);
    sessionStorage.setItem(USER_KEY, JSON.stringify(user));
}

function clearStoredAuth() {
    sessionStorage.removeItem(TOKEN_KEY);
    sessionStorage.removeItem(USER_KEY);
}

async function handleSubmit(event) {
    event.preventDefault();

    if (isLoading) {
        return;
    }

    clearMessage();

    const validationMessage = validateForm();
    if (validationMessage !== "") {
        showMessage("error", validationMessage);
        return;
    }

    setLoadingState(true);

    try {
        const data = await sendAuthRequest();
        saveAuthData(data.token, data.user);
        renderAuthState(data.user);
        showMessage(
            "success",
            currentMode === "sign-up"
                ? "Регистрация прошла успешно, аккаунт уже активен"
                : "Вход выполнен успешно"
        );
    } catch (error) {
        let message = "Сервис временно недоступен. Попробуйте ещё раз немного позже";

        if (error instanceof Error && error.message !== "Failed to fetch") {
            message = error.message;
        }

        showMessage("error", message);
    } finally {
        setLoadingState(false);
    }
}

function handleLogout() {
    clearStoredAuth();
    authForm.reset();
    renderGuestState();
    setMode("sign-in");
    showMessage("success", "Вы вышли из аккаунта");
}

init();
