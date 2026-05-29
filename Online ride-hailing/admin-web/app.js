const DEMO_DELAY = 120;

const state = {
  current: "dashboard",
  keyword: "",
  page: 1,
  pageSize: 10,
  apiOnline: false,
  theme: localStorage.getItem("admin_theme") || "blue",
  admin: {
    id: Number(localStorage.getItem("admin_id")) || 1,
    username: localStorage.getItem("admin_username") || "admin",
    role: "超级管理员",
  },
};

const navGroups = [
  {
    title: "运营中心",
    items: [
      { key: "dashboard", icon: "□", label: "运营概览" },
      { key: "orders", icon: "≡", label: "订单管理", count: "86" },
      { key: "dispatch", icon: "+", label: "后台派单" },
    ],
  },
  {
    title: "系统管理",
    items: [
      { key: "admins", icon: "◎", label: "管理员" },
      { key: "roles", icon: "◇", label: "角色权限" },
      { key: "menus", icon: "☰", label: "菜单管理" },
    ],
  },
  {
    title: "基础配置",
    items: [
      { key: "cities", icon: "⌂", label: "城市配置" },
      { key: "carTypes", icon: "▣", label: "车型配置" },
      { key: "settings", icon: "⚙", label: "接口设置" },
    ],
  },
];

const pages = {
  dashboard: { module: "运营中心", title: "运营概览" },
  orders: { module: "运营中心", title: "订单管理" },
  dispatch: { module: "运营中心", title: "后台派单" },
  admins: { module: "系统管理", title: "管理员" },
  roles: { module: "系统管理", title: "角色权限" },
  menus: { module: "系统管理", title: "菜单管理" },
  cities: { module: "基础配置", title: "城市配置" },
  carTypes: { module: "基础配置", title: "车型配置" },
  settings: { module: "基础配置", title: "接口设置" },
};

const tableConfigs = {
  orders: {
    endpoint: "/order/list",
    title: "订单列表",
    searchPlaceholder: "订单号、乘客、手机号",
    filters: [
      { key: "status", label: "全部状态", options: orderStatusOptions() },
    ],
    columns: [
      { key: "order_id", label: "订单号", width: "150px" },
      { key: "passenger_name", label: "乘客" },
      { key: "passenger_phone", label: "乘客手机", width: "126px" },
      { key: "driver_name", label: "司机" },
      { key: "car_type", label: "车型" },
      { key: "route", label: "行程", render: routeCell },
      { key: "status", label: "状态", width: "92px", render: orderStatusTag },
      { key: "final_price", label: "金额", width: "88px", render: moneyCell },
      { key: "created_at", label: "创建时间", width: "138px", render: timeCell },
    ],
    actions: [
      { label: "详情", type: "text", action: showOrderDetail },
      { label: "派单", type: "secondary", action: () => openDispatchDrawer() },
    ],
    demo: demoOrders,
  },
  admins: {
    endpoint: "/admin/list",
    title: "管理员列表",
    searchPlaceholder: "用户名、角色",
    createLabel: "新增管理员",
    create: () => openEntityDrawer("admins"),
    columns: [
      { key: "id", label: "ID", width: "70px" },
      { key: "username", label: "账号" },
      { key: "role_name", label: "角色" },
      { key: "status", label: "状态", render: normalStatusTag },
      { key: "last_login_time", label: "最近登录", render: timeCell },
    ],
    actions: [
      { label: "编辑", type: "text", action: (row) => openEntityDrawer("admins", row) },
      { label: "删除", type: "danger", action: (row) => deleteEntity("管理员", "/admin/delete", { admin_id: row.id }) },
    ],
    demo: demoAdmins,
  },
  roles: {
    endpoint: "/role/list",
    title: "角色列表",
    searchPlaceholder: "角色名称",
    createLabel: "新增角色",
    create: () => openEntityDrawer("roles"),
    columns: [
      { key: "id", label: "ID", width: "70px" },
      { key: "name", label: "角色名称" },
      { key: "remark", label: "备注" },
      { key: "menu_ids", label: "菜单数量", render: (row) => Array.isArray(row.menu_ids) ? row.menu_ids.length : row.menu_ids || 0 },
    ],
    actions: [
      { label: "授权", type: "secondary", action: (row) => openEntityDrawer("roles", row) },
      { label: "删除", type: "danger", action: (row) => deleteEntity("角色", "/role/delete", { role_id: row.id }) },
    ],
    demo: demoRoles,
  },
  menus: {
    endpoint: "/menu/list",
    title: "菜单树",
    searchPlaceholder: "菜单名称、路径",
    createLabel: "新增菜单",
    create: () => openEntityDrawer("menus"),
    columns: [
      { key: "name", label: "菜单名称", render: menuNameCell },
      { key: "path", label: "路径" },
      { key: "icon", label: "图标" },
      { key: "sort", label: "排序", width: "80px" },
    ],
    actions: [
      { label: "编辑", type: "text", action: (row) => openEntityDrawer("menus", row) },
      { label: "删除", type: "danger", action: (row) => deleteEntity("菜单", "/menu/delete", { menu_id: row.id }) },
    ],
    demo: demoMenus,
    tree: true,
  },
  cities: {
    endpoint: "/city/list",
    title: "城市配置",
    searchPlaceholder: "城市名称、编码",
    createLabel: "新增城市",
    create: () => openEntityDrawer("cities"),
    columns: [
      { key: "id", label: "ID", width: "70px" },
      { key: "city_code", label: "城市编码" },
      { key: "city_name", label: "城市名称" },
      { key: "status", label: "状态", render: enableStatusTag },
    ],
    actions: [
      { label: "编辑", type: "text", action: (row) => openEntityDrawer("cities", row) },
      { label: "删除", type: "danger", action: (row) => deleteEntity("城市", "/city/delete", { id: row.id }) },
    ],
    demo: demoCities,
  },
  carTypes: {
    endpoint: "/cartype/list",
    title: "车型配置",
    searchPlaceholder: "车型名称",
    createLabel: "新增车型",
    create: () => openEntityDrawer("carTypes"),
    columns: [
      { key: "id", label: "ID", width: "70px" },
      { key: "type_name", label: "车型" },
      { key: "base_price", label: "起步价", render: moneyCell },
      { key: "km_price", label: "公里价", render: moneyCell },
      { key: "minute_price", label: "时长价", render: moneyCell },
      { key: "status", label: "状态", render: enableStatusTag },
    ],
    actions: [
      { label: "编辑", type: "text", action: (row) => openEntityDrawer("carTypes", row) },
      { label: "删除", type: "danger", action: (row) => deleteEntity("车型", "/cartype/delete", { id: row.id }) },
    ],
    demo: demoCarTypes,
  },
};

