package sysout

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
)

// TenantModel 租户模型
type TenantModel struct {
	Id           uint64      `json:"id"`           // 租户ID
	Name         string      `json:"name"`         // 租户名称
	Code         string      `json:"code"`         // 租户编码
	Domain       string      `json:"domain"`       // 租户域名
	Status       int         `json:"status"`       // 状态：1=正常 2=锁定 3=禁用
	StatusName   string      `json:"statusName"`   // 状态名称
	MaxUsers     int         `json:"maxUsers"`     // 最大用户数
	StorageLimit int64       `json:"storageLimit"` // 存储限制(字节)
	ExpireAt     *gtime.Time `json:"expireAt"`     // 过期时间
	AdminUserId  uint64      `json:"adminUserId"`  // 租户管理员用户ID
	AdminName    string      `json:"adminName"`    // 管理员用户名
	Config       interface{} `json:"config"`       // 租户配置
	Remark       string      `json:"remark"`       // 备注
	CreatedAt    *gtime.Time `json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"`    // 更新时间
}

// TenantListModel 租户列表响应模型
type TenantListModel struct {
	List     []*TenantModel `json:"list"`     // 租户列表
	Total    int64          `json:"total"`    // 总数
	Page     int            `json:"page"`     // 当前页码
	PageSize int            `json:"pageSize"` // 每页数量
}

// TenantDetailModel 租户详情模型
type TenantDetailModel struct {
	*TenantModel
	Stats *TenantStatsModel `json:"stats"` // 统计信息
}

// TenantStatsModel 租户统计模型
type TenantStatsModel struct {
	UserCount       int         `json:"userCount"`       // 用户数量
	RoleCount       int         `json:"roleCount"`       // 角色数量
	MenuCount       int         `json:"menuCount"`       // 菜单数量
	StorageUsed     int64       `json:"storageUsed"`     // 已使用存储
	StorageUsedText string      `json:"storageUsedText"` // 已使用存储文本
	LastActiveTime  *gtime.Time `json:"lastActiveTime"`  // 最后活跃时间
}

// TenantOptionModel 租户选项模型
type TenantOptionModel struct {
	Value uint64 `json:"value"` // 租户ID
	Label string `json:"label"` // 租户名称
	Code  string `json:"code"`  // 租户编码
}

// TenantOptionsModel 租户选项列表模型
type TenantOptionsModel struct {
	List []*TenantOptionModel `json:"list"` // 租户选项列表
}

// 获取状态名称
func GetTenantStatusName(status int) string {
	switch status {
	case 1:
		return "正常"
	case 2:
		return "锁定"
	case 3:
		return "禁用"
	default:
		return "未知"
	}
}

// 格式化存储大小
func FormatStorageSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return "0 B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
