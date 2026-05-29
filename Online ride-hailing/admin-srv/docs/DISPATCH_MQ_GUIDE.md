# RocketMQ 后台派单功能 - 完整实现文档

## 📊 业务流程图（已实现）

```
┌─────────────────────────────────────────────────┐
│  乘客下单（其他服务/接口处理）                      │
│  ↓                                              │
│  INSERT INTO passenger_orders                    │
│  (order_id, ..., driver_id=NULL, status=0)       │
└──────────────────┬──────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────┐
│  管理员后台选择订单 + 选择司机                     │
│  ↓                                              │
│  DispatchOrder API (参数校验)                    │
└──────────────────┬──────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────┐
│  发送派单消息到 RocketMQ 消息队列                 │
│  (producer.go: SendDispatchMessage)             │
└──────────────────┬──────────────────────────────┘
                   ↓
┌─────────────────────────────────────────────────┐
│  MQ消费者异步接收消息                             │
│  (consumer.go: handleDispatchMessage)           │
│         ↓                                       │
│    ┌─────┴─────┐                                 │
│    ↓           ↓                                 │
│ 校验状态     其他状态                              │
│ (status=0?)   (status≠0)                         │
│    ↓           ↓                                 │
│   ✅ 是        ❌ 否                              │
│    ↓           ↓                                 │
│ 更新订单      丢弃消息                            │
│ 推送给司机     派单失败                           │
└─────────────────────────────────────────────────┘
```

---

## 🔧 新增/修改的文件清单

### 1. **模型层** - `model/order/passenger_order.go`
```go
// ✅ 新增状态常量
const (
	OrderStatusPendingDispatch = 0 // 待派单（用户下单后等待后台分配司机）← NEW
	OrderStatusPendingAccept   = 1 // 待接单（已派单给司机，等待司机接单）
	// ... 其他状态保持不变
)
```

### 2. **配置层** - `global/config.go`
```go
type AppConfig struct {
	Mysql
	Redis
	RocketMQ // ← 改为使用RocketMQ
}

// ✅ 新增配置结构体
type RocketMQ struct {
	NameServers []string `mapstructure:"name_servers"` // ["127.0.0.1:9876"]
	GroupName   string   `mapstructure:"group_name"`   // "dispatch_group"
	Topic       string   `mapstructure:"topic"`        // "dispatch_topic"
}
```

### 3. **MQ生产者** - `pkg/mq/producer.go` (NEW)
**职责**: 发送派单消息到队列

```go
func InitProducer() error          // 初始化生产者（程序启动时调用）
func SendDispatchMessage(orderID string, driverID int64) error  // 发送派单消息
func CloseProducer() error          // 关闭生产者（程序退出时调用）
```

### 4. **MQ消费者** - `pkg/mq/consumer.go` (NEW)
**职责**: 异步消费并处理派单逻辑

```go
func StartDispatchConsumer() error  // 启动消费者（程序启动时调用）
func handleDispatchMessage(ctx, msg) error  // 核心业务处理函数：
                                           // 1. 解析消息
                                           // 2. 查询订单
                                           // 3. 校验status=0?
                                           // 4. 更新driver_id/dispatch_time/status=1
                                           // 5. 推送通知给司机端
```

### 5. **业务逻辑层** - `internal/logic/dispatchorderlogic.go`
**改动**: 从"直接创建订单"改为"参数校验 + 发送MQ消息"

```go
func (l *DispatchOrderLogic) DispatchOrder(in *admin.DispatchOrderRequest) {
    // 步骤1: 参数校验 (order_id, driver_id)
    // 步骤2: 调用 mq.SendDispatchMessage() → 发送到RocketMQ
    // 返回: 异步处理中的响应
}
```

### 6. **初始化入口**
- `inits/mq.go` (NEW) - 初始化 MQ 生产者和消费者
- `inits/init.go` - 添加 InitMQ() 调用

---

## ⚙️ Nacos/配置中心配置示例

在 Nacos 的配置文件中添加以下内容：

```yaml
# RocketMQ 配置（用于后台派单消息队列）
rocketmq:
  name_servers:
    - "127.0.0.1:9876"  # RocketMQ NameServer地址（可多个）
  group_name: "dispatch_group"   # 生产者/消费者组名
  topic: "dispatch_topic"        # 主题名称
```

或使用环境变量：

```bash
export ROCKETMQ_NAME_SERVERS="127.0.0.1:9876"
export ROCKETMQ_GROUP_NAME="dispatch_group"
export ROCKETMQ_TOPIC="dispatch_topic"
```

---

## 🚀 API 调用示例

### 请求
```json
POST /api/admin/DispatchOrder

{
  "order_id": "P20260525110101000000112000001",
  "driver_id": 12345
}
```

### 响应（同步返回）
```json
{
  "success": true,
  "message": "派单请求已提交，系统正在异步处理中",
  "order_id": "P20260525110101000000112000001"
}
```

> **注意**: 此时只是消息发送成功，实际派单操作由 MQ 消费者异步完成！

---

## 📋 数据库表结构变更

```sql
-- passenger_orders 表状态字段含义更新：
-- status = 0 : 待派单（用户刚下单，等待管理员分配司机）
-- status = 1 : 待接单（已分配司机，等待司机接单）
-- status = 2 : 已接单（司机已接单）
-- status = 3 : 司机已到达
-- status = 4 : 行程中
-- status = 5 : 已完成
-- status = 6 : 已取消
```

---

## 🔍 日志输出示例

### 程序启动时
```
✅ RocketMQ Producer 启动成功
✅ RocketMQ Dispatch Consumer 启动成功，等待接收派单消息...
✅ RocketMQ 消息队列初始化完成
```

### 派单API调用后
```
🎯 后台派单请求: order_id=P20260525..., driver_id=12345
📤 派单消息发送成功: order_id=P20260525..., driver_id=12345, msg_id=ABC123...
✅ 派单消息已发送到MQ队列: order_id=P20260525..., driver_id=12345 (等待异步处理)
```

### MQ消费者处理后
```
📥 接收到派单消息: order_id=P20260525..., driver_id=12345
✅ 后台派单处理完成: order_id=P20260525... → driver_id=12345, dispatch_time=2026-05-25 10:30:00
📤 推送派单消息给司机端: order_id=P20260525..., driver_id=12345
```

---

## ⚠️ 注意事项

1. **RocketMQ 服务必须先启动**：确保 RocketMQ Broker 和 NameServer 正在运行
2. **幂等性设计**：消费者中使用乐观锁 (`WHERE status=0 AND driver_id IS NULL`) 避免重复派单
3. **错误重试**：生产者配置了失败重试2次；消费者遇到数据库错误会自动重试
4. **优雅关闭**：程序退出时会正确关闭 Producer 和 Consumer

---

## 🛠️ 下一步优化建议

- [ ] 实现真正的司机端推送（WebSocket/长连接）
- [ ] 添加派单超时机制（如5分钟未接单则取消）
- [ ] 实现派单记录日志表（记录每次派单操作）
- [ ] 添加监控告警（派单成功率、平均耗时等指标）
