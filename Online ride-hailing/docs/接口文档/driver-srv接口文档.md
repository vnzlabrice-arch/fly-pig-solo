# driver-srv 司机端系统 接口文档

---

## 一、司机认证模块 (auth)

### 1.1 登录注册

#### 1. 发送验证码
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/auth/send-code |
| **请求方式** | POST |
| **模块** | 司机认证 |

**请求参数：**
```json
{
  "phone": "string"     // 手机号，必填，11位
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "expires_at": 300  // 验证码有效期（秒）
  }
}
```

---

#### 2. 验证码登录/注册
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/auth/login |
| **请求方式** | POST |
| **模块** | 司机认证 |

**请求参数：**
```json
{
  "phone": "string",     // 手机号，必填
  "code": "string"       // 验证码，必填，4-6位
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "token": "string",
    "driver_id": 1,
    "phone": "string",
    "is_new": true,           // 是否新用户
    "is_certified": false     // 是否已认证
  }
}
```

---

#### 3. 退出登录
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/auth/logout |
| **请求方式** | POST |
| **模块** | 司机认证 |

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

## 二、司机信息模块 (driver)

### 2.1 个人信息

#### 4. 获取司机信息
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/info |
| **请求方式** | GET |
| **模块** | 个人信息 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "driver_id": 1,
    "phone": "string",
    "nickname": "string",
    "real_name": "string",
    "avatar_url": "string",
    "gender": 1,
    "status": 1,
    "service_score": 4.80,
    "order_count": 100,
    "is_certified": true,
    "is_bound_vehicle": true
  }
}
```

---

#### 5. 更新个人信息
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/update |
| **请求方式** | PUT |
| **模块** | 个人信息 |

**请求参数：**
```json
{
  "nickname": "string",    // 昵称
  "avatar": "string"       // 头像URL
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

### 2.2 司机认证

#### 6. 提交认证资料
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/cert/submit |
| **请求方式** | POST |
| **模块** | 司机认证 |

**请求参数：**
```json
{
  "real_name": "string",      // 真实姓名，必填
  "license_no": "string",     // 驾驶证号，必填
  "license_expire": "string", // 驾照到期日期，必填，格式：2025-01-01
  "drive_years": 3,           // 驾龄，必填
  "id_card_front": "string",  // 身份证正面照，必填，URL
  "id_card_back": "string",   // 身份证反面照，必填，URL
  "license_pic": "string",    // 驾驶证照片，必填，URL
  "drive_pic": "string"       // 行驶证照片，必填，URL
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "audit_id": 1
  }
}
```

---

#### 7. 获取认证状态
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/cert/status |
| **请求方式** | GET |
| **模块** | 司机认证 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "is_certified": true,
    "audit_status": 1,           // 审核状态：0-待审核，1-通过，2-拒绝
    "real_name": "string",
    "license_no": "string",
    "reject_reason": "string"    // 拒绝原因（审核失败时返回）
  }
}
```

---

### 2.3 车辆管理

#### 8. 绑定车辆
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/vehicle/bind |
| **请求方式** | POST |
| **模块** | 车辆管理 |

**请求参数：**
```json
{
  "plate_no": "string",    // 车牌号，必填
  "brand": "string",       // 品牌，必填
  "model": "string",       // 车型，必填
  "color": "string",       // 颜色，必填
  "license_pic": "string", // 行驶证照片，必填
  "car_pic": "string"      // 车辆照片，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "vehicle_id": 1
  }
}
```

---

#### 9. 获取车辆信息
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/vehicle/info |
| **请求方式** | GET |
| **模块** | 车辆管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "plate_no": "京A12345",
    "brand": "string",
    "model": "string",
    "color": "string",
    "status": 1
  }
}
```

---

#### 10. 更新车辆信息
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/vehicle/update |
| **请求方式** | PUT |
| **模块** | 车辆管理 |

**请求参数：**
```json
{
  "plate_no": "string",
  "brand": "string",
  "model": "string",
  "color": "string"
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

### 3.1 订单列表

#### 11. 获取进行中订单
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/current |
| **请求方式** | GET |
| **模块** | 订单管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "order_id": "string",
    "order_type": 1,
    "car_type": "string",
    "status": 3,
    "passenger_name": "string",
    "passenger_phone": "string",
    "passenger_avatar": "string",
    "start_address": "string",
    "start_lng": 116.404,
    "start_lat": 39.915,
    "end_address": "string",
    "end_lng": 116.504,
    "end_lat": 39.925,
    "pass_remark": "string",
    "book_time": "2024-01-01T10:00:00Z",
    "pickup_time": "2024-01-01T10:15:00Z",
    "start_time": "2024-01-01T10:20:00Z"
  }
}
```

