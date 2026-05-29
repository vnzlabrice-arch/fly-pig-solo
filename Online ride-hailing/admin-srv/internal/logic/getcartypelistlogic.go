package logic

import (
	"context"
	"errors"

	"admin-srv/global"
	"admin-srv/internal/svc"
	"admin-srv/model/system"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCarTypeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCarTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCarTypeListLogic {
	return &GetCarTypeListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCarTypeListLogic) GetCarTypeList(in *admin.GetCarTypeListRequest) (*admin.GetCarTypeListResponse, error) {
	// 默认分页参数
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// 构建查询
	query := global.DB.Model(&system.CarTypeConfig{})

	// 统计总数
	var total int64
	query.Count(&total)

	// 查询列表
	var carTypes []system.CarTypeConfig
	err := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&carTypes).Error
	if err != nil {
		return nil, errors.New("查询车型列表失败")
	}

	// 转换结果
	list := make([]*admin.CarTypeListItem, 0, len(carTypes))
	for _, carType := range carTypes {
		item := &admin.CarTypeListItem{
			Id:          int32(carType.ID),
			TypeName:    carType.TypeName,
			BasePrice:   carType.BasePrice,
			KmPrice:     carType.KmPrice,
			MinutePrice: carType.MinutePrice,
			Status:      int32(carType.Status),
		}
		list = append(list, item)
	}

	return &admin.GetCarTypeListResponse{
		Total: int32(total),
		List:  list,
	}, nil
}
