// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

func init() {

}

const (
	HTTPContentTypeXml    = "text/xml"
	HTTPContentTypeHtml   = "text/html"
	HTTPContentTypeStream = "text/event-stream"
	HTTPContentTypeJson   = "application/json"
)

const (
	HTTPHandlerTime50          = "< 50"
	HTTPHandlerTime200         = "< 200"
	HTTPHandlerTime200To500    = "BETWEEN 200 AND 500"
	HTTPHandlerTime500To1000   = "BETWEEN 500 AND 1000"
	HTTPHandlerTime1000To10000 = "BETWEEN 1000 AND 10000"
	HTTPHandlerTime10000UP     = "> 10000"
)