---

#### 12. 获取历史订单列表
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/history |
| **请求方式** | GET |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "status": 5    // 订单状态：4-行程中，5-已完成，6-已取消
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
        "order_id": "string",
        "order_type": 1,
        "car_type": "string",
        "status": 5,
        "start_address": "string",
        "end_address": "string",
        "final_price": 35.50,
        "actual_income": 28.40,
        "book_time": "2024-01-01T10:00:00Z",
        "end_time": "2024-01-01T10:45:00Z"
      }
    ],
    "total": 100
  }
}
```

---

#### 13. 订单详情
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/detail/:order_id |
| **请求方式** | GET |
| **模块** | 订单管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "order_id": "string",
    "order_type": 1,
    "car_type": "string",
    "status": 5,
    "start_address": "string",
    "start_lng": 116.404,
    "start_lat": 39.915,
    "end_address": "string",
    "end_lng": 116.504,
    "end_lat": 39.925,
    "passenger_name": "string",
    "passenger_phone": "string",
    "passenger_avatar": "string",
    "pass_remark": "string",
    "estimated_price": 40.00,
    "final_price": 35.50,
    "platform_fee": 7.10,
    "actual_income": 28.40,
    "book_time": "2024-01-01T10:00:00Z",
    "pickup_time": "2024-01-01T10:15:00Z",
    "start_time": "2024-01-01T10:20:00Z",
    "end_time": "2024-01-01T10:45:00Z",
    "cancel_reason": "string"
  }
}
```

---

### 3.2 接单与行程

#### 14. 抢单/接单
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/accept |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_id": "string"    // 订单ID，必填
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

#### 15. 到达乘客位置
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/arrive |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_id": "string"    // 订单ID，必填
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

#### 16. 开始行程
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/start |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_id": "string"    // 订单ID，必填
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

#### 17. 结束行程
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/end |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_id": "string",   // 订单ID，必填
  "end_lng": 116.504,    // 结束经度，必填
  "end_lat": 39.925      // 结束纬度，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "final_price": 35.50
  }
}
```

---

#### 18. 取消订单（司机取消）
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/order/cancel |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_id": "string",      // 订单ID，必填
  "reason": "string"         // 取消原因，必填
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

## 四、司机状态模块 (status)

### 4.1 上下线管理

#### 19. 司机上线
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/status/online |
| **请求方式** | POST |
| **模块** | 状态管理 |

**请求参数：**
```json
{
  "lng": 116.404,    // 当前经度，必填
  "lat": 39.915      // 当前纬度，必填
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

#### 20. 司机下线
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/status/offline |
| **请求方式** | POST |
| **模块** | 状态管理 |

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

#### 21. 开始接单
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/status/accept-order |
| **请求方式** | POST |
| **模块** | 状态管理 |

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

#### 22. 停止接单
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/status/stop-order |
| **请求方式** | POST |
| **模块** | 状态管理 |

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

#### 23. 更新位置
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/status/location |
| **请求方式** | PUT |
| **模块** | 状态管理 |

**请求参数：**
```json
{
  "lng": 116.404,    // 当前经度，必填
  "lat": 39.915      // 当前纬度，必填
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

### 4.2 服务区域

#### 24. 获取服务区域
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/service-area/list |
| **请求方式** | GET |
| **模块** | 服务区域 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "city_code": "010",
      "city_name": "北京",
      "status": 1
    }
  ]
}
```

---

#### 25. 设置服务区域
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/service-area/set |
| **请求方式** | POST |
| **模块** | 服务区域 |

**请求参数：**
```json
{
  "city_codes": ["010", "021"]   // 城市编码数组，必填
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

## 五、财务模块 (finance)

### 5.1 钱包

#### 26. 获取钱包信息
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/wallet/info |
| **请求方式** | GET |
| **模块** | 钱包 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "balance": 1000.00,       // 账户余额
    "withdrawable": 800.00,  // 可提现金额
    "frozen": 200.00         // 冻结金额
  }
}
```

---

### 5.2 收入明细

#### 27. 获取收入明细列表
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/income/list |
| **请求方式** | GET |
| **模块** | 收入明细 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "start_time": "2024-01-01",
  "end_time": "2024-01-31"
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
        "total_fee": 35.50,
        "platform_fee": 7.10,
        "actual_income": 28.40,
        "created_at": "2024-01-01T10:45:00Z"
      }
    ],
    "total": 100,
    "total_income": 2840.00
  }
}
```

---

#### 28. 收入统计
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/income/stat |
| **请求方式** | GET |
| **模块** | 收入明细 |

**请求参数：**
```json
{
  "type": "today"    // 统计类型：today-今日，week-本周，month-本月
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "order_count": 10,
    "total_fee": 350.00,
    "platform_fee": 70.00,
    "actual_income": 280.00
  }
}
```

---

### 5.3 提现

#### 29. 申请提现
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/withdraw/apply |
| **请求方式** | POST |
| **模块** | 提现 |

**请求参数：**
```json
{
  "amount": 500.00,        // 提现金额，必填
  "bank_card": "string"    // 银行卡号，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "withdraw_id": 1
  }
}
```

---

#### 30. 提现记录列表
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/withdraw/list |
| **请求方式** | GET |
| **模块** | 提现 |

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
        "amount": 500.00,
        "bank_card": "string",
        "status": 2,
        "created_at": "2024-01-01T10:00:00Z",
        "arrive_time": "2024-01-01T18:00:00Z"
      }
    ],
    "total": 10
  }
}
```

