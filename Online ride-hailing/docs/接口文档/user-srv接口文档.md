# user-srv 用户端系统 接口文档

---

## 一、用户认证模块 (auth)

### 1.1 登录注册

#### 1. 发送验证码
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/auth/send-code |
| **请求方式** | POST |
| **模块** | 用户认证 |

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
| **接口名** | /passenger/auth/login |
| **请求方式** | POST |
| **模块** | 用户认证 |

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
    "user_id": 1,
    "phone": "string",
    "is_new": true,           // 是否新用户
    "nickname": "string"
  }
}
```

---

#### 3. 退出登录
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/auth/logout |
| **请求方式** | POST |
| **模块** | 用户认证 |

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

## 二、用户信息模块 (user)

### 2.1 个人信息

#### 4. 获取用户信息
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/user/info |
| **请求方式** | GET |
| **模块** | 个人信息 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user_id": 1,
    "phone": "string",
    "nickname": "string",
    "avatar_url": "string",
    "real_name": "string",
    "gender": 1,
    "member_level": 0,
    "balance": 100.00,
    "flower_coin": 50,
    "integral": 1000,
    "status": 1
  }
}
```

---

#### 5. 更新个人信息
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/user/update |
| **请求方式** | PUT |
| **模块** | 个人信息 |

**请求参数：**
```json
{
  "nickname": "string",    // 昵称
  "avatar": "string",       // 头像URL
  "gender": 1,             // 性别：0-未知，1-男，2-女
  "real_name": "string"    // 真实姓名
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

### 2.2 实名认证

#### 6. 提交实名认证
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/user/realname |
| **请求方式** | POST |
| **模块** | 实名认证 |

**请求参数：**
```json
{
  "real_name": "string",          // 真实姓名，必填
  "id_card": "string"             // 身份证号，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "is_certified": true
  }
}
```

---

### 2.3 地址簿

#### 7. 地址簿列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/address/list |
| **请求方式** | GET |
| **模块** | 地址簿管理 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "tag": "家",
      "address": "北京市朝阳区xxx",
      "lng": 116.404,
      "lat": 39.915,
      "is_default": 1
    }
  ]
}
```

---

#### 8. 添加地址
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/address/add |
| **请求方式** | POST |
| **模块** | 地址簿管理 |

**请求参数：**
```json
{
  "tag": "string",         // 地址标签，必填
  "address": "string",     // 详细地址，必填
  "lng": 116.404,          // 经度，必填
  "lat": 39.915,           // 纬度，必填
  "is_default": 0          // 是否默认：0-否，1-是
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

#### 9. 编辑地址
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/address/update/:id |
| **请求方式** | PUT |
| **模块** | 地址簿管理 |

**请求参数：**
```json
{
  "tag": "string",
  "address": "string",
  "lng": 116.404,
  "lat": 39.915,
  "is_default": 0
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

#### 10. 删除地址
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/address/delete/:id |
| **请求方式** | DELETE |
| **模块** | 地址簿管理 |

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

#### 11. 设置默认地址
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/address/default/:id |
| **请求方式** | PUT |
| **模块** | 地址簿管理 |

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

## 三、订单模块 (order)

### 3.1 下单

#### 12. 预估价格
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/estimate |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "car_type": "string",     // 车型，必填
  "start_lng": 116.404,     // 起点经度，必填
  "start_lat": 39.915,      // 起点纬度，必填
  "end_lng": 116.504,       // 终点经度，必填
  "end_lat": 39.925         // 终点纬度，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "estimated_price": 35.50,
    "distance": 5000,        // 距离（米）
    "duration": 1200         // 预估时长（秒）
  }
}
```

---

#### 13. 创建订单
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/create |
| **请求方式** | POST |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "order_type": 1,          // 订单类型：1-即时用车，2-预约用车，3-拼车
  "car_type": "string",     // 车型，必填
  "start_lng": 116.404,     // 起点经度，必填
  "start_lat": 39.915,      // 起点纬度，必填
  "start_address": "string", // 起点地址，必填
  "end_lng": 116.504,       // 终点经度，必填
  "end_lat": 39.925,        // 终点纬度，必填
  "end_address": "string",  // 终点地址，必填
  "pass_remark": "string",   // 备注
  "appoint_time": "string", // 预约时间（仅预约单）
  "coupon_id": 1,           // 优惠券ID（可选）
  "address_id": 1           // 常用地址ID（可选）
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "order_id": "string",
    "estimated_price": 35.50,
    "discount_amount": 5.00
  }
}
```

---

### 3.2 订单状态

