package logic

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"user-srv/model"
	"user-srv/pkg"
	"user-srv/user"

	"user-srv/global"
	"user-srv/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddPassengerOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPassengerOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPassengerOrderLogic {
	return &AddPassengerOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddPassengerOrder 处理乘客下单请求
// 1. 验证请求参数（用户ID、地址、经纬度等必填项）
// 2. Redis分布式锁防重复下单（同一用户5秒内只能发起一次下单）
// 3. 验证用户是否存在
// 4. 计算预估价格和距离
// 5. 如使用优惠券，则验证优惠券有效性并计算优惠金额
// 6. 使用事务创建订单、更新优惠券状态、记录优惠券使用日志
// 7. 返回订单ID
func (l *AddPassengerOrderLogic) AddPassengerOrder(in *user.AddPassengerOrderReq) (*user.AddPassengerOrderResp, error) {
	if in.UserId == 0 {
		return &user.AddPassengerOrderResp{
			Code:    400,
			Message: "用户ID不能为空",
		}, nil
	}

	if in.StartAddress == "" || in.EndAddress == "" {
		return &user.AddPassengerOrderResp{
			Code:    400,
			Message: "出发地和目的地地址不能为空",
		}, nil
	}

	if in.StartLng == 0 || in.StartLat == 0 || in.EndLng == 0 || in.EndLat == 0 {
		return &user.AddPassengerOrderResp{
			Code:    400,
			Message: "经纬度坐标不能为零",
		}, nil
	}

	// Redis分布式锁：防止同一用户短时间内重复下单
	lockKey := fmt.Sprintf("order:create:%d", in.UserId)
	requestId := fmt.Sprintf("%d_%d", in.UserId, time.Now().UnixNano())
	distLock := pkg.NewRedisDistributedLock(global.RDB, l.ctx)

	ok, err := distLock.TryLock(lockKey, requestId, 10*time.Second, 5*time.Second)
	if err != nil || !ok {
		l.Infof("用户 %d 获取下单锁失败: err=%v ok=%v", in.UserId, err, ok)
		return &user.AddPassengerOrderResp{
			Code:    429,
			Message: "操作过于频繁，请稍后再试",
		}, nil
	}
	defer func() {
		if released, unlockErr := distLock.Unlock(lockKey, requestId); unlockErr != nil {
			l.Errorf("释放下单锁失败: key=%s err=%v", lockKey, unlockErr)
		} else {
			l.Infof("释放下单锁: key=%s released=%v", lockKey, released)
		}
	}()

	var passengerUser model.PassengerUser
	err = global.DB.Where("id = ?", in.UserId).First(&passengerUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &user.AddPassengerOrderResp{
				Code:    404,
				Message: "用户不存在",
			}, nil
		}
		l.Errorf("查询用户失败: %v", err)
		return &user.AddPassengerOrderResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	orderID := generateOrderID(in.UserId, in.CarType, float64(in.StartLng), float64(in.StartLat))
	carTypeStr := getCarTypeString(in.CarType)

	distance := calculateDistance(float64(in.StartLng), float64(in.StartLat), float64(in.EndLng), float64(in.EndLat))
	estimatedPrice := calculateEstimatedPrice(distance, in.CarType)

	var couponID int64
	var couponName string
	var reduceAmount float64

	if in.CouponId > 0 {
		reduceAmount, couponID, couponName, err = l.validateAndCalculateCoupon(in.UserId, in.CouponId, estimatedPrice, carTypeStr)
		if err != nil {
			return &user.AddPassengerOrderResp{
				Code:    400,
				Message: err.Error(),
			}, nil
		}
		estimatedPrice -= reduceAmount
		if estimatedPrice < 0 {
			estimatedPrice = 0
		}
	}

	passengerName := in.PassengerName
	if passengerName == "" {
		passengerName = passengerUser.RealName
	}
	if passengerName == "" {
		passengerName = passengerUser.Nickname
	}

	passengerPhone := in.PassengerPhone
	if passengerPhone == "" {
		passengerPhone = passengerUser.Phone
	}

	order := model.PassengerOrder{
		OrderID:        orderID,
		PassengerID:    in.UserId,
		PassengerName:  passengerName,
		PassengerPhone: passengerPhone,
		CouponID:       couponID,
		CouponName:     couponName,
		OrderType:      1,
		CarType:        carTypeStr,
		Status:         1,
		StartLng:       float64(in.StartLng),
		StartLat:       float64(in.StartLat),
		StartAddress:   in.StartAddress,
		EndLng:         float64(in.EndLng),
		EndLat:         float64(in.EndLat),
		EndAddress:     in.EndAddress,
		PassRemark:     in.Remark,
		EstimatedPrice: estimatedPrice,
		BookTime:       time.Now(),
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Create(&order).Error
	if err != nil {
		l.Errorf("创建订单失败: %v", err)
		return &user.AddPassengerOrderResp{
			Code:    500,
			Message: "系统错误",
		}, nil
	}

	if couponID > 0 {
		err = tx.Model(&model.UserCoupon{}).
			Where("id = ?", couponID).
			Updates(map[string]interface{}{
				"status":    2,
				"used_time": time.Now(),
				"order_no":  orderID,
			}).Error
		if err != nil {
			l.Errorf("更新优惠券状态失败: %v", err)
			tx.Rollback()
			return &user.AddPassengerOrderResp{
				Code:    500,
				Message: "系统错误",
			}, nil
		}

		couponUseLog := model.CouponUseLog{
			UserID:       in.UserId,
			TemplateID:   0,
			UserCouponID: couponID,
			OrderNo:      orderID,
			OrderAmount:  estimatedPrice + reduceAmount,
			ReduceAmount: reduceAmount,
		}
		err = tx.Create(&couponUseLog).Error
		if err != nil {
			l.Errorf("创建优惠券使用记录失败: %v", err)
			tx.Rollback()
			return &user.AddPassengerOrderResp{
				Code:    500,
				Message: "系统错误",
			}, nil
		}
	}

	tx.Commit()
	alipay := pkg.Alipay(orderID, estimatedPrice)
	return &user.AddPassengerOrderResp{
		Code:      200,
		Message:   "下单成功",
		AlipayUrl: alipay,
	}, nil
}

// validateAndCalculateCoupon 验证优惠券的有效性并计算优惠金额
// 验证规则：
//   - 优惠券是否存在
//   - 优惠券是否属于当前用户
//   - 优惠券状态是否为未使用
//   - 优惠券是否在有效期内
//   - 订单金额是否达到使用门槛
//   - 优惠券是否适用于当前车型
//
// 返回：优惠金额、优惠券ID、优惠券名称、错误信息
func (l *AddPassengerOrderLogic) validateAndCalculateCoupon(userID, couponID int64, orderAmount float64, carType string) (float64, int64, string, error) {
	var userCoupon model.UserCoupon
	err := global.DB.Where("id = ?", couponID).First(&userCoupon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, "", errors.New("优惠券不存在")
		}
		return 0, 0, "", errors.New("系统错误")
	}

	if userCoupon.UserID != userID {
		return 0, 0, "", errors.New("优惠券不属于当前用户")
	}

	if userCoupon.Status != 1 {
		return 0, 0, "", errors.New("优惠券已使用或已失效")
	}

	now := time.Now()
	if now.Before(userCoupon.StartTime) || now.After(userCoupon.EndTime) {
		return 0, 0, "", errors.New("优惠券不在有效期内")
	}

	var template model.CouponTemplate
	err = global.DB.Where("id = ?", userCoupon.TemplateID).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, "", errors.New("优惠券模板不存在")
		}
		return 0, 0, "", errors.New("系统错误")
	}

	if template.Status != 1 {
		return 0, 0, "", errors.New("优惠券模板已下架")
	}

	if orderAmount < template.MinAmount {
		return 0, 0, "", fmt.Errorf("订单金额未达到优惠券使用门槛（需满%.2f元）", template.MinAmount)
	}

	if template.CarType != "" {
		carTypes := strings.Split(template.CarType, ",")
		found := false
		for _, ct := range carTypes {
			if strings.TrimSpace(ct) == carType {
				found = true
				break
			}
		}
		if !found {
			return 0, 0, "", errors.New("该优惠券不适用于当前车型")
		}
	}

	reduceAmount := calculateCouponDiscount(orderAmount, &template)

	return reduceAmount, couponID, template.Name, nil
}

