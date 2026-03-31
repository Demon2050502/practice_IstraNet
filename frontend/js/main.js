(function () {
    const session = window.appSession;
    if (!session) {
        return;
    }

    const authForm = document.getElementById("auth-form");
    const nameField = document.getElementById("name-field");
    const nameInput = document.getElementById("full-name");
    const roleField = document.getElementById("role-field");
    const roleSelect = document.getElementById("role");
    const emailInput = document.getElementById("email");
    const passwordInput = document.getElementById("password");
    const submitButton = document.getElementById("submit-button");
    const messageBlock = document.getElementById("form-message");
    const accountState = document.getElementById("account-state");
    const accountName = document.getElementById("account-name");
    const accountRole = document.getElementById("account-role");
    const accountHomeLink = document.getElementById("account-home-link");
    const logoutButton = document.getElementById("logout-button");
    const authTabs = document.querySelector(".auth-tabs");
    const modeButtons = document.querySelectorAll("[data-mode-button]");

    if (!authForm || !emailInput || !passwordInput || !submitButton) {
        return;
    }

    let mode = "sign-in";

    function showMessage(type, message) {
        if (!messageBlock) {
            return;
        }

        if (!message) {
            messageBlock.textContent = "";
            messageBlock.className = "alert is-hidden";
            return;
        }

        messageBlock.textContent = message;
        messageBlock.className = `alert alert--${type}`;
    }

    function setLoading(isLoading) {
        submitButton.disabled = isLoading;
        submitButton.textContent = isLoading
            ? "Подождите..."
            : mode === "sign-up"
                ? "Создать аккаунт"
                : "Войти";
    }

    function renderAuthState(user) {
        const hasUser = Boolean(user);

        authForm.classList.toggle("is-hidden", hasUser);
        if (authTabs) {
            authTabs.classList.toggle("is-hidden", hasUser);
        }

        if (!accountState) {
            return;
        }

        accountState.classList.toggle("is-hidden", !hasUser);

        if (!hasUser) {
            return;
        }

        accountName.textContent = user.name || "-";
        accountRole.textContent = session.getRoleLabel(user.role);
        if (accountHomeLink) {
            accountHomeLink.href = session.getHomePathByRole(user.role);
        }
    }

    function switchMode(nextMode) {
        mode = nextMode;

        modeButtons.forEach((button) => {
            button.classList.toggle("auth-tab--active", button.dataset.modeButton === nextMode);
        });

        const isSignUp = nextMode === "sign-up";

        if (nameField) {
            nameField.classList.toggle("is-hidden", !isSignUp);
        }

        if (roleField) {
            roleField.classList.toggle("is-hidden", !isSignUp);
        }

        passwordInput.autocomplete = isSignUp ? "new-password" : "current-password";
        setLoading(false);
        showMessage("", "");
    }

    async function parseResponse(response) {
        try {
            return await response.json();
        } catch (error) {
            return {};
        }
    }

    async function handleSubmit(event) {
        event.preventDefault();
        showMessage("", "");

        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();
        const fullName = nameInput ? nameInput.value.trim() : "";

        if (!email || !password) {
            showMessage("error", "Заполните email и пароль.");
            return;
        }

        if (password.length < 6) {
            showMessage("error", "Пароль должен содержать минимум 6 символов.");
            return;
        }

        if (mode === "sign-up" && !fullName) {
            showMessage("error", "Укажите полное имя.");
            return;
        }

        const endpoint = mode === "sign-up" ? "/auth/sign-up" : "/auth/sign-in";
        const payload = mode === "sign-up"
            ? {
                email,
                password,
                full_name: fullName,
                role: roleSelect ? roleSelect.value.trim() : "user",
            }
            : {
                email,
                password,
            };

        setLoading(true);

        try {
            const response = await fetch(session.getRequestUrl(endpoint), {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                },
                body: JSON.stringify(payload),
            });

            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage("error", data.message || "Не удалось выполнить запрос.");
                return;
            }

            session.saveAuthData(data.token, data.user);
            showMessage("success", mode === "sign-up" ? "Аккаунт создан." : "Авторизация выполнена.");
            renderAuthState(data.user);

            window.setTimeout(() => {
                session.redirectToRoleHome(data.user.role);
            }, 600);
        } catch (error) {
            showMessage("error", "Не удалось связаться с сервером. Проверьте, что приложение запущено локально.");
        } finally {
            setLoading(false);
        }
    }

    modeButtons.forEach((button) => {
        button.addEventListener("click", () => {
            switchMode(button.dataset.modeButton || "sign-in");
        });
    });

    authForm.addEventListener("submit", handleSubmit);

    if (logoutButton) {
        logoutButton.addEventListener("click", () => {
            session.clearStoredAuth();
            renderAuthState(null);
            authForm.reset();
            switchMode("sign-in");
        });
    }

    const currentUser = session.getStoredUser();
    if (currentUser) {
        renderAuthState(currentUser);
    } else {
        renderAuthState(null);
        switchMode("sign-in");
    }
})();