#### 14. 获取进行中订单
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/current |
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
    "status": 2,
    "driver_id": 1,
    "driver_name": "string",
    "driver_phone": "string",
    "driver_avatar": "string",
    "driver_score": 4.80,
    "vehicle_plate": "京A12345",
    "start_address": "string",
    "start_lng": 116.404,
    "start_lat": 39.915,
    "end_address": "string",
    "end_lng": 116.504,
    "end_lat": 39.925,
    "pass_remark": "string",
    "estimated_price": 35.50,
    "book_time": "2024-01-01T10:00:00Z",
    "pickup_time": "2024-01-01T10:15:00Z",
    "driver_lng": 116.410,
    "driver_lat": 39.920
  }
}
```

---

#### 15. 取消订单（乘客取消）
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/cancel |
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

#### 16. 订单详情
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/detail/:order_id |
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
    "driver_id": 1,
    "driver_name": "string",
    "driver_phone": "string",
    "vehicle_plate": "京A12345",
    "pass_remark": "string",
    "estimated_price": 40.00,
    "final_price": 35.50,
    "coupon_discount": 5.00,
    "pay_status": 1,
    "pay_time": "2024-01-01T10:50:00Z",
    "book_time": "2024-01-01T10:00:00Z",
    "pickup_time": "2024-01-01T10:15:00Z",
    "start_time": "2024-01-01T10:20:00Z",
    "end_time": "2024-01-01T10:45:00Z",
    "fee_details": [
      {
        "fee_type": "里程费",
        "amount": 20.00,
        "description": "5公里"
      },
      {
        "fee_type": "时长费",
        "amount": 10.00,
        "description": "25分钟"
      },
      {
        "fee_type": "起步价",
        "amount": 8.00,
        "description": "基础费用"
      }
    ]
  }
}
```

---

### 3.3 历史订单

#### 17. 历史订单列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/history |
| **请求方式** | GET |
| **模块** | 订单管理 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "status": 5    // 订单状态：5-已完成，6-已取消
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
        "pay_status": 1,
        "book_time": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100
  }
}
```

---

### 3.4 支付

#### 18. 发起支付
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/order/pay |
| **请求方式** | POST |
| **模块** | 订单支付 |

**请求参数：**
```json
{
  "order_id": "string",       // 订单ID，必填
  "pay_channel": "balance",    // 支付渠道：balance-余额，alipay-支付宝，wechat-微信
  "coupon_id": 1              // 优惠券ID（可选）
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "pay_amount": 35.50,
    "balance": 100.00,
    "pay_result": true
  }
}
```

---

## 四、钱包模块 (wallet)

### 4.1 钱包信息

#### 19. 获取钱包信息
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/wallet/info |
| **请求方式** | GET |
| **模块** | 钱包 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "balance": 100.00,
    "flower_coin": 50,
    "integral": 1000
  }
}
```

---

### 4.2 充值

#### 20. 充值
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/wallet/recharge |
| **请求方式** | POST |
| **模块** | 充值 |

**请求参数：**
```json
{
  "amount": 100.00,           // 充值金额，必填
  "pay_channel": "alipay"     // 支付渠道：alipay-支付宝，wechat-微信
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "trade_no": "string",
    "pay_params": {}
  }
}
```

---

### 4.3 流水记录

#### 21. 流水记录列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/wallet/flow |
| **请求方式** | GET |
| **模块** | 流水记录 |

**请求参数：**
```json
{
  "page": 1,
  "page_size": 10,
  "flow_type": 1    // 流水类型：1-充值，2-消费，3-退款，5-奖励
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
        "flow_type": 2,
        "flow_type_name": "消费",
        "amount": -35.50,
        "balance": 64.50,
        "pay_channel": "余额",
        "order_id": "string",
        "remark": "行程消费",
        "created_at": "2024-01-01T10:50:00Z"
      }
    ],
    "total": 100
  }
}
```

---

## 五、优惠券模块 (coupon)

### 5.1 优惠券领取

#### 22. 领取优惠券
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/coupon/receive |
| **请求方式** | POST |
| **模块** | 优惠券 |

**请求参数：**
```json
{
  "template_id": 1    // 优惠券模板ID，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "user_coupon_id": 1
  }
}
```

---

#### 23. 优惠券列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/coupon/list |
| **请求方式** | GET |
| **模块** | 优惠券 |

**请求参数：**
```json
{
  "status": 1    // 状态：1-未使用，2-已使用，3-已过期
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
        "coupon_no": "string",
        "name": "新人专享券",
        "type": 1,
        "type_name": "满减券",
        "min_amount": 20.00,
        "reduce_amount": 5.00,
        "start_time": "2024-01-01T00:00:00Z",
        "end_time": "2024-01-31T23:59:59Z"
      }
    ],
    "total": 10
  }
}
```

---