// calculateCouponDiscount 根据优惠券模板类型计算优惠金额
// 优惠券类型说明：
//   - 类型1（满减）：订单金额满足门槛时，直接减免指定金额
//   - 类型2（折扣）：按折扣比例计算减免金额，最高不超过max_reduce
//   - 类型3（立减）：直接减免指定金额，不考虑门槛
//   - 类型4（新人券）：满足门槛后减免指定金额
//
// 返回：减免金额
func calculateCouponDiscount(orderAmount float64, template *model.CouponTemplate) float64 {
	switch template.Type {
	case 1:
		if orderAmount >= template.MinAmount {
			return template.ReduceAmount
		}
	case 2:
		if orderAmount >= template.MinAmount {
			discount := orderAmount * (1 - template.Discount/10)
			if discount > template.MaxReduce {
				return template.MaxReduce
			}
			return discount
		}
	case 3:
		return template.ReduceAmount
	case 4:
		if orderAmount >= template.MinAmount {
			return template.ReduceAmount
		}
	}
	return 0
}

// generateOrderID 生成可追溯的网约车订单号
// 订单号格式：P + 日期(8位) + 区域码(4位) + 用户ID后6位 + 车型代码(1位) + 时间(6位) + 序列号(4位) + 纳秒后4位
// 总长度：34位
//
// 参数说明：
//   - userID：用户ID，用于嵌入订单号便于追溯
//   - carType：车型代码，用于区分订单类型
//   - startLng：出发地经度，用于计算区域码
//   - startLat：出发地纬度，用于计算区域码
//
// 示例：P2026010531AF012345E14302500011234
func generateOrderID(userID int64, carType int32, startLng, startLat float64) string {
	now := time.Now()
	dateStr := now.Format("20060102")
	hmsStr := now.Format("150405")
	millis := strconv.FormatInt(time.Now().UnixNano()%100000000, 10)
	if len(millis) < 8 {
		millis = fmt.Sprintf("%08s", millis)
	}
	userIDStr := fmt.Sprintf("%06d", userID%1000000)
	carTypeChar := getCarTypeChar(carType)
	regionCode := calculateRegionCode(startLng, startLat)
	sequence := strconv.FormatInt(int64(time.Now().UnixNano()%10000), 10)
	if len(sequence) < 4 {
		sequence = fmt.Sprintf("%04s", sequence)
	}
	return fmt.Sprintf("P%s%s%s%s%s%s%s",
		dateStr,
		regionCode,
		userIDStr,
		carTypeChar,
		hmsStr,
		sequence,
		millis[:4])
}

