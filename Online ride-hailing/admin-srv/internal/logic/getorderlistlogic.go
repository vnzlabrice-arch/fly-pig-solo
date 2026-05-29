package logic

import (
	"context"
	"errors"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/model/order"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderListLogic {
	return &GetOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrderList 查询订单列表（支持分页、关键词搜索、状态筛选）
func (l *GetOrderListLogic) GetOrderList(in *admin.GetOrderListRequest) (*admin.GetOrderListResponse, error) {
	if global.DB == nil {
		return &admin.GetOrderListResponse{
			Code:    500,
			Message: "数据库未初始化",
			Total:   0,
			List:    nil,
		}, nil
	}

	// 设置默认分页参数
	page := in.Page
	pageSize := in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 构建查询条件
	query := global.DB.Model(&order.PassengerOrder{})

	// 状态筛选（0=全部，不筛选）
	if in.Status > 0 {
		query = query.Where("status = ?", in.Status)
	}

	// 关键词搜索：订单号/乘客电话/司机ID模糊匹配
	if in.Keyword != "" {
		query = query.Where("order_id LIKE ? OR passenger_phone LIKE ? OR driver_id LIKE ?",
			"%"+in.Keyword+"%", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	}

	// 查询总数
	var total int64
	query.Count(&total)

	// 分页查询
	var orders []order.PassengerOrder
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&orders).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.Errorf("查询订单列表失败: %v", err)
		return &admin.GetOrderListResponse{
			Code:    500,
			Message: "查询订单失败",
			Total:   0,
			List:    nil,
		}, nil
	}

	// 转换为响应结构
	list := make([]*admin.OrderListItem, 0, len(orders))
	for _, o := range orders {
		list = append(list, &admin.OrderListItem{
			OrderId:        o.OrderID,
			PassengerName:  o.PassengerName,
			PassengerPhone: maskPhone(o.PassengerPhone),
			CarType:        o.CarType,
			StartAddress:   o.StartAddress,
			EndAddress:     o.EndAddress,
			Status:         int32(o.Status),
			EstimatedPrice: o.EstimatedPrice,
			FinalPrice:     o.FinalPrice,
			CreatedAt:      o.CreatedAt.Unix(),
		})
	}

	return &admin.GetOrderListResponse{
		Code:    200,
		Message: "success",
		Total:   total,
		List:    list,
	}, nil
}

// maskPhone 手机号脱敏：138****1234
func maskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}
