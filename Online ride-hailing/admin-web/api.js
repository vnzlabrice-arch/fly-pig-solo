(function initAdminApi(global) {
  const DEFAULT_API_BASE_URL = "http://127.0.0.1:8888";

  function getApiBaseUrl() {
    return global.localStorage.getItem("admin_api_base_url") || DEFAULT_API_BASE_URL;
  }

  function setApiBaseUrl(url) {
    global.localStorage.setItem("admin_api_base_url", String(url || "").trim());
  }

  function isSuccessCode(code) {
    return Number(code) === 0 || Number(code) === 200;
  }

  function getToken() {
    return global.localStorage.getItem("admin_token") || "";
  }

  function saveSession(data) {
    if (!data) return;
    if (data.token) global.localStorage.setItem("admin_token", data.token);
    if (data.admin_id) global.localStorage.setItem("admin_id", String(data.admin_id));
    if (data.username) global.localStorage.setItem("admin_username", data.username);
  }

  async function request(path, payload = {}, options = {}) {
    const controller = new AbortController();
    const timeout = global.setTimeout(() => controller.abort(), options.timeoutMs || 1200);
    const headers = {
      "Content-Type": "application/json",
      ...(getToken() ? { Authorization: `Bearer ${getToken()}` } : {}),
    };

    let response;
    try {
      response = await fetch(`${getApiBaseUrl()}${path}`, {
        method: options.method || "POST",
        headers,
        body: JSON.stringify(payload),
        signal: controller.signal,
      });
    } catch (error) {
      if (error && error.name === "AbortError") {
        throw new Error("接口连接超时");
      }
      throw error;
    } finally {
      global.clearTimeout(timeout);
    }

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const result = await response.json();
    const code = Number(result.code);
    const message = result.message || result.msg || "";

    if (!isSuccessCode(code)) {
      throw new Error(message || "请求失败");
    }

    return {
      code,
      message,
      data: result.data || {},
      raw: result,
    };
  }

  async function safeRequest(path, payload, options) {
    try {
      const result = await request(path, payload, options);
      return { ok: true, ...result };
    } catch (error) {
      return {
        ok: false,
        code: 500,
        message: error instanceof Error ? error.message : "网络请求失败",
        data: null,
      };
    }
  }

  global.AdminApi = {
    getApiBaseUrl,
    setApiBaseUrl,
    isSuccessCode,
    saveSession,
    async login(payload) {
      const resp = await safeRequest("/admin/login", payload);
      if (resp.ok) saveSession(resp.data);
      return resp;
    },
    current(adminId) {
      return safeRequest("/admin/current", { admin_id: Number(adminId) || 1 });
    },
    list(endpoint, payload) {
      return safeRequest(endpoint, payload);
    },
    mutate(endpoint, payload) {
      return safeRequest(endpoint, payload);
    },
  };
})(window);
