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

type GetCityListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCityListLogic {
	return &GetCityListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCityListLogic) GetCityList(in *admin.GetCityListRequest) (*admin.GetCityListResponse, error) {

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
	query := global.DB.Model(&system.CityConfig{})

	// 关键词搜索
	if in.Keyword != "" {
		query = query.Where("city_name LIKE ? OR city_code LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 查询列表
	var cities []system.CityConfig
	// 👇 这里修复！必须加 .Error ！！！
	err := query.Offset(int(offset)).Limit(int(pageSize)).Order("id ASC").Find(&cities).Error
	if err != nil {
		return nil, errors.New("查询列表失败")
	}

	// 转换结果
	list := make([]*admin.CityListItem, 0, len(cities))
	for _, city := range cities {
		item := &admin.CityListItem{
			Id:       int32(city.ID),
			CityCode: city.CityCode,
			CityName: city.CityName,
			Status:   int32(city.Status),
		}
		list = append(list, item)
	}

	return &admin.GetCityListResponse{
		Total: int32(total),
		List:  list,
	}, nil
}
