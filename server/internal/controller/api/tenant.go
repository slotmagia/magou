package api

import (
	"client-app/internal/api/v1/tenant"
	"client-app/internal/service"
	"context"
)

// Tenant 租户控制器
type Tenant struct{}

// NewTenant 创建租户控制器实例
func NewTenant() *Tenant {
	return &Tenant{}
}

// GetTenantList 获取租户列表
func (c *Tenant) GetTenantList(ctx context.Context, req *tenant.TenantListReq) (res *tenant.TenantListRes, err error) {
	out, err := service.Tenant().GetTenantList(ctx, &req.TenantListInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantListRes{
		TenantListModel: out,
	}
	return res, nil
}

// CreateTenant 创建租户
func (c *Tenant) CreateTenant(ctx context.Context, req *tenant.CreateTenantReq) (res *tenant.CreateTenantRes, err error) {
	out, err := service.Tenant().CreateTenant(ctx, &req.CreateTenantInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.CreateTenantRes{
		TenantModel: out,
	}
	return res, nil
}

// UpdateTenant 更新租户
func (c *Tenant) UpdateTenant(ctx context.Context, req *tenant.UpdateTenantReq) (res *tenant.UpdateTenantRes, err error) {
	out, err := service.Tenant().UpdateTenant(ctx, &req.UpdateTenantInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.UpdateTenantRes{
		TenantModel: out,
	}
	return res, nil
}

// DeleteTenant 删除租户
func (c *Tenant) DeleteTenant(ctx context.Context, req *tenant.DeleteTenantReq) (res *tenant.DeleteTenantRes, err error) {
	err = service.Tenant().DeleteTenant(ctx, &req.DeleteTenantInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.DeleteTenantRes{}
	return res, nil
}

// GetTenantDetail 获取租户详情
func (c *Tenant) GetTenantDetail(ctx context.Context, req *tenant.TenantDetailReq) (res *tenant.TenantDetailRes, err error) {
	out, err := service.Tenant().GetTenantDetail(ctx, &req.TenantDetailInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantDetailRes{
		TenantDetailModel: out,
	}
	return res, nil
}

// UpdateTenantStatus 更新租户状态
func (c *Tenant) UpdateTenantStatus(ctx context.Context, req *tenant.TenantStatusReq) (res *tenant.TenantStatusRes, err error) {
	err = service.Tenant().UpdateTenantStatus(ctx, &req.TenantStatusInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantStatusRes{}
	return res, nil
}

// GetTenantStats 获取租户统计
func (c *Tenant) GetTenantStats(ctx context.Context, req *tenant.TenantStatsReq) (res *tenant.TenantStatsRes, err error) {
	out, err := service.Tenant().GetTenantStats(ctx, &req.TenantStatsInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantStatsRes{
		TenantStatsModel: out,
	}
	return res, nil
}

// UpdateTenantConfig 更新租户配置
func (c *Tenant) UpdateTenantConfig(ctx context.Context, req *tenant.TenantConfigReq) (res *tenant.TenantConfigRes, err error) {
	err = service.Tenant().UpdateTenantConfig(ctx, &req.TenantConfigInp)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantConfigRes{}
	return res, nil
}

// GetTenantOptions 获取租户选项
func (c *Tenant) GetTenantOptions(ctx context.Context, req *tenant.TenantOptionsReq) (res *tenant.TenantOptionsRes, err error) {
	out, err := service.Tenant().GetTenantOptions(ctx)
	if err != nil {
		return nil, err
	}

	res = &tenant.TenantOptionsRes{
		TenantOptionsModel: out,
	}
	return res, nil
}
