// Package main
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package main

import (
	"client-app/internal/cmd"
	"client-app/internal/global"
	_ "client-app/internal/logic"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var ctx = gctx.GetInitCtx()
	global.Init(ctx)
	cmd.Main.Run(ctx)
}