const formConfigs = {
  admins: {
    title: "管理员",
    createEndpoint: "/admin/create",
    updateEndpoint: "/admin/update",
    idKey: "admin_id",
    fields: [
      { key: "username", label: "账号", required: true },
      { key: "password", label: "密码", type: "password", createOnly: true, required: true },
      { key: "role_id", label: "角色", type: "select", options: [{ value: 1, label: "超级管理员" }, { value: 2, label: "运营管理员" }, { value: 3, label: "客服专员" }] },
      { key: "status", label: "状态", type: "select", options: [{ value: 1, label: "正常" }, { value: 0, label: "冻结" }] },
    ],
  },
  roles: {
    title: "角色",
    createEndpoint: "/role/create",
    updateEndpoint: "/role/update",
    idKey: "role_id",
    fields: [
      { key: "name", label: "角色名称", required: true },
      { key: "remark", label: "备注" },
      { key: "menu_ids", label: "菜单 ID", hint: "多个 ID 用英文逗号分隔，例如 1,2,3" },
    ],
    normalize(payload) {
      payload.menu_ids = String(payload.menu_ids || "")
        .split(",")
        .map((item) => Number(item.trim()))
        .filter(Boolean);
      return payload;
    },
  },
  menus: {
    title: "菜单",
    createEndpoint: "/menu/create",
    updateEndpoint: "/menu/update",
    idKey: "menu_id",
    fields: [
      { key: "parent_id", label: "父级 ID", type: "number", defaultValue: 0 },
      { key: "name", label: "菜单名称", required: true },
      { key: "path", label: "路由路径", required: true },
      { key: "icon", label: "图标" },
      { key: "sort", label: "排序", type: "number", defaultValue: 1 },
    ],
  },
  cities: {
    title: "城市",
    createEndpoint: "/city/create",
    updateEndpoint: "/city/update",
    idKey: "id",
    fields: [
      { key: "city_code", label: "城市编码", required: true },
      { key: "city_name", label: "城市名称", required: true },
      { key: "status", label: "状态", type: "select", options: [{ value: 1, label: "启用" }, { value: 0, label: "禁用" }] },
    ],
  },
  carTypes: {
    title: "车型",
    createEndpoint: "/cartype/create",
    updateEndpoint: "/cartype/update",
    idKey: "id",
    fields: [
      { key: "type_name", label: "车型名称", required: true },
      { key: "base_price", label: "起步价", type: "number", step: "0.01" },
      { key: "km_price", label: "公里单价", type: "number", step: "0.01" },
      { key: "minute_price", label: "时长单价", type: "number", step: "0.01" },
      { key: "status", label: "状态", type: "select", options: [{ value: 1, label: "启用" }, { value: 0, label: "禁用" }] },
    ],
  },
};

const els = {
  loginScreen: document.getElementById("loginScreen"),
  loginForm: document.getElementById("loginForm"),
  appShell: document.getElementById("appShell"),
  mainNav: document.getElementById("mainNav"),
  pageModule: document.getElementById("pageModule"),
  pageTitle: document.getElementById("pageTitle"),
  content: document.getElementById("content"),
  globalSearch: document.getElementById("globalSearch"),
  refreshBtn: document.getElementById("refreshBtn"),
  apiModeTag: document.getElementById("apiModeTag"),
  adminAvatar: document.getElementById("adminAvatar"),
  adminName: document.getElementById("adminName"),
  adminRole: document.getElementById("adminRole"),
  themeToggleBtn: document.getElementById("themeToggleBtn"),
  themeToggleText: document.getElementById("themeToggleText"),
  drawerMask: document.getElementById("drawerMask"),
  drawer: document.getElementById("drawer"),
  drawerType: document.getElementById("drawerType"),
  drawerTitle: document.getElementById("drawerTitle"),
  drawerForm: document.getElementById("drawerForm"),
  closeDrawerBtn: document.getElementById("closeDrawerBtn"),
  openSettingBtn: document.getElementById("openSettingBtn"),
  toast: document.getElementById("toast"),
};

