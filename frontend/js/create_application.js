(function () {
    const session = window.appSession;
    if (!session) {
        return;
    }

    const currentUser = session.requireAuth(["user"]);
    if (!currentUser) {
        return;
    }

    const userName = document.getElementById("workspace-user-name");
    const logoutButton = document.getElementById("logout-button");
    const form = document.getElementById("create-application-form");
    const messageBlock = document.getElementById("form-message");
    const titleInput = document.getElementById("title");
    const descriptionInput = document.getElementById("description");
    const prioritySelect = document.getElementById("priority");
    const phoneInput = document.getElementById("contact-phone");
    const addressInput = document.getElementById("contact-address");
    const submitButton = document.getElementById("submit-button");

    if (userName) {
        userName.textContent = currentUser.name || "-";
    }

    if (logoutButton) {
        logoutButton.addEventListener("click", () => {
            session.clearStoredAuth();
            window.location.replace("/auth");
        });
    }

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
        if (!submitButton) {
            return;
        }

        submitButton.disabled = isLoading;
        submitButton.textContent = isLoading ? "Отправка..." : "Отправить заявку";
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

        const title = titleInput.value.trim();
        const description = descriptionInput.value.trim();
        const priorityCode = prioritySelect.value;
        const contactPhone = phoneInput.value.trim();
        const contactAddress = addressInput.value.trim();

        if (!title || !description) {
            showMessage("error", "Заполните тему и описание заявки.");
            return;
        }

        setLoading(true);

        try {
            const response = await session.authorizedFetch("/applications/create-app", {
                method: "POST",
                body: JSON.stringify({
                    title,
                    description,
                    priority_code: priorityCode,
                    contact_phone: contactPhone || undefined,
                    contact_address: contactAddress || undefined,
                }),
            });

            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage("error", data.message || "Не удалось создать заявку.");
                return;
            }

            showMessage("success", "Заявка успешно создана.");
            form.reset();

            window.setTimeout(() => {
                const detailsPath = data && data.id ? `/account?id=${data.id}` : "/account";
                window.location.href = detailsPath;
            }, 700);
        } catch (error) {
            showMessage("error", "Не удалось отправить заявку. Попробуйте ещё раз.");
        } finally {
            setLoading(false);
        }
    }

    if (form) {
        form.addEventListener("submit", handleSubmit);
    }
})();
