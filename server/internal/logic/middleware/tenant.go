package middleware

import (
	"client-app/internal/library/contexts"
	"client-app/internal/library/response"
	"client-app/internal/model"
	"client-app/internal/model/entity"
	"client-app/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// TenantFilter 租户过滤中间件
func (s *sMiddleware) TenantFilter(r *ghttp.Request) {
	var (
		ctx      = r.Context()
		tenantId uint64
		err      error
	)

	// 获取租户ID的多种方式
	// 1. 从Header中获取
	tenantIdStr := r.Header.Get("X-Tenant-Id")
	if tenantIdStr != "" {
		tenantId = gconv.Uint64(tenantIdStr)
	}

	// 2. 从域名中获取
	if tenantId == 0 {
		host := r.Host
		// 移除端口号
		if gstr.Contains(host, ":") {
			host = gstr.SubStrRune(host, 0, gstr.PosRune(host, ":"))
		}

		// 根据域名查询租户
		var tenant *entity.Tenant
		err = g.DB().Model("tenants").Where("domain", host).Where("deleted_at IS NULL").Scan(&tenant)
		if err == nil && tenant != nil {
			tenantId = tenant.Id
		}
	}

	// 3. 从URL参数中获取（开发环境使用）
	if tenantId == 0 {
		tenantIdStr = r.Get("tenant_id").String()
		if tenantIdStr != "" {
			tenantId = gconv.Uint64(tenantIdStr)
		}
	}

	// 4. 使用默认租户（系统租户）
	if tenantId == 0 {
		tenantId = 1 // 系统租户ID
	}

	// 验证租户访问权限
	err = service.Tenant().ValidateTenantAccess(ctx, tenantId)
	if err != nil {
		response.JsonExit(r, 400, err.Error())
		return
	}

	// 设置租户上下文
	customCtx := &model.Context{
		Data: g.Map{
			"tenantId": tenantId,
		},
	}
	contexts.Init(r, customCtx)

	r.Middleware.Next()
}

// TenantAuth 租户权限验证中间件
func (s *sMiddleware) TenantAuth(r *ghttp.Request) {
	var (
		ctx = r.Context()
	)

	// 获取当前用户信息
	customCtx := contexts.Get(ctx)
	if customCtx == nil || customCtx.User == nil {
		response.JsonExit(r, 401, "请先登录")
		return
	}

	// 获取租户ID
	tenantId := gconv.Uint64(customCtx.Data["tenantId"])
	if tenantId == 0 {
		response.JsonExit(r, 400, "租户信息错误")
		return
	}

	// 验证用户是否属于当前租户
	if customCtx.User.TenantId != int64(tenantId) {
		// 检查是否为系统管理员（可以访问所有租户）
		if !customCtx.User.IsSystemAdmin() {
			response.JsonExit(r, 403, "无权限访问该租户")
			return
		}
	}

	// 更新用户上下文中的租户信息
	if customCtx.User.TenantId == 0 {
		customCtx.User.TenantId = int64(tenantId)
	}

	r.Middleware.Next()
}

// GetCurrentTenantId 获取当前请求的租户ID
func GetCurrentTenantId(r *ghttp.Request) uint64 {
	customCtx := contexts.Get(r.Context())
	if customCtx == nil {
		return 1 // 默认系统租户
	}
	return gconv.Uint64(customCtx.Data["tenantId"])
}

// SetTenantCondition 为查询条件添加租户过滤
func SetTenantCondition(r *ghttp.Request, model *gdb.Model) *gdb.Model {
	tenantId := GetCurrentTenantId(r)
	return model.Where("tenant_id", tenantId)
}

// IsTenantAdmin 判断当前用户是否为租户管理员
func IsTenantAdmin(r *ghttp.Request) bool {
	customCtx := contexts.Get(r.Context())
	if customCtx == nil || customCtx.User == nil {
		return false
	}
	return customCtx.User.IsTenantAdmin()
}

// IsSystemAdmin 判断当前用户是否为系统管理员
func IsSystemAdmin(r *ghttp.Request) bool {
	customCtx := contexts.Get(r.Context())
	if customCtx == nil || customCtx.User == nil {
		return false
	}
	return customCtx.User.IsSystemAdmin()
}
