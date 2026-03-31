(function () {
    const s = window.appSession;
    if (!s) return;
    const me = s.requireAuth(["admin"]);
    if (!me) return;

    const q = (id) => document.getElementById(id);
    const esc = (v) => String(v ?? "").replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll("\"", "&quot;").replaceAll("'", "&#39;");
    const msg = (id, type, text) => {
        const el = q(id);
        if (!text) {
            el.textContent = "";
            el.className = "alert is-hidden";
            return;
        }
        el.textContent = text;
        el.className = `alert alert--${type}`;
    };
    const json = async (r) => { try { return await r.json(); } catch { return {}; } };
    const badge = (code, fin) => fin || code === "resolved" || code === "closed" ? "badge badge--success" : code === "new" ? "badge" : "badge badge--soft";
    const PAGE = 6;
    const baseStatuses = [
        { id: 1, code: "new", name: "Новая", is_final: false },
        { id: 2, code: "in_progress", name: "В работе", is_final: false },
        { id: 3, code: "waiting", name: "Ожидание", is_final: false },
        { id: 4, code: "resolved", name: "Решена", is_final: true },
        { id: 5, code: "closed", name: "Закрыта", is_final: true },
    ];
    const st = {
        apps: [],
        users: [],
        statuses: [...baseStatuses],
        appPage: 1,
        userPage: 1,
        appFilter: "all",
        appSearch: "",
        userSearch: "",
        historySearch: "",
        appID: 0,
        userID: 0,
        app: null,
        user: null,
        history: [],
        statusID: 0,
    };

    q("workspace-user-name").textContent = me.name || "-";
    q("logout-button").addEventListener("click", () => {
        s.clearStoredAuth();
        window.location.replace("/auth");
    });

    const findStatus = (code) => st.statuses.find((x) => x.code === code) || null;
    const paginate = (items, page) => {
        const pages = Math.max(Math.ceil(items.length / PAGE), 1);
        const p = Math.min(Math.max(page, 1), pages);
        return { page: p, pages, items: items.slice((p - 1) * PAGE, p * PAGE) };
    };

    function syncStatuses() {
        const map = new Map(baseStatuses.map((x) => [x.code, { ...x }]));
        st.statuses.forEach((x) => x?.code && map.set(x.code, { ...x }));
        st.apps.forEach((x) => {
            if (!x.status?.code) return;
            const cur = map.get(x.status.code) || {};
            map.set(x.status.code, {
                id: cur.id || 0,
                code: x.status.code,
                name: x.status.name || cur.name || x.status.code,
                is_final: Boolean(cur.is_final || ["resolved", "closed"].includes(x.status.code)),
            });
        });
        st.statuses = Array.from(map.values()).sort((a, b) => (Number(a.id) || 999) - (Number(b.id) || 999));
    }

    function renderStats() {
        q("stat-users").textContent = String(st.users.filter((x) => x.role_code === "user").length);
        q("stat-operators").textContent = String(st.users.filter((x) => x.role_code === "operator").length);
        q("stat-admins").textContent = String(st.users.filter((x) => x.role_code === "admin").length);
        q("stat-applications").textContent = String(st.apps.length);
        q("stat-new").textContent = String(st.apps.filter((x) => x.status?.code === "new").length);
        q("stat-final").textContent = String(st.apps.filter((x) => findStatus(x.status?.code)?.is_final || ["resolved", "closed"].includes(x.status?.code)).length);
    }

    function fillOperators() {
        q("assign-operator").innerHTML = st.users
            .filter((x) => x.role_code === "operator")
            .map((x) => `<option value="${x.id}">${esc(x.full_name)} (${esc(x.email)})</option>`)
            .join("");
    }

    function fillStatuses() {
        q("application-status-select").innerHTML = st.statuses.map((x) => `<option value="${esc(x.code)}">${esc(x.name)}</option>`).join("");
    }

    function appRows() {
        return st.apps.filter((x) => {
            if (st.appFilter === "new" && x.status?.code !== "new") return false;
            if (st.appFilter === "unassigned" && x.assigned_to) return false;
            if (st.appFilter === "final" && !(findStatus(x.status?.code)?.is_final || ["resolved", "closed"].includes(x.status?.code))) return false;
            return [x.title, x.created_by?.name, x.assigned_to?.name, x.description].join(" ").toLowerCase().includes(st.appSearch.toLowerCase());
        });
    }

    function userRows() {
        return st.users.filter((x) => [x.full_name, x.email, x.role_name, x.role_code].join(" ").toLowerCase().includes(st.userSearch.toLowerCase()));
    }

    function renderAppsTable() {
        const p = paginate(appRows(), st.appPage);
        st.appPage = p.page;
        q("application-page-info").textContent = `Страница ${p.page} из ${p.pages}`;
        q("application-page-prev").disabled = p.page <= 1;
        q("application-page-next").disabled = p.page >= p.pages;
        q("application-list").innerHTML = p.items.length ? p.items.map((x) => `
            <tr class="${x.id === st.appID ? "is-active" : ""}" data-app-id="${x.id}">
                <td><button class="admin-table__button" type="button"><span class="admin-table__title">${esc(x.title || "Без названия")}</span><span class="admin-table__meta">ID: ${esc(x.id)}</span></button></td>
                <td>${esc(x.created_by?.name || "-")}</td>
                <td>${esc(x.assigned_to?.name || "Не назначен")}</td>
                <td><span class="${badge(x.status?.code, findStatus(x.status?.code)?.is_final)}">${esc(x.status?.name || "-")}</span></td>
                <td>${esc(s.formatDateTime(x.updated_at))}</td>
            </tr>`).join("") : '<tr><td class="admin-table__empty" colspan="5">Заявок пока нет.</td></tr>';
        q("application-list").querySelectorAll("[data-app-id]").forEach((row) => row.addEventListener("click", () => loadApp(Number(row.dataset.appId))));
    }

    function renderUsersTable() {
        const p = paginate(userRows(), st.userPage);
        st.userPage = p.page;
        q("user-page-info").textContent = `Страница ${p.page} из ${p.pages}`;
        q("user-page-prev").disabled = p.page <= 1;
        q("user-page-next").disabled = p.page >= p.pages;
        q("user-list").innerHTML = p.items.length ? p.items.map((x) => `
            <tr class="${x.id === st.userID ? "is-active" : ""}" data-user-id="${x.id}">
                <td>${esc(x.full_name || "-")}</td>
                <td>${esc(x.email || "-")}</td>
                <td><span class="${x.role_code === "admin" ? "badge badge--soft" : "badge"}">${esc(x.role_name || x.role_code || "-")}</span></td>
                <td><span class="${x.is_active ? "badge badge--success" : "badge badge--danger"}">${x.is_active ? "Активен" : "Отключён"}</span></td>
                <td>${esc(s.formatDate(x.created_at))}</td>
            </tr>`).join("") : '<tr><td class="admin-table__empty" colspan="5">Пользователей пока нет.</td></tr>';
        q("user-list").querySelectorAll("[data-user-id]").forEach((row) => row.addEventListener("click", () => loadUser(Number(row.dataset.userId))));
    }

    function renderHistory() {
        const items = st.history.filter((x) => [x.action, x.actor?.name, x.field, x.old_value, x.new_value, s.formatDateTime(x.created_at)].join(" ").toLowerCase().includes(st.historySearch.toLowerCase()));
        q("application-history").innerHTML = items.length ? items.map((x) => `
            <article class="timeline-card">
                <div class="timeline-card__head"><h5 class="timeline-card__title">${esc(x.action || "Действие")}</h5><span class="badge badge--muted">${esc(s.formatDateTime(x.created_at))}</span></div>
                <p class="timeline-card__meta"><strong>${esc(x.actor?.name || "Система")}</strong><br>${esc(x.old_value ? `Было: ${x.old_value}` : "")}${x.old_value && x.new_value ? " | " : ""}${esc(x.new_value ? `Стало: ${x.new_value}` : "Действие выполнено")}</p>
            </article>`).join("") : '<div class="empty-state">История пока пуста.</div>';
    }

    function renderComments(items) {
        q("application-comments").innerHTML = items.length ? items.map((x) => `
            <article class="timeline-card">
                <div class="timeline-card__head"><h5 class="timeline-card__title">${esc(x.author || "Автор")}</h5><span class="badge badge--muted">${esc(s.formatDateTime(x.created_at))}</span></div>
                <p class="timeline-card__meta">${esc(x.body || "-")}</p>
            </article>`).join("") : '<div class="empty-state">Комментариев пока нет.</div>';
    }

    function renderStatuses() {
        q("status-list").innerHTML = st.statuses.length ? st.statuses.map((x) => `
            <article class="application-card ${x.id === st.statusID ? "application-card--active" : ""}" data-status-id="${x.id}">
                <div class="application-card__head"><h3 class="application-card__title">${esc(x.name)}</h3><span class="${badge(x.code, x.is_final)}">${esc(x.code)}</span></div>
                <p class="application-card__meta">${x.is_final ? "Финальный статус" : "Рабочий статус"}</p>
            </article>`).join("") : '<div class="empty-state">Статусы пока не загружены.</div>';
        q("status-list").querySelectorAll("[data-status-id]").forEach((node) => node.addEventListener("click", () => {
            const x = st.statuses.find((item) => Number(item.id) === Number(node.dataset.statusId));
            if (!x) return;
            st.statusID = x.id;
            q("status-code").value = x.code || "";
            q("status-name").value = x.name || "";
            q("status-final").checked = Boolean(x.is_final);
            renderStatuses();
        }));
    }

    function clearApp() {
        st.appID = 0;
        st.app = null;
        st.history = [];
        q("application-placeholder").classList.remove("is-hidden");
        q("application-details").classList.add("is-hidden");
        renderComments([]);
        renderHistory();
    }

    function showApp(x) {
        st.appID = x.id;
        st.app = x;
        q("application-placeholder").classList.add("is-hidden");
        q("application-details").classList.remove("is-hidden");
        q("application-title").textContent = x.title || "-";
        q("application-description").textContent = x.description || "-";
        q("application-status").textContent = x.status?.name || "-";
        q("application-status").className = badge(x.status?.code, findStatus(x.status?.code)?.is_final);
        q("application-priority").textContent = x.priority?.name || "-";
        q("application-priority").className = "badge badge--soft";
        q("application-created-by").textContent = x.created_by?.name || "-";
        q("application-assigned-to").textContent = x.assigned_to?.name || "Не назначен";
        q("application-created-at").textContent = s.formatDateTime(x.created_at);
        q("application-updated-at").textContent = s.formatDateTime(x.updated_at);
        q("application-phone").textContent = x.contact_phone || "Не указан";
        q("application-address").textContent = x.contact_address || "Не указан";
        fillOperators();
        fillStatuses();
        if (x.assigned_to) q("assign-operator").value = String(x.assigned_to.id);
        if (x.status?.code) q("application-status-select").value = x.status.code;
        q("application-action-comment").value = "";
        renderComments(Array.isArray(x.comments) ? x.comments : []);
        renderHistory();
        renderAppsTable();
    }

    function clearUser() {
        st.userID = 0;
        st.user = null;
        q("user-placeholder").classList.remove("is-hidden");
        q("user-details").classList.add("is-hidden");
    }

    function showUser(x) {
        st.userID = x.id;
        st.user = x;
        q("user-placeholder").classList.add("is-hidden");
        q("user-details").classList.remove("is-hidden");
        q("user-name").textContent = x.full_name || "-";
        q("user-email").textContent = x.email || "-";
        q("user-id").textContent = String(x.id);
        q("user-created-at").textContent = s.formatDateTime(x.created_at);
        q("user-role-name").textContent = x.role_name || x.role_code || "-";
        q("user-role-badge").textContent = x.role_name || x.role_code || "-";
        q("user-role-badge").className = x.role_code === "admin" ? "badge badge--soft" : "badge";
        q("user-active-badge").textContent = x.is_active ? "Активен" : "Отключён";
        q("user-active-badge").className = x.is_active ? "badge badge--success" : "badge badge--danger";
        q("user-role-select").value = x.role_code;
        renderUsersTable();
    }

    async function loadApps(pref) {
        msg("application-list-message", "", "");
        try {
            const r = await s.authorizedFetch("/api/admin/applications/get-apps");
            const d = await json(r);
            if (!r.ok) {
                msg("application-list-message", "error", d.message || "Не удалось загрузить заявки.");
                st.apps = [];
                renderAppsTable();
                return;
            }
            st.apps = Array.isArray(d.items) ? d.items : [];
            syncStatuses();
            fillStatuses();
            renderStats();
            renderStatuses();
            renderAppsTable();
            if (!st.apps.length) return clearApp();
            const id = pref || st.appID || st.apps[0].id;
            await loadApp(st.apps.some((x) => x.id === id) ? id : st.apps[0].id);
        } catch {
            msg("application-list-message", "error", "Не удалось загрузить заявки.");
        }
    }

    async function loadApp(id) {
        if (!id) return clearApp();
        msg("application-detail-message", "", "");
        try {
            const [r1, r2] = await Promise.all([
                s.authorizedFetch(`/api/admin/applications/get-app?id=${id}`),
                s.authorizedFetch(`/api/admin/applications/get-history?id=${id}`),
            ]);
            const d1 = await json(r1);
            const d2 = await json(r2);
            if (r1.status === 404) {
                st.apps = st.apps.filter((x) => x.id !== id);
                renderAppsTable();
                return st.apps.length ? loadApp(st.apps[0].id) : clearApp();
            }
            if (!r1.ok) return msg("application-detail-message", "error", d1.message || "Не удалось загрузить карточку заявки.");
            st.history = Array.isArray(d2.items) ? d2.items : [];
            showApp(d1);
        } catch {
            msg("application-detail-message", "error", "Не удалось загрузить карточку заявки.");
        }
    }

    async function loadUsers(pref) {
        msg("user-list-message", "", "");
        try {
            const r = await s.authorizedFetch("/api/admin/users/get-users");
            const d = await json(r);
            if (!r.ok) {
                msg("user-list-message", "error", d.message || "Не удалось загрузить пользователей.");
                st.users = [];
                renderUsersTable();
                return;
            }
            st.users = Array.isArray(d.items) ? d.items : [];
            renderStats();
            fillOperators();
            renderUsersTable();
            if (pref && st.users.some((x) => x.id === pref)) await loadUser(pref);
        } catch {
            msg("user-list-message", "error", "Не удалось загрузить пользователей.");
        }
    }

    async function loadUser(id) {
        if (!id) return clearUser();
        msg("user-detail-message", "", "");
        try {
            const r = await s.authorizedFetch(`/api/admin/users/get-user?id=${id}`);
            const d = await json(r);
            if (r.status === 404) {
                st.users = st.users.filter((x) => x.id !== id);
                renderUsersTable();
                return clearUser();
            }
            if (!r.ok) return msg("user-detail-message", "error", d.message || "Не удалось загрузить пользователя.");
            showUser(d);
        } catch {
            msg("user-detail-message", "error", "Не удалось загрузить пользователя.");
        }
    }

    async function send(url, body, method, errTarget, okText, cb) {
        try {
            const r = await s.authorizedFetch(url, { method, body: JSON.stringify(body) });
            const d = await json(r);
            if (!r.ok) return msg(errTarget, "error", d.message || "Не удалось выполнить действие.");
            if (okText) msg(errTarget, "success", okText);
            if (cb) await cb(d);
        } catch {
            msg(errTarget, "error", "Не удалось выполнить действие.");
        }
    }

    q("application-page-prev").addEventListener("click", () => { st.appPage -= 1; renderAppsTable(); });
    q("application-page-next").addEventListener("click", () => { st.appPage += 1; renderAppsTable(); });
    q("user-page-prev").addEventListener("click", () => { st.userPage -= 1; renderUsersTable(); });
    q("user-page-next").addEventListener("click", () => { st.userPage += 1; renderUsersTable(); });
    q("application-search").addEventListener("input", () => { st.appSearch = q("application-search").value; st.appPage = 1; renderAppsTable(); });
    q("user-search").addEventListener("input", () => { st.userSearch = q("user-search").value; st.userPage = 1; renderUsersTable(); });
    q("history-search").addEventListener("input", () => { st.historySearch = q("history-search").value.toLowerCase(); renderHistory(); });
    q("application-filter-bar").querySelectorAll("[data-filter]").forEach((btn) => btn.addEventListener("click", () => {
        st.appFilter = btn.dataset.filter || "all";
        st.appPage = 1;
        q("application-filter-bar").querySelectorAll("[data-filter]").forEach((x) => x.classList.toggle("filter-chip--active", x === btn));
        renderAppsTable();
    }));

    q("assign-button").addEventListener("click", () => {
        if (!st.app || !q("assign-operator").value) return;
        send("/api/admin/applications/assign-app", { id: st.app.id, operator_id: Number(q("assign-operator").value) }, "PUT", "application-detail-message", "Оператор назначен.", async () => loadApps(st.app.id));
    });
    q("change-status-button").addEventListener("click", () => {
        if (!st.app || !q("application-status-select").value) return;
        send("/api/admin/applications/change-status", { id: st.app.id, status_code: q("application-status-select").value, comment: q("application-action-comment").value.trim() || undefined }, "PUT", "application-detail-message", "Статус обновлён.", async () => loadApps(st.app.id));
    });
    q("delete-application-button").addEventListener("click", () => {
        if (!st.app || !window.confirm("Удалить выбранную заявку?")) return;
        send("/api/admin/applications/delete-app", { id: st.app.id }, "DELETE", "application-detail-message", "", async () => {
            st.apps = st.apps.filter((x) => x.id !== st.app.id);
            renderStats();
            renderAppsTable();
            st.apps.length ? await loadApp(st.apps[0].id) : clearApp();
        });
    });
    q("change-role-button").addEventListener("click", () => {
        if (!st.user) return;
        send("/api/admin/users/change-role", { user_id: st.user.id, role_code: q("user-role-select").value }, "PUT", "user-detail-message", "Роль изменена.", async () => {
            await loadUsers(st.user.id);
            await loadApps(st.appID);
        });
    });
    q("delete-user-button").addEventListener("click", () => {
        if (!st.user || !window.confirm("Удалить выбранного пользователя?")) return;
        send("/api/admin/users/delete-user", { user_id: st.user.id }, "DELETE", "user-detail-message", "", async () => {
            st.users = st.users.filter((x) => x.id !== st.user.id);
            renderStats();
            renderUsersTable();
            fillOperators();
            clearUser();
        });
    });
    q("create-status-button").addEventListener("click", () => {
        const code = q("status-code").value.trim();
        const name = q("status-name").value.trim();
        if (!code || !name) return msg("status-message", "error", "Заполните code и название статуса.");
        send("/api/admin/dictionaries/create-status", { code, name, is_final: q("status-final").checked }, "POST", "status-message", "Статус создан.", async (d) => {
            st.statuses = [...st.statuses.filter((x) => x.code !== d.code), d];
            syncStatuses(); fillStatuses(); renderStatuses(); q("status-code").value = ""; q("status-name").value = ""; q("status-final").checked = false; st.statusID = 0;
        });
    });
    q("update-status-button").addEventListener("click", () => {
        const code = q("status-code").value.trim();
        const name = q("status-name").value.trim();
        if (!st.statusID) return msg("status-message", "error", "Сначала выберите статус из списка.");
        if (!code || !name) return msg("status-message", "error", "Заполните code и название статуса.");
        send("/api/admin/dictionaries/change-status", { id: st.statusID, code, name, is_final: q("status-final").checked }, "PUT", "status-message", "Статус обновлён.", async (d) => {
            st.statuses = st.statuses.map((x) => x.id === d.id ? d : x);
            syncStatuses(); fillStatuses(); renderStatuses(); q("status-code").value = ""; q("status-name").value = ""; q("status-final").checked = false; st.statusID = 0;
        });
    });
    q("delete-status-button").addEventListener("click", () => {
        if (!st.statusID || !window.confirm("Удалить выбранный статус?")) return;
        send("/api/admin/dictionaries/delete-status", { id: st.statusID }, "DELETE", "status-message", "Статус удалён.", async () => {
            st.statuses = st.statuses.filter((x) => Number(x.id) !== Number(st.statusID));
            syncStatuses(); fillStatuses(); renderStatuses(); q("status-code").value = ""; q("status-name").value = ""; q("status-final").checked = false; st.statusID = 0;
        });
    });
    q("reset-status-button").addEventListener("click", () => {
        st.statusID = 0;
        q("status-code").value = "";
        q("status-name").value = "";
        q("status-final").checked = false;
        renderStatuses();
    });

    Promise.all([loadUsers(), loadApps()]).then(() => {
        renderStats();
        renderStatuses();
        fillOperators();
        fillStatuses();
    });
})();
