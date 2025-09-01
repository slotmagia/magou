// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Tenant is the golang structure for table tenants.
type Tenant struct {
	Id           uint64      `json:"id"           description:"租户ID"`
	Name         string      `json:"name"         description:"租户名称"`
	Code         string      `json:"code"         description:"租户编码"`
	Domain       string      `json:"domain"       description:"租户域名"`
	Status       int         `json:"status"       description:"状态：1=正常 2=锁定 3=禁用"`
	MaxUsers     int         `json:"maxUsers"  db:"max_users"    description:"最大用户数"`
	StorageLimit int64       `json:"storageLimit" description:"存储限制(字节)"`
	ExpireAt     *gtime.Time `json:"expireAt"     description:"过期时间"`
	AdminUserId  uint64      `json:"adminUserId"  description:"租户管理员用户ID"`
	Config       string      `json:"config"       description:"租户配置"`
	Remark       string      `json:"remark"       description:"备注"`
	CreatedBy    uint64      `json:"createdBy"    description:"创建者"`
	UpdatedBy    uint64      `json:"updatedBy"    description:"更新者"`
	CreatedAt    *gtime.Time `json:"createdAt"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    description:"更新时间"`
	DeletedAt    *gtime.Time `json:"deletedAt"    description:"删除时间"`
}

// TenantConfig 租户配置结构
type TenantConfig struct {
	Features    map[string]bool `json:"features"`    // 功能特性开关
	Limitations map[string]int  `json:"limitations"` // 资源限制
	Settings    map[string]any  `json:"settings"`    // 自定义设置
}

// TenantStats 租户统计信息
type TenantStats struct {
	UserCount      int `json:"userCount"`      // 用户数量
	RoleCount      int `json:"roleCount"`      // 角色数量
	MenuCount      int `json:"menuCount"`      // 菜单数量
	StorageUsed    int64 `json:"storageUsed"`  // 已使用存储
	LastActiveTime *gtime.Time `json:"lastActiveTime"` // 最后活跃时间
}

// 租户状态常量
const (
	TenantStatusNormal   = 1 // 正常
	TenantStatusLocked   = 2 // 锁定
	TenantStatusDisabled = 3 // 禁用
)

// IsNormal 判断租户是否正常状态
func (t *Tenant) IsNormal() bool {
	return t.Status == TenantStatusNormal
}

// IsLocked 判断租户是否被锁定
func (t *Tenant) IsLocked() bool {
	return t.Status == TenantStatusLocked
}

// IsDisabled 判断租户是否被禁用
func (t *Tenant) IsDisabled() bool {
	return t.Status == TenantStatusDisabled
}

// IsExpired 判断租户是否过期
func (t *Tenant) IsExpired() bool {
	if t.ExpireAt == nil {
		return false
	}
	return gtime.Now().After(t.ExpireAt)
}

// IsSystemTenant 判断是否为系统租户
func (t *Tenant) IsSystemTenant() bool {
	return t.Code == "system"
}

// CanCreateUser 判断是否可以创建用户
func (t *Tenant) CanCreateUser(currentUserCount int) bool {
	return currentUserCount < t.MaxUsers
}

// CanUseStorage 判断是否还有存储空间
func (t *Tenant) CanUseStorage(usedStorage int64, needStorage int64) bool {
	return usedStorage+needStorage <= t.StorageLimit
}
