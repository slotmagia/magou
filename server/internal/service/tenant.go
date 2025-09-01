package service

import (
	"context"
	"client-app/internal/model/input/sysin"
	"client-app/internal/model/output/sysout"
)

// ITenant 租户服务接口
type ITenant interface {
	// GetTenantList 获取租户列表
	GetTenantList(ctx context.Context, in *sysin.TenantListInp) (*sysout.TenantListModel, error)
	
	// CreateTenant 创建租户
	CreateTenant(ctx context.Context, in *sysin.CreateTenantInp) (*sysout.TenantModel, error)
	
	// UpdateTenant 更新租户
	UpdateTenant(ctx context.Context, in *sysin.UpdateTenantInp) (*sysout.TenantModel, error)
	
	// DeleteTenant 删除租户
	DeleteTenant(ctx context.Context, in *sysin.DeleteTenantInp) error
	
	// GetTenantDetail 获取租户详情
	GetTenantDetail(ctx context.Context, in *sysin.TenantDetailInp) (*sysout.TenantDetailModel, error)
	
	// UpdateTenantStatus 更新租户状态
	UpdateTenantStatus(ctx context.Context, in *sysin.TenantStatusInp) error
	
	// GetTenantStats 获取租户统计信息
	GetTenantStats(ctx context.Context, in *sysin.TenantStatsInp) (*sysout.TenantStatsModel, error)
	
	// UpdateTenantConfig 更新租户配置
	UpdateTenantConfig(ctx context.Context, in *sysin.TenantConfigInp) error
	
	// GetTenantOptions 获取租户选项列表
	GetTenantOptions(ctx context.Context) (*sysout.TenantOptionsModel, error)
	
	// GetTenantByCode 根据编码获取租户信息
	GetTenantByCode(ctx context.Context, code string) (*sysout.TenantModel, error)
	
	// ValidateTenantAccess 验证租户访问权限
	ValidateTenantAccess(ctx context.Context, tenantId uint64) error
}

var localTenant ITenant

// Tenant 获取租户服务实例
func Tenant() ITenant {
	if localTenant == nil {
		panic("implement not found for interface ITenant, forgot register?")
	}
	return localTenant
}

// RegisterTenant 注册租户服务实现
func RegisterTenant(i ITenant) {
	localTenant = i
}
