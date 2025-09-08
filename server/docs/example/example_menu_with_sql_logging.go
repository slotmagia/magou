package examples

import (
	"context"
	"time"

	"client-app/utility/logger"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// 模拟菜单服务，展示如何集成SQL日志功能
type MenuService struct{}

// GetMenuList 获取菜单列表（带SQL日志）
func (s *MenuService) GetMenuList(ctx context.Context, status int) ([]map[string]interface{}, error) {
	start := time.Now()

	// 模拟SQL查询
	sql := "SELECT `id`,`tenant_id`,`parent_id`,`menu_code`,`icon`,`path`,`component`,`permission`,`sort_order`,`visible`,`status`,`created_at`,`updated_at`,`redirect`,`active_menu`,`always_show`,`breadcrumb`,`remark`,`title`,`menu_type` FROM `sys_menus` WHERE (`status`=?) AND `deleted_at` IS NULL ORDER BY `sort_order` ASC,`id` ASC LIMIT 10"
	args := []interface{}{status}

	// 模拟数据库查询时间
	time.Sleep(25 * time.Millisecond)

	// 计算执行时间
	duration := time.Since(start)

	// 记录SQL日志（带调用者信息）
	logger.LogSQLWithCaller(ctx, sql, args, duration)

	// 模拟返回数据
	return []map[string]interface{}{
		{"id": 1, "title": "系统管理", "status": status},
		{"id": 2, "title": "用户管理", "status": status},
	}, nil
}

// GetMenuDetail 获取菜单详情（带SQL日志）
func (s *MenuService) GetMenuDetail(ctx context.Context, id int) (map[string]interface{}, error) {
	start := time.Now()

	// 模拟SQL查询
	sql := "SELECT `id`,`tenant_id`,`parent_id`,`menu_code`,`icon`,`path`,`component`,`permission`,`sort_order`,`visible`,`status`,`created_at`,`updated_at`,`redirect`,`active_menu`,`always_show`,`breadcrumb`,`remark`,`title`,`menu_type` FROM `sys_menus` WHERE (`id`=?) AND `deleted_at` IS NULL LIMIT 1"
	args := []interface{}{id}

	// 模拟数据库查询时间
	time.Sleep(15 * time.Millisecond)

	// 计算执行时间
	duration := time.Since(start)

	// 记录SQL日志（带调用者信息）
	logger.LogSQLWithCaller(ctx, sql, args, duration)

	// 模拟返回数据
	return map[string]interface{}{
		"id":    id,
		"title": "菜单详情",
		"path":  "/menu/detail",
	}, nil
}

// GetMenuTree 获取菜单树（带SQL日志）
func (s *MenuService) GetMenuTree(ctx context.Context) ([]map[string]interface{}, error) {
	start := time.Now()

	// 模拟SQL查询
	sql := "SELECT `id`,`tenant_id`,`parent_id`,`menu_code`,`icon`,`path`,`component`,`permission`,`sort_order`,`visible`,`status`,`created_at`,`updated_at`,`redirect`,`active_menu`,`always_show`,`breadcrumb`,`remark`,`title`,`menu_type` FROM `sys_menus` WHERE (`status`=1) AND (`menu_type` IN (1,2)) ORDER BY `sort_order` ASC,`id` ASC"
	args := []interface{}{}

	// 模拟数据库查询时间
	time.Sleep(35 * time.Millisecond)

	// 计算执行时间
	duration := time.Since(start)

	// 记录SQL日志（带调用者信息）
	logger.LogSQLWithCaller(ctx, sql, args, duration)

	// 模拟返回数据
	return []map[string]interface{}{
		{"id": 1, "title": "系统管理", "parent_id": 0},
		{"id": 2, "title": "用户管理", "parent_id": 1},
		{"id": 3, "title": "角色管理", "parent_id": 1},
	}, nil
}

// 模拟慢查询
func (s *MenuService) GetMenuWithSlowQuery(ctx context.Context) ([]map[string]interface{}, error) {
	start := time.Now()

	// 模拟复杂的SQL查询
	sql := "SELECT m.*, r.name as role_name FROM `sys_menus` m LEFT JOIN `sys_role_menus` rm ON m.id = rm.menu_id LEFT JOIN `sys_roles` r ON rm.role_id = r.id WHERE m.status = 1 AND r.status = 1 ORDER BY m.sort_order ASC, m.id ASC"
	args := []interface{}{}

	// 模拟慢查询
	time.Sleep(150 * time.Millisecond)

	// 计算执行时间
	duration := time.Since(start)

	// 记录SQL日志（带调用者信息）
	logger.LogSQLWithCaller(ctx, sql, args, duration)

	// 模拟返回数据
	return []map[string]interface{}{
		{"id": 1, "title": "系统管理", "role_name": "管理员"},
		{"id": 2, "title": "用户管理", "role_name": "管理员"},
	}, nil
}

func RunSQLLoggingDemo() {
	// 初始化上下文
	ctx := gctx.New()

	// 启用SQL演示功能
	logger.EnableSQLDemo()

	// 创建菜单服务
	menuService := &MenuService{}

	g.Log().Info(ctx, "开始演示SQL日志功能...")

	// 演示1: 查询菜单列表
	g.Log().Info(ctx, "演示1: 查询菜单列表")
	menus, err := menuService.GetMenuList(ctx, 1)
	if err != nil {
		g.Log().Errorf(ctx, "查询菜单列表失败: %v", err)
	} else {
		g.Log().Infof(ctx, "查询到 %d 条菜单记录", len(menus))
	}

	// 演示2: 查询菜单详情
	g.Log().Info(ctx, "演示2: 查询菜单详情")
	menu, err := menuService.GetMenuDetail(ctx, 1)
	if err != nil {
		g.Log().Errorf(ctx, "查询菜单详情失败: %v", err)
	} else {
		g.Log().Infof(ctx, "菜单详情: %v", menu)
	}

	// 演示3: 查询菜单树
	g.Log().Info(ctx, "演示3: 查询菜单树")
	tree, err := menuService.GetMenuTree(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "查询菜单树失败: %v", err)
	} else {
		g.Log().Infof(ctx, "查询到 %d 条菜单树记录", len(tree))
	}

	// 演示4: 慢查询
	g.Log().Info(ctx, "演示4: 慢查询（会显示红色警告）")
	slowResult, err := menuService.GetMenuWithSlowQuery(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "慢查询失败: %v", err)
	} else {
		g.Log().Infof(ctx, "慢查询完成，结果数量: %d", len(slowResult))
	}

	// 演示5: 直接使用SQL追踪器
	g.Log().Info(ctx, "演示5: 直接使用SQL追踪器")
	demonstrateDirectSQLTracing(ctx)

	g.Log().Info(ctx, "SQL日志功能演示完成")
}