#### 24. 可用优惠券（下单时查询）
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/coupon/available |
| **请求方式** | GET |
| **模块** | 优惠券 |

**请求参数：**
```json
{
  "order_amount": 35.50,    // 订单金额，必填
  "car_type": "string"       // 车型（可选）
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "满20减5",
      "type": 1,
      "reduce_amount": 5.00,
      "min_amount": 20.00
    }
  ]
}
```

---

## 六、客服模块 (support)

### 6.1 客服会话

#### 25. 发起客服会话
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/support/create |
| **请求方式** | POST |
| **模块** | 客服 |

**请求参数：**
```json
{
  "order_id": "string",      // 关联订单ID（可选）
  "type": 1                  // 类型：1-行程问题，2-支付问题，3-其他
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "chat_id": 1
  }
}
```

---

#### 26. 发送消息
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/support/message |
| **请求方式** | POST |
| **模块** | 客服 |

**请求参数：**
```json
{
  "chat_id": 1,              // 会话ID，必填
  "content": "string"         // 消息内容，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "message_id": 1
  }
}
```

---

#### 27. 获取聊天记录
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/support/history/:chat_id |
| **请求方式** | GET |
| **模块** | 客服 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "chat_id": 1,
    "status": 1,
    "messages": [
      {
        "id": 1,
        "sender_type": 1,
        "sender_name": "我",
        "content": "string",
        "created_at": "2024-01-01T10:00:00Z"
      },
      {
        "id": 2,
        "sender_type": 2,
        "sender_name": "客服小猪",
        "content": "string",
        "created_at": "2024-01-01T10:01:00Z"
      }
    ]
  }
}
```

---

#### 28. 关闭会话
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/support/close/:chat_id |
| **请求方式** | PUT |
| **模块** | 客服 |

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

## 七、安全模块 (safety)

### 7.1 行程安全

#### 29. 行程分享
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/safety/share |
| **请求方式** | POST |
| **模块** | 行程安全 |

**请求参数：**
```json
{
  "order_id": "string",       // 订单ID，必填
  "contact_phone": "string"   // 紧急联系人电话，必填
}
```

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "share_url": "string"
  }
}
```

---

#### 30. 紧急求助
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/safety/emergency |
| **请求方式** | POST |
| **模块** | 行程安全 |

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
  "data": {
    "emergency_id": 1,
    "police_phone": "110",
    "help_status": 1
  }
}
```

---

#### 31. 安全记录列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/safety/logs |
| **请求方式** | GET |
| **模块** | 行程安全 |

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
        "log_type": 1,
        "log_type_name": "行程分享",
        "created_at": "2024-01-01T10:30:00Z"
      }
    ],
    "total": 10
  }
}
```

---

## 八、文章模块 (article)

### 8.1 文章公告

#### 32. 文章分类列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/article/categories |
| **请求方式** | GET |
| **模块** | 文章 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "name": "帮助中心"
    },
    {
      "id": 2,
      "name": "用户协议"
    }
  ]
}
```

---

#### 33. 文章列表
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/article/list |
| **请求方式** | GET |
| **模块** | 文章 |

**请求参数：**
```json
{
  "category_id": 1,
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
        "title": "如何预约用车",
        "category_id": 1,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 10
  }
}
```

---

#### 34. 文章详情
| 项目 | 内容 |
|------|------|
| **接口名** | /passenger/article/detail/:id |
| **请求方式** | GET |
| **模块** | 文章 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "title": "如何预约用车",
    "content": "string",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

## 九、公共模块 (common)

### 9.1 城市与车型

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
      "minute_price": 0.50,
      "description": "滴滴快车"
    }
  ]
}
```

---

### 9.2 基础信息

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

#### 38. 获取启动页/广告
| 项目 | 内容 |
|------|------|
| **接口名** | /common/splash |
| **请求方式** | GET |
| **模块** | 公共接口 |

**请求参数：** 无

**响应参数：**
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "image_url": "string",
    "link_url": "string",
    "duration": 3
  }
}
```

---

## 接口数量统计

| 模块 | 接口数量 |
|------|----------|
| 用户认证 | 3 |
| 个人信息 | 2 |
| 实名认证 | 1 |
| 地址簿管理 | 5 |
| 预估价格 | 1 |
| 创建订单 | 1 |
| 订单状态 | 3 |
| 历史订单 | 1 |
| 订单支付 | 1 |
| 钱包信息 | 1 |
| 充值 | 1 |
| 流水记录 | 1 |
| 优惠券领取 | 1 |
| 优惠券列表 | 2 |
| 客服会话 | 4 |
| 行程安全 | 3 |
| 文章分类 | 1 |
| 文章列表 | 2 |
| 文章详情 | 1 |
| 公共接口 | 3 |
| **总计** | **39** |
