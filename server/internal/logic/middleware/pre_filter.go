package middleware

import (
	"client-app/internal/global"
	"client-app/internal/library/response"
	"client-app/utility/validate"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"reflect"
)

// PreFilter 请求输入预处理
// api使用gf规范路由并且XxxReq结构体实现了validate.Filter接口即可
func (s *sMiddleware) PreFilter(r *ghttp.Request) {
	router := global.GetRequestRoute(r)
	if router == nil {
		r.Middleware.Next()
		return
	}

	funcInfo := router.Handler.Info

	// 非规范路由不处理
	if funcInfo.Type.NumIn() != 2 {
		r.Middleware.Next()
		return
	}

	inputType := funcInfo.Type.In(1)
	var inputObject reflect.Value
	if inputType.Kind() == reflect.Ptr {
		inputObject = reflect.New(inputType.Elem())
	} else {
		inputObject = reflect.New(inputType.Elem()).Elem()
	}

	// 先验证基本校验规则
	if err := r.Parse(inputObject.Interface()); err != nil {
		resp := gerror.Code(err)
		response.JsonExit(r, resp.Code(), gerror.Current(err).Error(), resp.Detail())
		return
	}

	// 没有实现预处理
	if _, ok := inputObject.Interface().(validate.Filter); !ok {
		r.Middleware.Next()
		return
	}

	// 执行预处理
	if err := validate.PreFilter(r.Context(), inputObject.Interface()); err != nil {
		resp := gerror.Code(err)
		response.JsonExit(r, resp.Code(), gerror.Current(err).Error(), resp.Detail())
		return
	}

	// 过滤后的参数
	r.SetParamMap(gconv.Map(inputObject.Interface()))
	r.Middleware.Next()
}
