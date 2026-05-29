package logic

import (
	"context"
	"errors"

	"admin-srv/internal/svc"
	"admin-srv/pb/admin"
	"admin-srv/pkg/mq"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchOrderLogic {
	return &DispatchOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DispatchOrder 后台派单接口（修改模式）
//
// 业务流程（按老师要求）:
// 1. 管理员后台选择订单(已存在, status=0待派单) + 选择司机
// 2. 派单API参数校验 (本函数)
// 3. 发送派单消息到 RocketMQ 消息队列 (本函数)
// 4. MQ消费者异步接收消息 (consumer.go)
// 5. 校验订单状态是否是待派单(status=0)
// 6. 更新订单: driver_id / dispatch_time / status=1(待接单)
// 7. 推送派单消息给司机端
func (l *DispatchOrderLogic) DispatchOrder(in *admin.DispatchOrderRequest) (*admin.DispatchOrderResponse, error) {
	// ====== 步骤1: 参数校验 ======
	if in.OrderId == "" {
		return nil, errors.New("订单号不能为空")
	}
	if in.DriverId <= 0 {
		return nil, errors.New("请指定目标司机ID")
	}

	l.Infof("🎯 后台派单请求: order_id=%s, driver_id=%d", in.OrderId, in.DriverId)

	// ====== 步骤2: 发送派单消息到 MQ 消息队列 ======
	err := mq.SendDispatchMessage(in.OrderId, in.DriverId)
	if err != nil {
		l.Errorf("❌ 发送派单消息失败: %v", err)
		return &admin.DispatchOrderResponse{
			Success: false,
			Message: "派单请求已提交但消息队列异常，请稍后查看结果",
			OrderId: in.OrderId,
		}, nil
	}

	// ====== 返回成功响应（异步处理中）=====
	l.Infof("✅ 派单消息已发送到MQ队列: order_id=%s, driver_id=%d (等待异步处理)",
		in.OrderId, in.DriverId)

	return &admin.DispatchOrderResponse{
		Success: true,
		Message: "派单请求已提交，系统正在处理中",
		OrderId: in.OrderId,
	}, nil
}