---

## 六、消息模块 (message)

### 6.1 系统消息

#### 31. 消息列表
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/message/list |
| **请求方式** | GET |
| **模块** | 消息 |

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
        "title": "string",
        "content": "string",
        "is_read": 0,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "unread_count": 5
  }
}
```

---

#### 32. 标记已读
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/message/read |
| **请求方式** | PUT |
| **模块** | 消息 |

**请求参数：**
```json
{
  "message_id": 1    // 消息ID，必填
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

#### 33. 标记全部已读
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/message/read-all |
| **请求方式** | PUT |
| **模块** | 消息 |

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

## 七、违规模块 (violation)

### 7.1 违规记录

#### 34. 违规记录列表
| 项目 | 内容 |
|------|------|
| **接口名** | /driver/violation/list |
| **请求方式** | GET |
| **模块** | 违规管理 |

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
        "violation_type": 1,
        "violation_name": "string",
        "score_deduct": 5,
        "reason": "string",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 10,
    "total_deduct": 20
  }
}
```

---

## 八、公共模块 (common)

### 8.1 城市与车型

#### 35. 获取支持城市列表
| 项目 | 内容 |
|------|------|
| **接口名** | /common/cities |
| **请求方式** | GET |
| **模块** | 公共接口 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "city_code": "010",
      "city_name": "北京"
    }
  ]
}
```

---

#### 36. 获取车型列表
| 项目 | 内容 |
|------|------|
| **接口名** | /common/car-types |
| **请求方式** | GET |
| **模块** | 公共接口 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "type_name": "经济型",
      "base_price": 8.00,
      "km_price": 1.50,
      "minute_price": 0.50
    }
  ]
}
```

---

### 8.2 基础信息

#### 37. 版本更新检查
| 项目 | 内容 |
|------|------|
| **接口名** | /common/version-check |
| **请求方式** | GET |
| **模块** | 公共接口 |

**请求参数：**
```json
{
  "platform": 1,      // 平台：1-Android，2-iOS
  "version": "1.0.0"  // 当前版本号
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "has_update": true,
    "version": "1.1.0",
    "download_url": "string",
    "force_update": false,
    "content": "string"
  }
}
```

---

## 接口数量统计

| 模块 | 接口数量 |
|------|----------|
| 司机认证 | 3 |
| 个人信息 | 2 |
| 司机认证 | 2 |
| 车辆管理 | 3 |
| 订单列表 | 3 |
| 接单与行程 | 5 |
| 上下线管理 | 5 |
| 服务区域 | 2 |
| 钱包 | 1 |
| 收入明细 | 2 |
| 提现 | 2 |
| 消息 | 3 |
| 违规记录 | 1 |
| 公共接口 | 3 |
| **总计** | **39** |