// getCarTypeChar 将车型数字代码转换为单字符标识
// 车型代码对照：
//   - 1（经济型） -> E
//   - 2（舒适型） -> C
//   - 3（商务型） -> B
//   - 4（豪华型） -> L
//   - 其他 -> E（默认为经济型）
//
// 返回：车型对应的单字符代码
func getCarTypeChar(carType int32) string {
	switch carType {
	case 1:
		return "E"
	case 2:
		return "C"
	case 3:
		return "B"
	case 4:
		return "L"
	default:
		return "E"
	}
}

// calculateRegionCode 根据经纬度坐标计算区域码
// 使用经纬度的整数部分模256得到一个16进制的区域标识
// 可用于大致区分不同城市或城区
//
// 参数说明：
//   - lng：经度
//   - lat：纬度
//
// 返回：2位十六进制区域码
func calculateRegionCode(lng, lat float64) string {
	lngInt := int64((lng + 180) * 10000)
	latInt := int64((lat + 90) * 10000)
	return fmt.Sprintf("%02X%02X", latInt%256, lngInt%256)
}

// parseOrderID 解析订单号，提取各部分信息
// 将34位订单号分解为日期、区域码、用户ID、车型代码、时间、序列号、毫秒等部分
//
// 参数说明：
//   - orderID：订单号字符串
//
// 返回：包含各部分信息的map，key包括date/region_code/user_id/car_type/time/sequence/millis
func parseOrderID(orderID string) (map[string]string, error) {
	result := make(map[string]string)
	if len(orderID) < 20 {
		return nil, errors.New("订单号格式不正确")
	}
	if orderID[0] != 'P' {
		return nil, errors.New("订单号不是有效的乘客订单")
	}
	result["date"] = orderID[1:9]
	result["region_code"] = orderID[9:13]
	result["user_id"] = orderID[13:19]
	result["car_type"] = orderID[19:20]
	result["time"] = orderID[20:26]
	result["sequence"] = orderID[26:30]
	result["millis"] = orderID[30:34]
	return result, nil
}