boot();

function boot() {
  applyTheme();
  renderNav();
  bindEvents();

  if (localStorage.getItem("admin_token")) {
    enterApp();
  }
}

function bindEvents() {
  els.loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    const username = document.getElementById("loginUsername").value.trim();
    const password = document.getElementById("loginPassword").value.trim();
    const resp = await AdminApi.login({ username, password });

    if (resp.ok) {
      state.apiOnline = true;
      state.admin.username = resp.data.username || username;
      state.admin.id = Number(resp.data.admin_id) || 1;
      toast("登录成功，已连接本地接口");
    } else {
      state.apiOnline = false;
      localStorage.setItem("admin_token", "demo-token");
      localStorage.setItem("admin_username", username || "admin");
      toast("本地接口未连接，已进入演示模式");
    }

    enterApp();
  });

  els.globalSearch.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      state.keyword = event.currentTarget.value.trim();
      state.page = 1;
      renderCurrent();
    }
  });

  els.refreshBtn.addEventListener("click", () => renderCurrent());
  els.themeToggleBtn.addEventListener("click", toggleTheme);
  els.closeDrawerBtn.addEventListener("click", closeDrawer);
  els.drawerMask.addEventListener("click", closeDrawer);
  els.openSettingBtn.addEventListener("click", () => navigate("settings"));
}

function applyTheme() {
  document.documentElement.dataset.theme = state.theme === "slate" ? "slate" : "blue";
  els.themeToggleText.textContent = state.theme === "slate" ? "石墨" : "蓝白";
}

function toggleTheme() {
  state.theme = state.theme === "slate" ? "blue" : "slate";
  localStorage.setItem("admin_theme", state.theme);
  applyTheme();
  toast(`已切换为${state.theme === "slate" ? "石墨青蓝" : "清爽蓝白"}主题`);
}

function enterApp() {
  els.loginScreen.classList.add("hidden");
  els.appShell.classList.remove("hidden");
  state.admin.username = localStorage.getItem("admin_username") || state.admin.username;
  updateAdminChip();
  renderCurrent();
}

function updateAdminChip() {
  els.adminName.textContent = state.admin.username;
  els.adminRole.textContent = state.admin.role;
  els.adminAvatar.textContent = state.admin.username.slice(0, 1).toUpperCase();
  els.apiModeTag.textContent = state.apiOnline ? "接口在线" : "演示数据";
}

function renderNav() {
  els.mainNav.innerHTML = navGroups.map((group) => `
    <section class="nav-section">
      <div class="nav-section-title">${group.title}</div>
      ${group.items.map((item) => `
        <button class="nav-item" data-nav="${item.key}" type="button">
          <span class="nav-icon">${item.icon}</span>
          <span>${item.label}</span>
          <small class="nav-count">${item.count || ""}</small>
        </button>
      `).join("")}
    </section>
  `).join("");

  els.mainNav.addEventListener("click", (event) => {
    const button = event.target.closest("[data-nav]");
    if (!button) return;
    navigate(button.dataset.nav);
  });
}

function navigate(key) {
  state.current = key;
  state.page = 1;
  state.keyword = "";
  els.globalSearch.value = "";
  renderCurrent();
}

function renderCurrent() {
  const page = pages[state.current];
  els.pageModule.textContent = page.module;
  els.pageTitle.textContent = page.title;

  document.querySelectorAll("[data-nav]").forEach((button) => {
    button.classList.toggle("active", button.dataset.nav === state.current);
  });

  if (state.current === "dashboard") return renderDashboard();
  if (state.current === "dispatch") return renderDispatch();
  if (state.current === "settings") return renderSettings();
  return renderTablePage(state.current);
}

function renderDashboard() {
  const recentOrdersConfig = { ...tableConfigs.orders, actions: null };
  els.content.innerHTML = `
    <section class="metric-grid">
      ${metricCard("今日订单", "328", "+12.4%", "≡")}
      ${metricCard("成交金额", "¥46,280", "+8.7%", "￥")}
      ${metricCard("在线司机", "124", "+5.1%", "◎")}
      ${metricCard("售后待办", "12", "-2", "!")}
    </section>
    <section class="dashboard-grid">
      <div class="panel-card">
        <h3>订单状态分布</h3>
        <div class="status-list">
          ${statusLine("待接单", 42)}
          ${statusLine("司机已接单", 68)}
          ${statusLine("行程中", 56)}
          ${statusLine("已完成", 84)}
          ${statusLine("已取消", 16)}
        </div>
      </div>
      <div class="panel-card">
        <h3>运营待办</h3>
        <div class="todo-list">
          ${todoItem("售后工单待审核", "12")}
          ${todoItem("城市价格配置待确认", "3")}
          ${todoItem("新增管理员待授权", "2")}
          ${todoItem("异常支付单", "5")}
        </div>
      </div>
    </section>
    <section class="table-panel" style="margin-top:18px">
      <div class="table-head">
        <h3>最近订单</h3>
        <button class="text-btn" type="button" data-jump-orders>查看全部</button>
      </div>
      ${buildTable(recentOrdersConfig, demoOrders().slice(0, 5))}
    </section>
  `;
  els.content.querySelector("[data-jump-orders]").addEventListener("click", () => navigate("orders"));
}

