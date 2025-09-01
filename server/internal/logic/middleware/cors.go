package middleware

import "github.com/gogf/gf/v2/net/ghttp"

//------------------------------------------

// CORS allows Cross-origin resource sharing.
func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
