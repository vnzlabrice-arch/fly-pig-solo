# admin-srv 后台管理系统 接口文档

---

## 一、系统管理模块 (system)

### 1.1 管理员认证

#### 1. 管理员登录
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/login |
| **请求方式** | POST |
| **模块** | 管理员认证 |

**请求参数：**
```json
{
  "username": "string",    // 用户名，必填
  "password": "string"     // 密码，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "string",           // JWT令牌
    "admin_id": 1,                // 管理员ID
    "username": "string",        // 用户名
    "role_id": 1,                // 角色ID
    "role_name": "string",       // 角色名称
    "expires_at": "2024-01-01T00:00:00Z"  // 过期时间
  }
}
```

---

#### 2. 管理员登出
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/logout |
| **请求方式** | POST |
| **模块** | 管理员认证 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 3. 获取当前管理员信息
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/info |
| **请求方式** | GET |
| **模块** | 管理员认证 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "admin_id": 1,
    "username": "string",
    "nickname": "string",
    "avatar": "string",
    "role_id": 1,
    "role_name": "string",
    "menus": ["string"]
  }
}
```

---

### 1.2 管理员管理

#### 4. 管理员列表
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/list |
| **请求方式** | GET |
| **模块** | 管理员管理 |

**请求参数：**
```json
{
  "page": 1,           // 页码，默认1
  "page_size": 10,     // 每页数量，默认10
  "keyword": "string"  // 搜索关键词（用户名/手机号）
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "string",
        "nickname": "string",
        "role_id": 1,
        "role_name": "string",
        "status": 1,
        "last_login_time": "2024-01-01T00:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

#### 5. 管理员详情
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/detail/:id |
| **请求方式** | GET |
| **模块** | 管理员管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "username": "string",
    "nickname": "string",
    "avatar": "string",
    "role_id": 1,
    "role_name": "string",
    "status": 1,
    "last_login_time": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 6. 创建管理员
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/create |
| **请求方式** | POST |
| **模块** | 管理员管理 |

**请求参数：**
```json
{
  "username": "string",    // 用户名，必填
  "password": "string",    // 密码，必填
  "nickname": "string",    // 昵称
  "role_id": 1,            // 角色ID，必填
  "status": 1              // 状态：1-正常，2-冻结
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 7. 编辑管理员
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/update/:id |
| **请求方式** | PUT |
| **模块** | 管理员管理 |

**请求参数：**
```json
{
  "nickname": "string",    // 昵称
  "role_id": 1,            // 角色ID
  "status": 1              // 状态：1-正常，2-冻结
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 8. 删除管理员
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/delete/:id |
| **请求方式** | DELETE |
| **模块** | 管理员管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 9. 修改密码
| 项目 | 内容 |
|------|------|
| **接口名** | /admin/password |
| **请求方式** | PUT |
| **模块** | 管理员管理 |

**请求参数：**
```json
{
  "old_password": "string",  // 旧密码，必填
  "new_password": "string"    // 新密码，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.3 角色管理

#### 10. 角色列表
| 项目 | 内容 |
|------|------|
| **接口名** | /role/list |
| **请求方式** | GET |
| **模块** | 角色管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "string",
        "remark": "string",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 11. 角色详情
| 项目 | 内容 |
|------|------|
| **接口名** | /role/detail/:id |
| **请求方式** | GET |
| **模块** | 角色管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "string",
    "remark": "string",
    "menu_ids": [1, 2, 3]
  }
}
```

---

#### 12. 创建角色
| 项目 | 内容 |
|------|------|
| **接口名** | /role/create |
| **请求方式** | POST |
| **模块** | 角色管理 |

**请求参数：**
```json
{
  "name": "string",       // 角色名称，必填
  "remark": "string",     // 备注
  "menu_ids": [1, 2, 3]   // 菜单ID数组
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 13. 编辑角色
| 项目 | 内容 |
|------|------|
| **接口名** | /role/update/:id |
| **请求方式** | PUT |
| **模块** | 角色管理 |

**请求参数：**
```json
{
  "name": "string",
  "remark": "string",
  "menu_ids": [1, 2, 3]
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 14. 删除角色
| 项目 | 内容 |
|------|------|
| **接口名** | /role/delete/:id |
| **请求方式** | DELETE |
| **模块** | 角色管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.4 菜单管理

#### 15. 菜单列表（树形）
| 项目 | 内容 |
|------|------|
| **接口名** | /menu/list |
| **请求方式** | GET |
| **模块** | 菜单管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "name": "string",
      "path": "string",
      "icon": "string",
      "sort": 1,
      "children": [
        {
          "id": 2,
          "parent_id": 1,
          "name": "string",
          "path": "string",
          "icon": "string",
          "sort": 1
        }
      ]
    }
  ]
}
```

---

#### 16. 创建菜单
| 项目 | 内容 |
|------|------|
| **接口名** | /menu/create |
| **请求方式** | POST |
| **模块** | 菜单管理 |

**请求参数：**
```json
{
  "parent_id": 0,      // 父菜单ID，0为顶级
  "name": "string",    // 菜单名称，必填
  "path": "string",    // 菜单路径
  "icon": "string",    // 菜单图标
  "sort": 1            // 排序
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 17. 编辑菜单
| 项目 | 内容 |
|------|------|
| **接口名** | /menu/update/:id |
| **请求方式** | PUT |
| **模块** | 菜单管理 |

**请求参数：**
```json
{
  "parent_id": 0,
  "name": "string",
  "path": "string",
  "icon": "string",
  "sort": 1
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 18. 删除菜单
| 项目 | 内容 |
|------|------|
| **接口名** | /menu/delete/:id |
| **请求方式** | DELETE |
| **模块** | 菜单管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.5 城市配置

#### 19. 城市列表
| 项目 | 内容 |
|------|------|
| **接口名** | /city/list |
| **请求方式** | GET |
| **模块** | 城市配置 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "keyword": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "city_code": "010",
        "city_name": "北京",
        "status": 1
      }
    ],
    "total": 100
  }
}
```

---

#### 20. 创建城市
| 项目 | 内容 |
|------|------|
| **接口名** | /city/create |
| **请求方式** | POST |
| **模块** | 城市配置 |

**请求参数：**
```json
{
  "city_code": "string",   // 城市编码，必填
  "city_name": "string",   // 城市名称，必填
  "status": 1              // 状态：1-启用，2-禁用
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 21. 编辑城市
| 项目 | 内容 |
|------|------|
| **接口名** | /city/update/:id |
| **请求方式** | PUT |
| **模块** | 城市配置 |

**请求参数：**
```json
{
  "city_code": "string",
  "city_name": "string",
  "status": 1
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 22. 删除城市
| 项目 | 内容 |
|------|------|
| **接口名** | /city/delete/:id |
| **请求方式** | DELETE |
| **模块** | 城市配置 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.6 车型配置

#### 23. 车型列表
| 项目 | 内容 |
|------|------|
| **接口名** | /car-type/list |
| **请求方式** | GET |
| **模块** | 车型配置 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "type_name": "经济型",
        "base_price": 8.00,
        "km_price": 1.50,
        "minute_price": 0.50,
        "status": 1
      }
    ],
    "total": 10
  }
}
```

---

#### 24. 创建车型
| 项目 | 内容 |
|------|------|
| **接口名** | /car-type/create |
| **请求方式** | POST |
| **模块** | 车型配置 |

**请求参数：**
```json
{
  "type_name": "string",     // 车型名称，必填
  "base_price": 8.00,       // 起步价，必填
  "km_price": 1.50,         // 公里单价，必填
  "minute_price": 0.50,     // 时长单价，必填
  "status": 1               // 状态
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 25. 编辑车型
| 项目 | 内容 |
|------|------|
| **接口名** | /car-type/update/:id |
| **请求方式** | PUT |
| **模块** | 车型配置 |

**请求参数：**
```json
{
  "type_name": "string",
  "base_price": 8.00,
  "km_price": 1.50,
  "minute_price": 0.50,
  "status": 1
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 26. 删除车型
| 项目 | 内容 |
|------|------|
| **接口名** | /car-type/delete/:id |
| **请求方式** | DELETE |
| **模块** | 车型配置 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.7 系统配置

#### 27. 系统配置详情
| 项目 | 内容 |
|------|------|
| **接口名** | /system-config/detail |
| **请求方式** | GET |
| **模块** | 系统配置 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "configs": [
      {
        "id": 1,
        "config_key": "platform_name",
        "config_value": "花小猪打车",
        "remark": "平台名称"
      }
    ]
  }
}
```

---

#### 28. 编辑系统配置
| 项目 | 内容 |
|------|------|
| **接口名** | /system-config/update |
| **请求方式** | PUT |
| **模块** | 系统配置 |

**请求参数：**
```json
{
  "configs": [
    {
      "config_key": "platform_name",
      "config_value": "花小猪打车"
    }
  ]
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.8 风控配置

#### 29. 风控配置详情
| 项目 | 内容 |
|------|------|
| **接口名** | /risk-config/detail |
| **请求方式** | GET |
| **模块** | 风控配置 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "need_realname": 1,
    "blacklist_days": 7
  }
}
```

---

#### 30. 编辑风控配置
| 项目 | 内容 |
|------|------|
| **接口名** | /risk-config/update |
| **请求方式** | PUT |
| **模块** | 风控配置 |

**请求参数：**
```json
{
  "need_realname": 1,    // 是否需要实名：1-是，0-否
  "blacklist_days": 7     // 黑名单天数
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.9 文章管理

#### 31. 文章分类列表
| 项目 | 内容 |
|------|------|
| **接口名** | /article-category/list |
| **请求方式** | GET |
| **模块** | 文章管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "string",
      "parent_id": 0,
      "status": 1
    }
  ]
}
```

---

#### 32. 创建文章分类
| 项目 | 内容 |
|------|------|
| **接口名** | /article-category/create |
| **请求方式** | POST |
| **模块** | 文章管理 |

**请求参数：**
```json
{
  "name": "string",     // 分类名称，必填
  "parent_id": 0,       // 父分类ID
  "status": 1           // 状态
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 33. 编辑文章分类
| 项目 | 内容 |
|------|------|
| **接口名** | /article-category/update/:id |
| **请求方式** | PUT |
| **模块** | 文章管理 |

**请求参数：**
```json
{
  "name": "string",
  "parent_id": 0,
  "status": 1
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 34. 删除文章分类
| 项目 | 内容 |
|------|------|
| **接口名** | /article-category/delete/:id |
| **请求方式** | DELETE |
| **模块** | 文章管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 35. 文章列表
| 项目 | 内容 |
|------|------|
| **接口名** | /article/list |
| **请求方式** | GET |
| **模块** | 文章管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "category_id": 1,
  "keyword": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "string",
        "category_id": 1,
        "category_name": "string",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 36. 文章详情
| 项目 | 内容 |
|------|------|
| **接口名** | /article/detail/:id |
| **请求方式** | GET |
| **模块** | 文章管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "title": "string",
    "category_id": 1,
    "content": "string",
    "status": 1,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

---

#### 37. 创建文章
| 项目 | 内容 |
|------|------|
| **接口名** | /article/create |
| **请求方式** | POST |
| **模块** | 文章管理 |

**请求参数：**
```json
{
  "title": "string",       // 文章标题，必填
  "category_id": 1,         // 分类ID，必填
  "content": "string",     // 文章内容，必填
  "status": 1              // 状态：1-发布，2-草稿
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

#### 38. 编辑文章
| 项目 | 内容 |
|------|------|
| **接口名** | /article/update/:id |
| **请求方式** | PUT |
| **模块** | 文章管理 |

**请求参数：**
```json
{
  "title": "string",
  "category_id": 1,
  "content": "string",
  "status": 1
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

#### 39. 删除文章
| 项目 | 内容 |
|------|------|
| **接口名** | /article/delete/:id |
| **请求方式** | DELETE |
| **模块** | 文章管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

### 1.10 APP升级配置

#### 40. APP升级配置详情
| 项目 | 内容 |
|------|------|
| **接口名** | /app-update/detail |
| **请求方式** | GET |
| **模块** | APP升级配置 |

**请求参数：**
```json
{
  "platform": 1  // 平台：1-Android，2-iOS
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "platform": 1,
    "version": "1.0.0",
    "download_url": "string",
    "force_update": 0,
    "content": "string"
  }
}
```

---

#### 41. 编辑APP升级配置
| 项目 | 内容 |
|------|------|
| **接口名** | /app-update/update |
| **请求方式** | PUT |
| **模块** | APP升级配置 |

**请求参数：**
```json
{
  "platform": 1,            // 平台，必填
  "version": "string",      // 版本号，必填
  "download_url": "string", // 下载地址，必填
  "force_update": 0,        // 是否强制更新：0-否，1-是
  "content": "string"      // 更新内容
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

## 二、售后模块 (aftersale)

### 2.1 售后工单管理

#### 42. 售后工单列表
| 项目 | 内容 |
|------|------|
| **接口名** | /aftersale/list |
| **请求方式** | GET |
| **模块** | 售后工单管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "type": 1,             // 类型：1-退款，2-投诉，3-其他
  "status": 1,            // 状态：1-待处理，2-处理中，3-已完成，4-已拒绝
  "keyword": "string",    // 关键词（工单号/订单号）
  "start_time": "string", // 开始时间
  "end_time": "string"    // 结束时间
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "aftersale_no": "string",
        "order_id": "string",
        "type": 1,
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 43. 售后工单详情
| 项目 | 内容 |
|------|------|
| **接口名** | /aftersale/detail/:id |
| **请求方式** | GET |
| **模块** | 售后工单管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "aftersale_no": "string",
    "order_id": "string",
    "type": 1,
    "status": 1,
    "reason": "string",
    "images": ["string"],
    "created_at": "2024-01-01T00:00:00Z",
    "audit_logs": [
      {
        "admin_id": 1,
        "admin_name": "string",
        "status": 1,
        "note": "string",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

#### 44. 审核售后工单
| 项目 | 内容 |
|------|------|
| **接口名** | /aftersale/audit |
| **请求方式** | POST |
| **模块** | 售后工单管理 |

**请求参数：**
```json
{
  "aftersale_id": 1,       // 售后工单ID，必填
  "status": 2,              // 审核状态：2-处理中，3-已完成，4-已拒绝，必填
  "note": "string"          // 审核备注
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

---

## 三、订单模块 (order)

### 3.1 发票管理

#### 45. 发票列表
| 项目 | 内容 |
|------|------|
| **接口名** | /invoice/list |
| **请求方式** | GET |
| **模块** | 发票管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "status": 1,              // 状态：1-待开票，2-已开票，3-已作废
  "keyword": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "invoice_no": "string",
        "order_id": "string",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 46. 开具发票
| 项目 | 内容 |
|------|------|
| **接口名** | /invoice/create |
| **请求方式** | POST |
| **模块** | 发票管理 |

**请求参数：**
```json
{
  "order_id": "string"      // 订单ID，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

### 3.2 结算管理

#### 47. 结算记录列表
| 项目 | 内容 |
|------|------|
| **接口名** | /settle/list |
| **请求方式** | GET |
| **模块** | 结算管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "status": 1,              // 状态：1-待结算，2-已结算，3-结算失败
  "start_time": "string",
  "end_time": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "order_id": "string",
        "real_amount": 100.00,
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

### 3.3 坏账管理

#### 48. 坏账订单列表
| 项目 | 内容 |
|------|------|
| **接口名** | /bad-debt/list |
| **请求方式** | GET |
| **模块** | 坏账管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "order_id": "string",
        "debt_amount": 100.00,
        "overdue_days": 7,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

### 3.4 付款单管理

#### 49. 付款单列表
| 项目 | 内容 |
|------|------|
| **接口名** | /payment/list |
| **请求方式** | GET |
| **模块** | 付款单管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "status": 1,              // 状态：1-待付款，2-已付款，3-付款失败
  "keyword": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "payment_no": "string",
        "amount": 100.00,
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 50. 创建付款单（退款/赔付）
| 项目 | 内容 |
|------|------|
| **接口名** | /payment/create |
| **请求方式** | POST |
| **模块** | 付款单管理 |

**请求参数：**
```json
{
  "order_id": "string",      // 订单ID，必填
  "amount": 100.00,          // 金额，必填
  "type": 1,                 // 类型：1-退款，2-赔付
  "note": "string"           // 备注
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1
  }
}
```

---

## 四、日志模块 (logs)

### 4.1 登录日志

#### 51. 登录日志列表
| 项目 | 内容 |
|------|------|
| **接口名** | /log/login/list |
| **请求方式** | GET |
| **模块** | 日志管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "admin_id": 1,
  "start_time": "string",
  "end_time": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "admin_id": 1,
        "username": "string",
        "ip": "string",
        "login_time": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

### 4.2 操作日志

#### 52. 操作日志列表
| 项目 | 内容 |
|------|------|
| **接口名** | /log/operation/list |
| **请求方式** | GET |
| **模块** | 日志管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "admin_id": 1,
  "module": "string",
  "start_time": "string",
  "end_time": "string"
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "admin_id": 1,
        "username": "string",
        "module": "string",
        "operation": "string",
        "detail": "string",
        "ip": "string",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

## 接口数量统计

| 模块 | 接口数量 |
|------|----------|
| 管理员认证 | 3 |
| 管理员管理 | 5 |
| 角色管理 | 5 |
| 菜单管理 | 4 |
| 城市配置 | 4 |
| 车型配置 | 4 |
| 系统配置 | 2 |
| 风控配置 | 2 |
| 文章管理 | 9 |
| APP升级配置 | 2 |
| 售后工单管理 | 3 |
| 发票管理 | 2 |
| 结算管理 | 1 |
| 坏账管理 | 1 |
| 付款单管理 | 2 |
| 日志管理 | 2 |
| **总计** | **50** |