async function renderTablePage(key) {
  const config = tableConfigs[key];
  els.content.innerHTML = `
    <section class="toolbar">
      <div class="filters">
        <input id="pageKeyword" type="search" placeholder="${config.searchPlaceholder || "请输入关键字"}" value="${escapeAttr(state.keyword)}">
        ${(config.filters || []).map(renderFilter).join("")}
        <button class="ghost-btn" id="queryBtn" type="button">查询</button>
      </div>
      <div class="toolbar-actions">
        <button class="ghost-btn" id="resetBtn" type="button">重置</button>
        ${config.createLabel ? `<button class="primary-btn" id="createBtn" type="button">${config.createLabel}</button>` : ""}
      </div>
    </section>
    <section class="table-panel">
      <div class="table-head">
        <h3>${config.title}</h3>
        <span id="tableSummary">正在加载</span>
      </div>
      <div id="tableWrap">${emptyState("正在读取数据")}</div>
    </section>
  `;

  els.content.querySelector("#queryBtn").addEventListener("click", () => {
    state.keyword = els.content.querySelector("#pageKeyword").value.trim();
    state.page = 1;
    renderTablePage(key);
  });
  els.content.querySelector("#resetBtn").addEventListener("click", () => {
    state.keyword = "";
    state.page = 1;
    renderTablePage(key);
  });
  const createBtn = els.content.querySelector("#createBtn");
  if (createBtn) createBtn.addEventListener("click", config.create);

  const payload = collectFilters({ page: state.page, page_size: state.pageSize, keyword: state.keyword });
  const { list, total, online } = await fetchList(config, payload);
  state.apiOnline = online;
  updateAdminChip();

  const rows = config.tree ? flattenTree(list) : list;
  els.content.querySelector("#tableSummary").textContent = `共 ${total} 条`;
  els.content.querySelector("#tableWrap").innerHTML = buildTable(config, rows, total);
  bindTableActions(config, rows);
}

function renderDispatch() {
  els.content.innerHTML = `
    <section class="dashboard-grid">
      <div class="panel-card">
        <h3>后台派单</h3>
        <p>为客服或运营人员创建即时订单，并指定司机或进入调度池。</p>
        <form class="drawer-form" id="dispatchForm" style="padding:18px 0 0">
          ${dispatchFields()}
          <div class="drawer-actions">
            <button class="ghost-btn" type="reset">清空</button>
            <button class="primary-btn" type="submit">提交派单</button>
          </div>
        </form>
      </div>
      <div class="panel-card">
        <h3>派单建议</h3>
        <div class="todo-list">
          ${todoItem("建议选择舒适型，附近可用司机较多", "6 位")}
          ${todoItem("当前城市高峰计费已开启", "广州")}
          ${todoItem("预计接驾距离", "2.1 km")}
        </div>
      </div>
    </section>
  `;

  els.content.querySelector("#dispatchForm").addEventListener("submit", submitDispatch);
}

function renderSettings() {
  els.content.innerHTML = `
    <section class="panel-card">
      <h3>接口连接</h3>
      <p>管理端默认连接 admin-api：${AdminApi.getApiBaseUrl()}。你可以在这里切换后端地址。</p>
      <div class="setting-block">
        <label class="field">
          <span class="field-label">Admin API Base URL</span>
          <input class="setting-input" id="apiBaseInput" value="${escapeAttr(AdminApi.getApiBaseUrl())}">
        </label>
        <div>
          <button class="primary-btn" id="saveApiBtn" type="button">保存接口地址</button>
          <button class="ghost-btn" id="demoModeBtn" type="button">${state.apiOnline ? "使用演示数据" : "使用真实数据"}</button>
        </div>
      </div>
    </section>
  `;
  els.content.querySelector("#saveApiBtn").addEventListener("click", () => {
    AdminApi.setApiBaseUrl(els.content.querySelector("#apiBaseInput").value);
    toast("接口地址已保存");
  });
  els.content.querySelector("#demoModeBtn").addEventListener("click", () => {
    state.apiOnline = !state.apiOnline;
    updateAdminChip();
    // 刷新按钮文案和页面
    els.content.querySelector("#demoModeBtn").textContent = state.apiOnline ? "使用演示数据" : "使用真实数据";
    toast(state.apiOnline ? "已切换为真实数据模式" : "已切换为演示数据模式");
    if (state.apiOnline) renderCurrent();
  });
}

