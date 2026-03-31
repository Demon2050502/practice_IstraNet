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
    const listMessage = document.getElementById("list-message");
    const applicationList = document.getElementById("application-list");
    const placeholder = document.getElementById("application-placeholder");
    const details = document.getElementById("application-details");
    const detailMessage = document.getElementById("detail-message");
    const form = document.getElementById("application-form");
    const saveButton = document.getElementById("save-button");
    const deleteButton = document.getElementById("delete-button");
    const commentsList = document.getElementById("comments-list");

    const detailTitle = document.getElementById("detail-title");
    const detailDescription = document.getElementById("detail-description");
    const detailStatus = document.getElementById("detail-status");
    const detailPriority = document.getElementById("detail-priority");
    const detailCreatedAt = document.getElementById("detail-created-at");
    const detailUpdatedAt = document.getElementById("detail-updated-at");
    const detailAssignedTo = document.getElementById("detail-assigned-to");
    const detailCategory = document.getElementById("detail-category");
    const detailPhone = document.getElementById("detail-phone");
    const detailAddress = document.getElementById("detail-address");

    const titleInput = document.getElementById("application-title");
    const descriptionInput = document.getElementById("application-description");
    const commentInput = document.getElementById("application-comment");

    const statTotal = document.getElementById("stat-total");
    const statActive = document.getElementById("stat-active");
    const statResolved = document.getElementById("stat-resolved");
    const statClosed = document.getElementById("stat-closed");

    const state = {
        applications: [],
        selectedID: 0,
        selectedApplication: null,
    };

    if (userName) {
        userName.textContent = currentUser.name || "-";
    }

    if (logoutButton) {
        logoutButton.addEventListener("click", () => {
            session.clearStoredAuth();
            window.location.replace("/auth");
        });
    }

    function escapeHTML(value) {
        return String(value ?? "")
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll("\"", "&quot;")
            .replaceAll("'", "&#39;");
    }

    function showMessage(node, type, message) {
        if (!node) {
            return;
        }

        if (!message) {
            node.textContent = "";
            node.className = "alert is-hidden";
            return;
        }

        node.textContent = message;
        node.className = `alert alert--${type}`;
    }

    function setSelectedID(id) {
        state.selectedID = id || 0;

        const url = new URL(window.location.href);
        if (state.selectedID > 0) {
            url.searchParams.set("id", String(state.selectedID));
        } else {
            url.searchParams.delete("id");
        }

        window.history.replaceState({}, "", url);
    }

    function getStatusBadgeClass(name) {
        const value = String(name || "").toLowerCase();
        if (value.includes("закры") || value.includes("реш")) {
            return "badge badge--success";
        }
        if (value.includes("нов")) {
            return "badge";
        }
        return "badge badge--soft";
    }

    function renderStats() {
        const items = state.applications;
        const total = items.length;
        const resolved = items.filter((item) => String(item.status || "").toLowerCase().includes("реш")).length;
        const closed = items.filter((item) => String(item.status || "").toLowerCase().includes("закры")).length;
        const active = Math.max(total - resolved - closed, 0);

        if (statTotal) {
            statTotal.textContent = String(total);
        }
        if (statActive) {
            statActive.textContent = String(active);
        }
        if (statResolved) {
            statResolved.textContent = String(resolved);
        }
        if (statClosed) {
            statClosed.textContent = String(closed);
        }
    }

    function renderComments(items) {
        if (!commentsList) {
            return;
        }

        if (!items.length) {
            commentsList.innerHTML = '<div class="empty-state">Комментариев пока нет.</div>';
            return;
        }

        commentsList.innerHTML = items.map((item) => `
            <article class="timeline-card">
                <div class="timeline-card__head">
                    <h5 class="timeline-card__title">${escapeHTML(item.author || "Автор")}</h5>
                    <span class="badge badge--muted">${escapeHTML(session.formatDateTime(item.created_at))}</span>
                </div>
                <p class="timeline-card__meta">${escapeHTML(item.body || "-")}</p>
            </article>
        `).join("");
    }

    function clearDetails() {
        state.selectedApplication = null;
        setSelectedID(0);
        placeholder.classList.remove("is-hidden");
        details.classList.add("is-hidden");
        renderComments([]);
        showMessage(detailMessage, "", "");
    }

    function renderDetails(application) {
        state.selectedApplication = application;
        setSelectedID(application.id);
        placeholder.classList.add("is-hidden");
        details.classList.remove("is-hidden");

        detailTitle.textContent = application.title || "-";
        detailDescription.textContent = application.description || "-";
        detailStatus.textContent = application.status?.name || "-";
        detailStatus.className = getStatusBadgeClass(application.status?.name);
        detailPriority.textContent = application.priority?.name || "-";
        detailPriority.className = "badge badge--soft";
        detailCreatedAt.textContent = session.formatDateTime(application.created_at);
        detailUpdatedAt.textContent = session.formatDateTime(application.updated_at);
        detailAssignedTo.textContent = application.assigned_to?.name || "Не назначен";
        detailCategory.textContent = application.category?.name || "Не указана";
        detailPhone.textContent = application.contact_phone || "Не указан";
        detailAddress.textContent = application.contact_address || "Не указан";

        titleInput.value = application.title || "";
        descriptionInput.value = application.description || "";
        commentInput.value = "";

        renderComments(Array.isArray(application.comments) ? application.comments : []);
    }

    function renderList() {
        renderStats();

        if (!state.applications.length) {
            applicationList.innerHTML = '<div class="empty-state">У вас пока нет обращений.</div>';
            clearDetails();
            return;
        }

        applicationList.innerHTML = state.applications.map((item) => `
            <article class="application-card ${item.id === state.selectedID ? "application-card--active" : ""}" data-app-id="${item.id}">
                <div class="application-card__head">
                    <h3 class="application-card__title">${escapeHTML(item.title || "Без названия")}</h3>
                    <span class="${getStatusBadgeClass(item.status)}">${escapeHTML(item.status || "-")}</span>
                </div>
                <p class="application-card__meta">
                    Приоритет: ${escapeHTML(item.priority || "-")}<br>
                    Создана: ${escapeHTML(session.formatDateTime(item.created_at))}
                </p>
            </article>
        `).join("");

        applicationList.querySelectorAll("[data-app-id]").forEach((node) => {
            node.addEventListener("click", () => {
                loadApplication(Number(node.dataset.appId));
            });
        });
    }

    async function parseResponse(response) {
        try {
            return await response.json();
        } catch (error) {
            return {};
        }
    }

    async function loadApplications(preferredID) {
        showMessage(listMessage, "", "");

        try {
            const response = await session.authorizedFetch("/applications/get-apps");
            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(listMessage, "error", data.message || "Не удалось загрузить список заявок.");
                state.applications = [];
                renderList();
                return;
            }

            state.applications = Array.isArray(data.items) ? data.items : [];
            renderList();

            if (!state.applications.length) {
                return;
            }

            const urlID = Number(new URL(window.location.href).searchParams.get("id")) || 0;
            const nextID = preferredID || state.selectedID || urlID || state.applications[0].id;
            const hasNext = state.applications.some((item) => item.id === nextID);

            await loadApplication(hasNext ? nextID : state.applications[0].id);
        } catch (error) {
            showMessage(listMessage, "error", "Не удалось загрузить список заявок.");
            state.applications = [];
            renderList();
        }
    }

    async function loadApplication(id) {
        if (!id) {
            clearDetails();
            return;
        }

        showMessage(detailMessage, "", "");

        try {
            const response = await session.authorizedFetch(`/applications/get-app?id=${id}`);
            const data = await parseResponse(response);

            if (response.status === 404) {
                state.applications = state.applications.filter((item) => item.id !== id);
                renderList();

                if (state.applications.length) {
                    await loadApplication(state.applications[0].id);
                } else {
                    clearDetails();
                }
                return;
            }

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось загрузить карточку заявки.");
                return;
            }

            renderDetails(data);
            renderList();
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось загрузить карточку заявки.");
        }
    }

    async function handleSave(event) {
        event.preventDefault();

        if (!state.selectedApplication) {
            return;
        }

        const title = titleInput.value.trim();
        const description = descriptionInput.value.trim();
        const comment = commentInput.value.trim();

        if (!title || !description) {
            showMessage(detailMessage, "error", "Заполните тему и описание заявки.");
            return;
        }

        saveButton.disabled = true;

        try {
            const response = await session.authorizedFetch("/applications/change-app", {
                method: "PUT",
                body: JSON.stringify({
                    id: state.selectedApplication.id,
                    title,
                    description,
                    comment: comment || undefined,
                }),
            });

            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось обновить заявку.");
                return;
            }

            showMessage(detailMessage, "success", "Изменения сохранены.");
            await loadApplications(state.selectedApplication.id);
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось обновить заявку.");
        } finally {
            saveButton.disabled = false;
        }
    }

    async function handleDelete() {
        if (!state.selectedApplication) {
            return;
        }

        if (!window.confirm("Удалить выбранную заявку?")) {
            return;
        }

        deleteButton.disabled = true;

        try {
            const deletingID = state.selectedApplication.id;
            const response = await session.authorizedFetch("/applications/delete-app", {
                method: "DELETE",
                body: JSON.stringify({ id: deletingID }),
            });

            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось удалить заявку.");
                return;
            }

            state.applications = state.applications.filter((item) => item.id !== deletingID);
            renderList();

            if (state.applications.length) {
                await loadApplication(state.applications[0].id);
            } else {
                clearDetails();
            }
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось удалить заявку.");
        } finally {
            deleteButton.disabled = false;
        }
    }

    form.addEventListener("submit", handleSave);
    deleteButton.addEventListener("click", handleDelete);

    loadApplications();
})();
