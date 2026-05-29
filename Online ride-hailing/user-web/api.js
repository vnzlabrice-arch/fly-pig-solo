(function initUserApi(global) {
  const DEFAULT_API_BASE_URL = "http://127.0.0.1:8888";

  function getApiBaseUrl() {
    const saved = global.localStorage.getItem("api_base_url");
    if (saved) return saved;
    return DEFAULT_API_BASE_URL;
  }

  function isSuccessCode(code) {
    return Number(code) === 0 || Number(code) === 200;
  }

  async function request(path, payload, options = {}) {
    const baseUrl = getApiBaseUrl();
    const token = global.localStorage.getItem("user_token") || "";
    const method = (options.method || "POST").toUpperCase();
    const isGetOrHead = method === "GET" || method === "HEAD";

    const headers = {};
    if (!isGetOrHead) {
      headers["Content-Type"] = "application/json";
    }
    if (options.withAuth && token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const fetchOptions = {
      method: method,
    };
    if (Object.keys(headers).length > 0) {
      fetchOptions.headers = headers;
    }
    if (payload && !isGetOrHead) {
      fetchOptions.body = JSON.stringify(payload);
    }

    const response = await fetch(`${baseUrl}${path}`, fetchOptions);

    if (!response.ok) {
      throw new Error(`请求失败: HTTP ${response.status}`);
    }

    return response.json();
  }

       async function safeRequest(path, payload, options) {
    try {
      const result = await request(path, payload, options);
      const code = Number(result.code);
      const message = result.message || "请求成功";

      let data = result.data;
      if (data === undefined || data === null) {
        data = {};
        for (const key in result) {
          if (key !== 'code' && key !== 'message') {
            data[key] = result[key];
          }
        }
      }

      return {
        code,
        message,
        data,
      };
    } catch (error) {
      return {
        code: 500,
        message: error instanceof Error ? error.message : "网络请求失败",
        data: null,
      };
    }
  }

  global.UserApi = {
    isSuccessCode,
    getApiBaseUrl,
    setApiBaseUrl(url) {
      global.localStorage.setItem("api_base_url", String(url || "").trim());
    },

    async sendCode(req) {
      const resp = await safeRequest("/user/sendCode", { phone: req.phone });
      return resp;
    },

    async register(req) {
      const resp = await safeRequest("/user/register", req);
      return resp;
    },

    async login(req) {
      const resp = await safeRequest("/user/login", req);
      if (isSuccessCode(resp.code)) {
        let userId = resp.data?.UserId || resp.data?.user_id || 0;
        const token = resp.data?.Token || resp.data?.token || "";

        if (typeof userId === 'string') {
          if (userId.includes(':')) {
            userId = userId.split(':')[0];
          }
          userId = parseInt(userId, 10) || 0;
        }

        if (userId > 0 && token) {
          localStorage.setItem("user_token", String(token));
          localStorage.setItem("user_id", String(userId));
        }
      }
      return resp;
    },

    async getUserDetail(req) {
      // 严格解析 user_id
      let userId = req?.user_id || req?.UserId || 0;
      
      if (typeof userId === 'string') {
        userId = userId.trim();
        // 处理异常格式如 "2:1" 或 ":1"
        if (userId.includes(':')) {
          userId = userId.split(':')[0];
        }
        // 只保留数字
        userId = userId.replace(/[^0-9]/g, '');
        userId = parseInt(userId, 10) || 0;
      }
      
      // 如果参数中的 userId 无效，尝试从 localStorage 获取
      if (!userId || userId <= 0) {
        let saved = localStorage.getItem("user_id");
        if (saved) {
          saved = String(saved).trim();
          if (saved.includes(':')) {
            saved = saved.split(':')[0];
          }
          saved = saved.replace(/[^0-9]/g, '');
          userId = parseInt(saved, 10) || 0;
        }
      }
      
      if (!userId || userId <= 0) {
        console.error('getUserDetail: user_id 无效');
        return { code: 400, message: "用户ID无效，请重新登录", data: null };
      }
      
      const resp = await safeRequest(
        "/user/detail",
        { user_id: userId },
        { method: "POST", withAuth: true }
      );
      return resp;
    },

    async updatePhone(req) {
      return safeRequest("/user/updatePhone", req, { withAuth: true });
    },

    async realNameAuth(req) {
      return safeRequest("/user/realNameAuth", req, { withAuth: true });
    },

    async getUserCouponList(req) {
      const params = new URLSearchParams();
      // go-zero form 标签默认必填，缺失字段会导致 400，必须全部带上
      params.append("userId", req.userId ?? 0);
      params.append("status", req.status ?? 0);
      params.append("page", req.page ?? 1);
      params.append("pageSize", req.pageSize ?? 20);
      
      const url = `/user/getUserCoupon?${params.toString()}`;
      return safeRequest(url, null, { method: "GET" });
    },

    async receiveCoupon(req) {
      return safeRequest("/user/addCoupon", req, { withAuth: true });
    },

    async addPassengerOrder(req) {
      return safeRequest("/user/addPassengerOrderReq", req, { withAuth: true });
    },

    async getUserOrderList(req) {
      const params = new URLSearchParams();
      // go-zero form 标签默认必填，缺失字段会导致 400，必须全部带上
      params.append("userId", req.userId ?? 0);
      params.append("status", req.status ?? 0);
      params.append("page", req.page ?? 1);
      params.append("pageSize", req.pageSize ?? 20);
      
      const url = `/user/getUserOrderList?${params.toString()}`;
      return safeRequest(url, null, { method: "GET" });
    },

    async getPassengerOrderDetail(req) {
      const params = new URLSearchParams();
      // go-zero form 标签默认必填，缺失字段会导致 400，必须全部带上
      params.append("userId", req.userId ?? 0);
      params.append("orderId", req.orderId ?? "");
      
      const url = `/user/getPassengerOrderDetail?${params.toString()}`;
      return safeRequest(url, null, { method: "GET" });
    },
  };
})(window);