async function fetchList(config, payload) {
  const resp = await AdminApi.list(config.endpoint, payload);
  if (resp.ok) {
    const data = resp.data || {};
    const list = data.list || (Array.isArray(data) ? data : []);
    return {
      list,
      total: Number(data.total) || list.length,
      online: true,
    };
  }

  await wait(DEMO_DELAY);
  const demo = config.demo();
  // 树形结构不参与分页，只做关键字过滤后整体返回
  if (config.tree) {
    const filtered = filterRows(demo, payload.keyword);
    return { list: filtered, total: filtered.length, online: false };
  }
  // 普通列表模式：关键字过滤 + 客户端分页
  const filtered = filterRows(demo, payload.keyword);
  const total = filtered.length;
  const offset = (Number(payload.page) - 1) * Number(payload.page_size);
  const list = filtered.slice(offset, offset + Number(payload.page_size));
  return { list, total, online: false };
}

function buildTable(config, rows, total = rows.length) {
  if (!rows.length) return emptyState("暂无数据");
  const totalPages = Math.max(1, Math.ceil(total / state.pageSize));
  return `
    <table class="data-table">
      <thead>
        <tr>
          ${config.columns.map((col) => `<th style="${col.width ? `width:${col.width}` : ""}">${col.label}</th>`).join("")}
          ${config.actions ? `<th style="width:${Math.max(config.actions.length * 74, 96)}px">操作</th>` : ""}
        </tr>
      </thead>
      <tbody>
        ${rows.map((row, index) => `
          <tr>
            ${config.columns.map((col) => `<td title="${escapeAttr(rawCell(row, col))}">${cell(row, col)}</td>`).join("")}
            ${config.actions ? `
              <td>
                ${config.actions.map((action, actionIndex) => actionButton(action, index, actionIndex)).join("")}
              </td>` : ""}
          </tr>
        `).join("")}
      </tbody>
    </table>
    <div class="pagination">
      <button class="page-btn" data-page="first" ${state.page <= 1 ? "disabled" : ""} type="button">«</button>
      <button class="page-btn" data-page="prev" ${state.page <= 1 ? "disabled" : ""} type="button">‹</button>
      ${renderPageNumbers(state.page, totalPages)}
      <button class="page-btn" data-page="next" ${state.page >= totalPages ? "disabled" : ""} type="button">›</button>
      <button class="page-btn" data-page="last" ${state.page >= totalPages ? "disabled" : ""} type="button">»</button>
      <span class="page-info">第 ${state.page} / ${totalPages} 页，共 ${total} 条</span>
    </div>
  `;
}

function bindTableActions(config, rows) {
  els.content.querySelectorAll("[data-row-action]").forEach((button) => {
    button.addEventListener("click", () => {
      const row = rows[Number(button.dataset.row)];
      const action = config.actions[Number(button.dataset.action)];
      action.action(row);
    });
  });

  // 分页按钮事件
  els.content.querySelectorAll(".page-btn").forEach((button) => {
    button.addEventListener("click", () => {
      const dir = button.dataset.page;
      const totalPages = Math.max(1, Math.ceil(Number(els.content.querySelector("#tableSummary").textContent.replace("共 ", "").replace(" 条", "")) / state.pageSize));
      if (dir === "first") state.page = 1;
      else if (dir === "prev" && state.page > 1) state.page--;
      else if (dir === "next" && state.page < totalPages) state.page++;
      else if (dir === "last") state.page = totalPages;
      else if (dir === "goto") state.page = Number(button.dataset.goto);
      renderTablePage(state.current);
    });
  });
}

function renderPageNumbers(current, total) {
  let html = "";
  const maxVisible = 5;
  let start = Math.max(1, current - Math.floor(maxVisible / 2));
  let end = Math.min(total, start + maxVisible - 1);
  if (end - start + 1 < maxVisible) start = Math.max(1, end - maxVisible + 1);
  for (let i = start; i <= end; i++) {
    html += `<button class="page-btn ${i === current ? "current" : ""}" data-page="goto" data-goto="${i}" type="button">${i}</button>`;
  }
  return html;
}

function actionButton(action, rowIndex, actionIndex) {
  const className = action.type === "danger" ? "danger-btn" : action.type === "secondary" ? "secondary-btn" : "text-btn";
  return `<button class="${className}" data-row-action data-row="${rowIndex}" data-action="${actionIndex}" type="button">${action.label}</button>`;
}

function openEntityDrawer(key, row = null) {
  const config = formConfigs[key];
  const isEdit = Boolean(row);
  openDrawer(`${isEdit ? "编辑" : "新增"}${config.title}`, isEdit ? "编辑" : "新增");

  els.drawerForm.innerHTML = config.fields
    .filter((field) => !field.createOnly || !isEdit)
    .map((field) => fieldHtml(field, row))
    .join("") + drawerActionsHtml();

  els.drawerForm.onsubmit = async (event) => {
    event.preventDefault();
    let payload = collectForm(els.drawerForm, config.fields);
    if (row) payload[config.idKey] = row.id;
    if (config.normalize) payload = config.normalize(payload);

    const endpoint = isEdit ? config.updateEndpoint : config.createEndpoint;
    const resp = await AdminApi.mutate(endpoint, payload);
    if (!resp.ok) {
      toast("接口未连接，已在演示模式中保存动作");
    } else {
      state.apiOnline = true;
      toast("保存成功");
    }
    closeDrawer();
    renderCurrent();
  };
}

