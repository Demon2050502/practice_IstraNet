(function () {
    const session = window.appSession;
    if (!session) {
        return;
    }

    const currentUser = session.requireAuth(["operator"]);
    if (!currentUser) {
        return;
    }

    const userName = document.getElementById("workspace-user-name");
    const logoutButton = document.getElementById("logout-button");
    const filterBar = document.getElementById("filter-bar");
    const searchInput = document.getElementById("search-input");
    const listMessage = document.getElementById("list-message");
    const detailMessage = document.getElementById("detail-message");
    const applicationList = document.getElementById("application-list");
    const placeholder = document.getElementById("application-placeholder");
    const details = document.getElementById("application-details");
    const actionComment = document.getElementById("action-comment");
    const actionsBox = document.getElementById("operator-actions");
    const commentsList = document.getElementById("comments-list");
    const historyList = document.getElementById("history-list");

    const detailTitle = document.getElementById("detail-title");
    const detailDescription = document.getElementById("detail-description");
    const detailStatus = document.getElementById("detail-status");
    const detailPriority = document.getElementById("detail-priority");
    const detailCreatedBy = document.getElementById("detail-created-by");
    const detailAssignedTo = document.getElementById("detail-assigned-to");
    const detailCreatedAt = document.getElementById("detail-created-at");
    const detailUpdatedAt = document.getElementById("detail-updated-at");
    const detailPhone = document.getElementById("detail-phone");
    const detailAddress = document.getElementById("detail-address");

    const statTotal = document.getElementById("stat-total");
    const statNew = document.getElementById("stat-new");
    const statProgress = document.getElementById("stat-progress");
    const statWaiting = document.getElementById("stat-waiting");
    const statResolved = document.getElementById("stat-resolved");
    const statUnassigned = document.getElementById("stat-unassigned");

    const state = {
        applications: [],
        selectedID: 0,
        selectedApplication: null,
        history: [],
        filter: "all",
        search: "",
    };

    if (userName) {
        userName.textContent = currentUser.name || "-";
    }

    logoutButton.addEventListener("click", () => {
        session.clearStoredAuth();
        window.location.replace("/auth");
    });

    function escapeHTML(value) {
        return String(value ?? "")
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll("\"", "&quot;")
            .replaceAll("'", "&#39;");
    }

    function showMessage(node, type, message) {
        if (!message) {
            node.textContent = "";
            node.className = "alert is-hidden";
            return;
        }

        node.textContent = message;
        node.className = `alert alert--${type}`;
    }

    function getStatusBadgeClass(code) {
        switch (code) {
            case "resolved":
            case "closed":
                return "badge badge--success";
            case "waiting":
                return "badge badge--soft";
            case "new":
                return "badge";
            default:
                return "badge badge--muted";
        }
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

    function renderStats() {
        const items = state.applications;
        const countByCode = (code) => items.filter((item) => item.status?.code === code).length;

        statTotal.textContent = String(items.length);
        statNew.textContent = String(countByCode("new"));
        statProgress.textContent = String(countByCode("in_progress"));
        statWaiting.textContent = String(countByCode("waiting"));
        statResolved.textContent = String(items.filter((item) => ["resolved", "closed"].includes(item.status?.code)).length);
        statUnassigned.textContent = String(items.filter((item) => !item.assigned_to).length);
    }

    function getFilteredApplications() {
        return state.applications.filter((item) => {
            if (state.filter === "new" && item.status?.code !== "new") {
                return false;
            }
            if (state.filter === "my" && item.assigned_to?.id !== currentUser.id) {
                return false;
            }
            if (state.filter === "waiting" && item.status?.code !== "waiting") {
                return false;
            }
            if (state.filter === "resolved" && !["resolved", "closed"].includes(item.status?.code)) {
                return false;
            }

            const query = state.search.trim().toLowerCase();
            if (!query) {
                return true;
            }

            const haystack = [item.title, item.created_by?.name, item.description].join(" ").toLowerCase();
            return haystack.includes(query);
        });
    }

    function renderComments(items) {
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

    function renderHistory(items) {
        if (!items.length) {
            historyList.innerHTML = '<div class="empty-state">История действий пока пуста.</div>';
            return;
        }

        historyList.innerHTML = items.map((item) => `
            <article class="timeline-card">
                <div class="timeline-card__head">
                    <h5 class="timeline-card__title">${escapeHTML(item.action || "Действие")}</h5>
                    <span class="badge badge--muted">${escapeHTML(session.formatDateTime(item.created_at))}</span>
                </div>
                <p class="timeline-card__meta">
                    <strong>${escapeHTML(item.actor?.name || "Система")}</strong><br>
                    ${escapeHTML(item.old_value ? `Было: ${item.old_value}` : "")}
                    ${item.old_value && item.new_value ? " | " : ""}
                    ${escapeHTML(item.new_value ? `Стало: ${item.new_value}` : "Действие выполнено")}
                </p>
            </article>
        `).join("");
    }

    function clearDetails() {
        state.selectedApplication = null;
        setSelectedID(0);
        placeholder.classList.remove("is-hidden");
        details.classList.add("is-hidden");
        renderComments([]);
        renderHistory([]);
        actionsBox.innerHTML = "";
    }

    function renderActions(application) {
        const buttons = [];
        const isMine = application.assigned_to?.id === currentUser.id;
        const isFinal = ["resolved", "closed"].includes(application.status?.code);

        if (!application.assigned_to) {
            buttons.push('<button class="button button--primary" type="button" data-action="take">Взять в работу</button>');
        }
        if (isMine && application.status?.code !== "in_progress") {
            buttons.push('<button class="button button--secondary" type="button" data-action="status" data-status="in_progress">В работу</button>');
        }
        if (isMine && application.status?.code !== "waiting" && !isFinal) {
            buttons.push('<button class="button button--secondary" type="button" data-action="status" data-status="waiting">Ожидание</button>');
        }
        if (isMine && application.status?.code !== "resolved" && application.status?.code !== "closed") {
            buttons.push('<button class="button button--secondary" type="button" data-action="status" data-status="resolved">Решена</button>');
        }
        if (isMine && application.status?.code !== "closed") {
            buttons.push('<button class="button button--ghost" type="button" data-action="close">Закрыть</button>');
        }

        if (!buttons.length) {
            actionsBox.innerHTML = '<div class="empty-state">Для текущего статуса быстрые действия недоступны.</div>';
            return;
        }

        actionsBox.innerHTML = buttons.join("");
        actionsBox.querySelectorAll("[data-action]").forEach((button) => {
            button.addEventListener("click", async () => {
                if (button.dataset.action === "take") {
                    await takeApplication();
                    return;
                }
                if (button.dataset.action === "status") {
                    await changeStatus(button.dataset.status);
                    return;
                }
                await closeApplication();
            });
        });
    }

    function renderDetails(application) {
        state.selectedApplication = application;
        setSelectedID(application.id);

        placeholder.classList.add("is-hidden");
        details.classList.remove("is-hidden");

        detailTitle.textContent = application.title || "-";
        detailDescription.textContent = application.description || "-";
        detailStatus.textContent = application.status?.name || "-";
        detailStatus.className = getStatusBadgeClass(application.status?.code);
        detailPriority.textContent = application.priority?.name || "-";
        detailPriority.className = "badge badge--soft";
        detailCreatedBy.textContent = application.created_by?.name || "-";
        detailAssignedTo.textContent = application.assigned_to?.name || "Не назначен";
        detailCreatedAt.textContent = session.formatDateTime(application.created_at);
        detailUpdatedAt.textContent = session.formatDateTime(application.updated_at);
        detailPhone.textContent = application.contact_phone || "Не указан";
        detailAddress.textContent = application.contact_address || "Не указан";
        actionComment.value = "";

        renderComments(Array.isArray(application.comments) ? application.comments : []);
        renderHistory(state.history);
        renderActions(application);
    }

    function renderList() {
        renderStats();

        const items = getFilteredApplications();
        if (!items.length) {
            applicationList.innerHTML = '<div class="empty-state">Подходящих заявок пока нет.</div>';
            if (!state.applications.length) {
                clearDetails();
            }
            return;
        }

        applicationList.innerHTML = items.map((item) => `
            <article class="application-card ${item.id === state.selectedID ? "application-card--active" : ""}" data-app-id="${item.id}">
                <div class="application-card__head">
                    <h3 class="application-card__title">${escapeHTML(item.title || "Без названия")}</h3>
                    <span class="${getStatusBadgeClass(item.status?.code)}">${escapeHTML(item.status?.name || "-")}</span>
                </div>
                <p class="application-card__meta">
                    Автор: ${escapeHTML(item.created_by?.name || "-")}<br>
                    Исполнитель: ${escapeHTML(item.assigned_to?.name || "Не назначен")}<br>
                    Обновлена: ${escapeHTML(session.formatDateTime(item.updated_at))}
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
            const response = await session.authorizedFetch("/api/operator/applications/get-apps");
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
                clearDetails();
                return;
            }

            const urlID = Number(new URL(window.location.href).searchParams.get("id")) || 0;
            const nextID = preferredID || state.selectedID || urlID || state.applications[0].id;
            const existing = state.applications.some((item) => item.id === nextID);
            await loadApplication(existing ? nextID : state.applications[0].id);
        } catch (error) {
            showMessage(listMessage, "error", "Не удалось загрузить список заявок.");
        }
    }

    async function loadApplication(id) {
        if (!id) {
            clearDetails();
            return;
        }

        showMessage(detailMessage, "", "");

        try {
            const [appResponse, historyResponse] = await Promise.all([
                session.authorizedFetch(`/api/operator/applications/get-app?id=${id}`),
                session.authorizedFetch(`/api/operator/applications/get-history?id=${id}`),
            ]);

            const appData = await parseResponse(appResponse);
            const historyData = await parseResponse(historyResponse);

            if (appResponse.status === 404) {
                state.applications = state.applications.filter((item) => item.id !== id);
                renderList();

                if (state.applications.length) {
                    await loadApplication(state.applications[0].id);
                } else {
                    clearDetails();
                }
                return;
            }

            if (!appResponse.ok) {
                showMessage(detailMessage, "error", appData.message || "Не удалось загрузить заявку.");
                return;
            }

            state.history = Array.isArray(historyData.items) ? historyData.items : [];
            renderDetails(appData);
            renderList();
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось загрузить заявку.");
        }
    }

    async function takeApplication() {
        if (!state.selectedApplication) {
            return;
        }

        try {
            const response = await session.authorizedFetch("/api/operator/applications/take-app", {
                method: "PUT",
                body: JSON.stringify({ id: state.selectedApplication.id }),
            });
            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось взять заявку в работу.");
                return;
            }

            showMessage(detailMessage, "success", "Заявка взята в работу.");
            await loadApplications(state.selectedApplication.id);
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось взять заявку в работу.");
        }
    }

    async function changeStatus(statusCode) {
        if (!state.selectedApplication || !statusCode) {
            return;
        }

        try {
            const response = await session.authorizedFetch("/api/operator/applications/change-status", {
                method: "PUT",
                body: JSON.stringify({
                    id: state.selectedApplication.id,
                    status_code: statusCode,
                    comment: actionComment.value.trim() || undefined,
                }),
            });
            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось изменить статус.");
                return;
            }

            showMessage(detailMessage, "success", "Статус обновлён.");
            await loadApplications(state.selectedApplication.id);
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось изменить статус.");
        }
    }

    async function closeApplication() {
        if (!state.selectedApplication) {
            return;
        }

        try {
            const response = await session.authorizedFetch("/api/operator/applications/close-app", {
                method: "PUT",
                body: JSON.stringify({
                    id: state.selectedApplication.id,
                    comment: actionComment.value.trim() || undefined,
                }),
            });
            const data = await parseResponse(response);

            if (!response.ok) {
                showMessage(detailMessage, "error", data.message || "Не удалось закрыть заявку.");
                return;
            }

            showMessage(detailMessage, "success", "Заявка закрыта.");
            await loadApplications(state.selectedApplication.id);
        } catch (error) {
            showMessage(detailMessage, "error", "Не удалось закрыть заявку.");
        }
    }

    filterBar.querySelectorAll("[data-filter]").forEach((button) => {
        button.addEventListener("click", () => {
            state.filter = button.dataset.filter || "all";
            filterBar.querySelectorAll("[data-filter]").forEach((item) => {
                item.classList.toggle("filter-chip--active", item === button);
            });
            renderList();
        });
    });

    searchInput.addEventListener("input", () => {
        state.search = searchInput.value;
        renderList();
    });

    loadApplications();
})();
