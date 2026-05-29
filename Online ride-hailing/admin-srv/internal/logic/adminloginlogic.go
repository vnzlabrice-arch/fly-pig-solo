package logic

import (
	"admin-srv/global"
	"admin-srv/pkg"
	"admin-srv/pkg/security"
	"admin-srv/model/system"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"admin-srv/internal/svc"
	"admin-srv/pb/admin"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AdminLogin 管理员登录（含防暴力破解功能）
//
// 功能说明:
// - 🔒 防暴力破解：连续失败5次锁定15分钟
func (l *AdminLoginLogic) AdminLogin(in *admin.AdminLoginRequest) (*admin.AdminLoginResponse, error) {
	// ========== 步骤1: 基础参数校验 ==========
	if in.Username == "" || in.Password == "" {
		return nil, errors.New("用户名或密码不能为空")
	}

	if global.DB == nil {
		return nil, errors.New("数据库未初始化")
	}

	// ========== 步骤2: 🔒 防暴力破解检查 ==========
	isLocked, remainingSec, err := security.CheckLoginLimit(l.ctx, in.Username)
	if err != nil {
		l.Errorf("Redis连接异常: %v", err)
	}
	if isLocked {
		lockTime := security.FormatLockTime(remainingSec)
		return nil, fmt.Errorf("登录失败次数过多，账号已锁定，请%s后重试", lockTime)
	}

	// ========== 步骤3: 查询管理员 ==========
	var adminUser system.AdminUser
	result := global.DB.Where("username = ?", in.Username).First(&adminUser)
	if result.Error != nil {
		// 用户不存在 - 统一提示，不暴露具体原因
		failCount, _ := security.RecordLoginFail(l.ctx, in.Username)
		remaining := security.MaxFailCount - failCount

		if failCount >= security.MaxFailCount {
			return nil, errors.New("用户名或密码错误，账号已被锁定")
		}

		return nil, fmt.Errorf("用户名或密码错误 (剩余%d次尝试机会)", remaining)
	}

	// ========== 步骤4: 验证密码 ==========
	if pkg.MD5(in.Password) != adminUser.Password {
		failCount, _ := security.RecordLoginFail(l.ctx, in.Username)
		remaining := security.MaxFailCount - failCount

		if failCount >= security.MaxFailCount {
			return nil, fmt.Errorf("密码错误，账号已锁定（%d次失败）", failCount)
		}

		return nil, fmt.Errorf("用户名或密码错误 (剩余%d次尝试机会)", remaining)
	}

	// ========== 步骤5: 检查账号状态 ==========
	if adminUser.Status != 1 {
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	// ========== 步骤6: 🔒 清除失败记录 + 生成Token ==========
	security.ClearLoginFail(l.ctx, in.Username)

	// 生成JWT Token
	token, err := pkg.TokenHandler(strconv.FormatInt(adminUser.ID, 10))
	if err != nil {
		l.Errorf("生成Token失败: %v", err)
		return nil, errors.New("登录失败，Token生成异常")
	}

	// ========== 步骤7: 更新最后登录时间 ==========
	now := time.Now()
	global.DB.Model(&adminUser).Updates(map[string]interface{}{
		"last_login_time": now,
	})

	l.Infof("✅ 登录成功: user=%s(%d)", adminUser.Username, adminUser.ID)

	return &admin.AdminLoginResponse{
		Token:    token,
		AdminId:  int32(adminUser.ID),
		Username: adminUser.Username,
	}, nil
}