function openDispatchDrawer(row) {
  openDrawer("后台派单", "新增");
  els.drawerForm.innerHTML = dispatchFields(row) + drawerActionsHtml("提交派单");
  els.drawerForm.onsubmit = submitDispatch;
}

async function submitDispatch(event) {
  event.preventDefault();
  const payload = collectForm(event.currentTarget);
  payload.car_type = Number(payload.car_type || 1);
  payload.driver_id = Number(payload.driver_id || 0);
  payload.start_lng = Number(payload.start_lng || 113.2644);
  payload.start_lat = Number(payload.start_lat || 23.1291);
  payload.end_lng = Number(payload.end_lng || 113.3245);
  payload.end_lat = Number(payload.end_lat || 23.1065);
  payload.estimated_price = Number(payload.estimated_price || 0);

  const resp = await AdminApi.mutate("/order/dispatch", payload);
  toast(resp.ok ? "派单已提交" : "接口未连接，已模拟提交派单");
  closeDrawer();
}

function showOrderDetail(row) {
  openDrawer("订单详情", "查看");
  els.drawerForm.innerHTML = `
    <div class="setting-block">
      <strong>${row.order_id}</strong>
      <span>${row.passenger_name} ${row.passenger_phone}</span>
      <span>${row.start_address} → ${row.end_address}</span>
      <span>司机：${row.driver_name || "待分配"} ${row.driver_phone || ""}</span>
      <span>金额：${moneyCell(row, { key: "final_price" })}</span>
    </div>
    <div class="drawer-actions">
      <button class="ghost-btn" type="button" data-close>关闭</button>
      <button class="primary-btn" type="button" data-dispatch>重新派单</button>
    </div>
  `;
  els.drawerForm.querySelector("[data-close]").addEventListener("click", closeDrawer);
  els.drawerForm.querySelector("[data-dispatch]").addEventListener("click", () => openDispatchDrawer(row));
}

async function deleteEntity(name, endpoint, payload) {
  if (!confirm(`确认删除该${name}？`)) return;
  const resp = await AdminApi.mutate(endpoint, payload);
  toast(resp.ok ? "删除成功" : "接口未连接，已模拟删除动作");
  renderCurrent();
}

function openDrawer(title, type) {
  els.drawerTitle.textContent = title;
  els.drawerType.textContent = type;
  els.drawer.classList.remove("hidden");
  els.drawerMask.classList.remove("hidden");
}

function closeDrawer() {
  els.drawer.classList.add("hidden");
  els.drawerMask.classList.add("hidden");
  els.drawerForm.innerHTML = "";
  els.drawerForm.onsubmit = null;
}

function fieldHtml(field, row = null) {
  const value = row ? normalizeFieldValue(row[field.key]) : field.defaultValue || "";
  const hint = field.hint ? `<small class="field-hint">${field.hint}</small>` : "";
  if (field.type === "select") {
    return `
      <label class="field">
        <span class="field-label">${field.label}</span>
        <select name="${field.key}" ${field.required ? "required" : ""}>
          ${(field.options || []).map((option) => `<option value="${option.value}" ${String(option.value) === String(value) ? "selected" : ""}>${option.label}</option>`).join("")}
        </select>
        ${hint}
      </label>
    `;
  }
  if (field.type === "textarea") {
    return `
      <label class="field">
        <span class="field-label">${field.label}</span>
        <textarea name="${field.key}" ${field.required ? "required" : ""}>${escapeHtml(value)}</textarea>
        ${hint}
      </label>
    `;
  }
  return `
    <label class="field">
      <span class="field-label">${field.label}</span>
      <input name="${field.key}" type="${field.type || "text"}" step="${field.step || ""}" value="${escapeAttr(value)}" ${field.required ? "required" : ""}>
      ${hint}
    </label>
  `;
}

function dispatchFields(row = {}) {
  return `
    <div class="form-row two">
      ${fieldHtml({ key: "passenger_name", label: "乘客姓名", required: true }, row)}
      ${fieldHtml({ key: "passenger_phone", label: "乘客手机", required: true }, row)}
    </div>
    ${fieldHtml({ key: "start_address", label: "出发地址", required: true }, row)}
    ${fieldHtml({ key: "end_address", label: "目的地", required: true }, row)}
    <div class="form-row two">
      ${fieldHtml({ key: "car_type", label: "车型", type: "select", options: [{ value: 1, label: "经济型" }, { value: 2, label: "舒适型" }, { value: 3, label: "商务型" }, { value: 4, label: "豪华型" }] }, row)}
      ${fieldHtml({ key: "driver_id", label: "指定司机 ID", type: "number" }, row)}
    </div>
    <div class="form-row two">
      ${fieldHtml({ key: "estimated_price", label: "预估价格", type: "number", step: "0.01" }, row)}
      ${fieldHtml({ key: "remark", label: "备注" }, row)}
    </div>
    <input type="hidden" name="start_lng" value="113.2644">
    <input type="hidden" name="start_lat" value="23.1291">
    <input type="hidden" name="end_lng" value="113.3245">
    <input type="hidden" name="end_lat" value="23.1065">
  `;
}