// demonstrateDirectSQLTracing 演示直接使用SQL追踪器
func demonstrateDirectSQLTracing(ctx context.Context) {
	// 模拟不同的SQL查询
	queries := []struct {
		sql   string
		args  []interface{}
		delay time.Duration
	}{
		{
			sql:   "SELECT COUNT(*) FROM sys_menus WHERE status = ?",
			args:  []interface{}{1},
			delay: 5 * time.Millisecond,
		},
		{
			sql:   "INSERT INTO sys_menus (title, path, status) VALUES (?, ?, ?)",
			args:  []interface{}{"新菜单", "/new-menu", 1},
			delay: 8 * time.Millisecond,
		},
		{
			sql:   "UPDATE sys_menus SET title = ? WHERE id = ?",
			args:  []interface{}{"更新菜单", 1},
			delay: 12 * time.Millisecond,
		},
		{
			sql:   "DELETE FROM sys_menus WHERE id = ?",
			args:  []interface{}{999},
			delay: 3 * time.Millisecond,
		},
	}

	for i, query := range queries {
		start := time.Now()

		// 模拟查询执行时间
		time.Sleep(query.delay)

		duration := time.Since(start)

		// 记录SQL日志
		logger.LogSQLWithCaller(ctx, query.sql, query.args, duration)

		g.Log().Infof(ctx, "执行查询 %d 完成", i+1)
	}
}