// TraceOrderID 将订单号转换为可读的可追溯信息
// 用于日志追踪和问题排查，可显示订单的完整解析信息
//
// 参数说明：
//   - orderID：订单号字符串
//
// 返回：格式化的订单追溯信息，包含订单号、日期、时间、用户ID、车型、序列号等
//
//	如订单号格式错误，则返回错误信息
func TraceOrderID(orderID string) string {
	info, err := parseOrderID(orderID)
	if err != nil {
		return fmt.Sprintf("无效订单号: %s, 错误: %v", orderID, err)
	}
	carTypeMap := map[string]string{
		"E": "经济型",
		"C": "舒适型",
		"B": "商务型",
		"L": "豪华型",
	}
	dateStr := info["date"]
	year := dateStr[:4]
	month := dateStr[4:6]
	day := dateStr[6:8]
	carType := info["car_type"]
	carTypeName := carTypeMap[carType]
	if carTypeName == "" {
		carTypeName = "未知"
	}
	userID := info["user_id"]
	timeStr := info["time"]
	hour := timeStr[:2]
	minute := timeStr[2:4]
	second := timeStr[4:6]
	sequence := info["sequence"]
	return fmt.Sprintf("订单号: %s | 日期: %s-%s-%s | 时间: %s:%s:%s | 用户ID: %s | 车型: %s(%s) | 序列号: %s",
		orderID, year, month, day, hour, minute, second, userID, carType, carTypeName, sequence)
}

// getCarTypeString 将车型数字代码转换为中文名称
// 车型代码对照：
//   - 1 -> 经济型
//   - 2 -> 舒适型
//   - 3 -> 商务型
//   - 4 -> 豪华型
//   - 其他 -> 经济型（默认为经济型）
//
// 返回：车型对应的中文名称
func getCarTypeString(carType int32) string {
	switch carType {
	case 1:
		return "经济型"
	case 2:
		return "舒适型"
	case 3:
		return "商务型"
	case 4:
		return "豪华型"
	default:
		return "经济型"
	}
}

// calculateDistance 使用Haversine公式计算两点间的球面距离
// 根据地球表面的两个经纬度坐标点，计算它们之间的直线距离
// 常用于计算起点到终点的预估里程
//
// 参数说明：
//   - startLng：起点经度
//   - startLat：起点纬度
//   - endLng：终点经度
//   - endLat：终点纬度
//
// 返回：两点间的距离，单位为米
func calculateDistance(startLng, startLat, endLng, endLat float64) float64 {
	const earthRadius = 6371.0

	startRadLat := startLat * math.Pi / 180
	endRadLat := endLat * math.Pi / 180
	deltaLat := (endLat - startLat) * math.Pi / 180
	deltaLng := (endLng - startLng) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(startRadLat)*math.Cos(endRadLat)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c * 1000
}

// calculateEstimatedPrice 根据行驶距离和车型计算预估价格
// 采用起步价 + 里程费的计费模式
// 起步价内包含3公里，超出部分按每公里单价计费
//
// 车型对应价格：
//   - 经济型：起步价8元，每公里1.5元
//   - 舒适型：起步价10元，每公里2元
//   - 商务型：起步价15元，每公里3元
//   - 豪华型：起步价20元，每公里4.5元
//
// 参数说明：
//   - distance：行驶距离，单位为米
//   - carType：车型代码（1-经济型，2-舒适型，3-商务型，4-豪华型）
//
// 返回：预估价格，单位为元
func calculateEstimatedPrice(distance float64, carType int32) float64 {
	var basePrice float64
	var perKmPrice float64

	switch carType {
	case 1:
		basePrice = 8.0
		perKmPrice = 1.5
	case 2:
		basePrice = 10.0
		perKmPrice = 2.0
	case 3:
		basePrice = 15.0
		perKmPrice = 3.0
	case 4:
		basePrice = 20.0
		perKmPrice = 4.5
	default:
		basePrice = 8.0
		perKmPrice = 1.5
	}

	km := distance / 1000.0
	if km < 3 {
		return basePrice
	}

	return basePrice + (km-3)*perKmPrice
}