function drawerActionsHtml(submitText = "保存") {
  return `
    <div class="drawer-actions">
      <button class="ghost-btn" type="button" data-cancel>取消</button>
      <button class="primary-btn" type="submit">${submitText}</button>
    </div>
  `;
}

document.addEventListener("click", (event) => {
  if (event.target.matches("[data-cancel]")) closeDrawer();
});

function collectForm(form, fields) {
  const formData = new FormData(form);
  const payload = {};
  for (const [key, value] of formData.entries()) {
    const field = fields?.find((item) => item.key === key);
    payload[key] = field?.type === "number" ? Number(value || 0) : value;
  }
  return payload;
}

function collectFilters(base) {
  els.content.querySelectorAll("[data-filter]").forEach((filter) => {
    if (filter.value !== "") base[filter.dataset.filter] = Number(filter.value) || filter.value;
  });
  return base;
}

function renderFilter(filter) {
  return `
    <select data-filter="${filter.key}">
      <option value="">${filter.label}</option>
      ${filter.options.map((option) => `<option value="${option.value}">${option.label}</option>`).join("")}
    </select>
  `;
}

function metricCard(label, value, trend, icon) {
  return `
    <div class="metric-card">
      <div class="metric-title">
        <span>${label}</span>
        <i class="metric-icon">${icon}</i>
      </div>
      <strong>${value}</strong>
      <em>${trend}</em>
    </div>
  `;
}

function statusLine(label, value) {
  return `
    <div class="status-line">
      <span class="status-name">${label}</span>
      <span class="bar"><i style="width:${value}%"></i></span>
      <strong>${value}%</strong>
    </div>
  `;
}

function todoItem(label, value) {
  return `<div class="todo-item"><span>${label}</span><strong>${value}</strong></div>`;
}

function cell(row, col) {
  if (col.render) return col.render(row, col);
  return escapeHtml(row[col.key] ?? "");
}

function rawCell(row, col) {
  if (col.key === "route") return `${row.start_address || ""} -> ${row.end_address || ""}`;
  return row[col.key] ?? "";
}

function routeCell(row) {
  return `${escapeHtml(row.start_address || "")} → ${escapeHtml(row.end_address || "")}`;
}

function menuNameCell(row) {
  const prefix = row.level ? "　".repeat(row.level) + "└ " : "";
  return `${prefix}${escapeHtml(row.name || "")}`;
}

function moneyCell(row, col) {
  const value = Number(row[col.key] || 0);
  return `¥${value.toFixed(2)}`;
}

function timeCell(row, col) {
  const value = row[col.key];
  if (!value) return "-";
  if (typeof value === "number") return new Date(value * 1000).toLocaleString("zh-CN", { hour12: false });
  return String(value).replace("T", " ").replace("Z", "").slice(0, 19);
}

function normalStatusTag(row, col) {
  return Number(row[col.key]) === 1 ? tag("正常", "green") : tag("冻结", "red");
}

function enableStatusTag(row, col) {
  return Number(row[col.key]) === 1 ? tag("启用", "green") : tag("禁用", "red");
}

function orderStatusTag(row, col) {
  const map = {
    1: ["待接单", "orange"],
    2: ["已接单", "blue"],
    3: ["行程中", "green"],
    4: ["已完成", "green"],
    5: ["已取消", "red"],
  };
  const item = map[Number(row[col.key])] || ["未知", "blue"];
  return tag(item[0], item[1]);
}

function tag(text, color) {
  return `<span class="tag ${color}">${text}</span>`;
}

function orderStatusOptions() {
  return [
    { value: 1, label: "待接单" },
    { value: 2, label: "已接单" },
    { value: 3, label: "行程中" },
    { value: 4, label: "已完成" },
    { value: 5, label: "已取消" },
  ];
}

function flattenTree(list, level = 0) {
  return (list || []).flatMap((item) => {
    const current = { ...item, level };
    return [current, ...flattenTree(item.children || [], level + 1)];
  });
}

function filterRows(rows, keyword) {
  if (!keyword) return rows;
  const lower = keyword.toLowerCase();
  return rows.filter((row) => JSON.stringify(row).toLowerCase().includes(lower));
}

function normalizeFieldValue(value) {
  if (Array.isArray(value)) return value.join(",");
  return value ?? "";
}

function emptyState(text) {
  return `<div class="empty-state">${text}</div>`;
}

function toast(message) {
  els.toast.textContent = message;
  els.toast.classList.add("show");
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => els.toast.classList.remove("show"), 2200);
}

function escapeHtml(value) {
  return String(value ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function escapeAttr(value) {
  return escapeHtml(value);
}

function wait(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function demoOrders() {
  return [
    { order_id: "FP202605250001", passenger_name: "李晓雨", passenger_phone: "13800001231", driver_name: "陈师傅", driver_phone: "13600008881", car_type: "舒适型", start_address: "广州塔", end_address: "珠江新城", status: 3, estimated_price: 42.5, final_price: 45.8, created_at: "2026-05-25T09:40:00Z" },
    { order_id: "FP202605250002", passenger_name: "王先生", passenger_phone: "13900004567", driver_name: "周师傅", driver_phone: "13600008882", car_type: "经济型", start_address: "天河客运站", end_address: "白云机场 T2", status: 2, estimated_price: 128, final_price: 0, created_at: "2026-05-25T09:30:00Z" },
    { order_id: "FP202605250003", passenger_name: "赵敏", passenger_phone: "13700006789", driver_name: "黄师傅", driver_phone: "13600008883", car_type: "商务型", start_address: "琶洲会展中心", end_address: "广州南站", status: 4, estimated_price: 96, final_price: 101.3, created_at: "2026-05-25T08:38:00Z" },
    { order_id: "FP202605250004", passenger_name: "孙浩", passenger_phone: "13500009876", driver_name: "", driver_phone: "", car_type: "经济型", start_address: "北京路步行街", end_address: "体育西路", status: 1, estimated_price: 35.6, final_price: 0, created_at: "2026-05-25T08:10:00Z" },
    { order_id: "FP202605250005", passenger_name: "林可", passenger_phone: "13200002345", driver_name: "吴师傅", driver_phone: "13600008885", car_type: "豪华型", start_address: "花城汇", end_address: "番禺广场", status: 5, estimated_price: 115, final_price: 0, created_at: "2026-05-25T07:42:00Z" },
    { order_id: "FP202605250006", passenger_name: "何悦", passenger_phone: "13100007654", driver_name: "梁师傅", driver_phone: "13600008886", car_type: "舒适型", start_address: "大学城北", end_address: "科学城", status: 4, estimated_price: 68.2, final_price: 70.4, created_at: "2026-05-25T07:18:00Z" },
  ];
}

function demoAdmins() {
  return [
    { id: 1, username: "admin", role_id: 1, role_name: "超级管理员", status: 1, last_login_time: "2026-05-25T09:12:00Z" },
    { id: 2, username: "operation", role_id: 2, role_name: "运营管理员", status: 1, last_login_time: "2026-05-24T18:40:00Z" },
    { id: 3, username: "service", role_id: 3, role_name: "客服专员", status: 1, last_login_time: "2026-05-24T16:22:00Z" },
    { id: 4, username: "audit", role_id: 4, role_name: "审核专员", status: 0, last_login_time: "2026-05-20T10:08:00Z" },
  ];
}

function demoRoles() {
  return [
    { id: 1, name: "超级管理员", remark: "拥有全部菜单与操作权限", menu_ids: [1, 2, 3, 4, 5, 6, 7] },
    { id: 2, name: "运营管理员", remark: "负责订单、城市、车型配置", menu_ids: [1, 2, 6, 7] },
    { id: 3, name: "客服专员", remark: "负责售后、派单与订单查询", menu_ids: [1, 2, 3] },
  ];
}

function demoMenus() {
  return [
    { id: 1, parent_id: 0, name: "运营中心", path: "/operation", icon: "dashboard", sort: 1, children: [
      { id: 11, parent_id: 1, name: "订单管理", path: "/operation/orders", icon: "orders", sort: 1 },
      { id: 12, parent_id: 1, name: "后台派单", path: "/operation/dispatch", icon: "plus", sort: 2 },
    ] },
    { id: 2, parent_id: 0, name: "系统管理", path: "/system", icon: "system", sort: 2, children: [
      { id: 21, parent_id: 2, name: "管理员", path: "/system/admins", icon: "user", sort: 1 },
      { id: 22, parent_id: 2, name: "角色权限", path: "/system/roles", icon: "role", sort: 2 },
      { id: 23, parent_id: 2, name: "菜单管理", path: "/system/menus", icon: "menu", sort: 3 },
    ] },
    { id: 3, parent_id: 0, name: "基础配置", path: "/config", icon: "setting", sort: 3, children: [
      { id: 31, parent_id: 3, name: "城市配置", path: "/config/cities", icon: "city", sort: 1 },
      { id: 32, parent_id: 3, name: "车型配置", path: "/config/car-types", icon: "car", sort: 2 },
    ] },
  ];
}

function demoCities() {
  return [
    { id: 1, city_code: "020", city_name: "广州", status: 1 },
    { id: 2, city_code: "0755", city_name: "深圳", status: 1 },
    { id: 3, city_code: "010", city_name: "北京", status: 1 },
    { id: 4, city_code: "021", city_name: "上海", status: 0 },
  ];
}

function demoCarTypes() {
  return [
    { id: 1, type_name: "经济型", base_price: 8, km_price: 1.5, minute_price: 0.5, status: 1 },
    { id: 2, type_name: "舒适型", base_price: 12, km_price: 2.2, minute_price: 0.8, status: 1 },
    { id: 3, type_name: "商务型", base_price: 18, km_price: 3.2, minute_price: 1.2, status: 1 },
    { id: 4, type_name: "豪华型", base_price: 28, km_price: 5.6, minute_price: 2.0, status: 0 },
  ];
}
