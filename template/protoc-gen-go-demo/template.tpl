var srv {{ $.InterfaceName }}

// HttpServer接口
type {{ $.InterfaceName }} interface {
{{range .MethodSet}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

type {{$.Name}} struct{
	server {{ $.InterfaceName }}
	router gin.IRouter
	resp  interface {
		Error(ctx *gin.Context, err error)
		ParamsError (ctx *gin.Context, err error)
		Success(ctx *gin.Context, data interface{})
	}
}

// Resp 返回值
type default{{$.Name}}Resp struct {}

func response(ctx *gin.Context, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code,
		"msg": msg,
		"data": data,
	})
}

// Error 返回错误信息
func Error(ctx *gin.Context, err error) {
	code := -1
	status := 500
	msg := "未知错误"

	if err == nil {
		msg += ", err is nil"
		response(ctx, status, code, msg, nil)
		return
	}

	type iCode interface{
		HTTPCode() int
		Message() string
		Code() int
	}

	var c iCode
	if errors.As(err, &c) {
		status = c.HTTPCode()
		code = c.Code()
		msg = c.Message()
	}

	_ = ctx.Error(err)

	response(ctx, status, code, msg, nil)
}
// Forbidden 返回禁止信息
func Forbidden(ctx *gin.Context, msg string) {
    code := -1
    response(ctx, http.StatusForbidden, code, msg, nil)
}

// ParamsError 参数错误
func ParamsError (ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	response(ctx, 400, 400, "参数错误", nil)
}

// Success 返回成功信息
func Success(ctx *gin.Context, data interface{}) {
	response(ctx, 200, 0, "成功", data)
}


{{range .Methods}}
func {{ .HandlerName }} (ctx *gin.Context) {
	var in {{.Request}}
{{if .HasPathParams }}
	if err := ctx.ShouldBindUri(&in); err != nil {
		ParamsError(ctx, err)
		return
	}
{{end}}
{{if eq .Method "GET" "DELETE" }}
	if err := ctx.ShouldBindQuery(&in); err != nil {
		ParamsError(ctx, err)
		return
	}
{{else if eq .Method "POST" "PUT" }}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ParamsError(ctx, err)
		return
	}
{{else}}
	if err := ctx.ShouldBind(&in); err != nil {
		ParamsError(ctx, err)
		return
	}
{{end}}
	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := srv.{{.Name}}(newCtx, &in)
	if err != nil {
		Error(ctx, err)
		return
	}

	Success(ctx, out)
}
{{end}}



func RegisterRouter(g *gin.Engine, svc {{ $.InterfaceName }}) *gin.Engine {
	srv = svc
	g.Use(gin.LoggerWithFormatter(LoggerFormat),gin.RecoveryWithWriter(nil, HandlerRecovery))
	m :=g.Group("/")
    m.Use(HandlerPermission)
	{{range .Methods}}
	        {{if eq .HandlerName "Login_0" }}
            g.Handle("{{.Method}}", "{{.Path}}", {{ .HandlerName }})
	        {{else}}
    		m.Handle("{{.Method}}", "{{.Path}}", {{ .HandlerName }})
    		{{end}}
    {{end}}
	return g
}
// 权限控制
var HandlerPermission func(g *gin.Context)
// recovery的配置
func HandlerRecovery(g *gin.Context, e interface{}) {
	if e != nil {
		g.AbortWithStatus(http.StatusBadRequest)
		_ = g.Error(errors.New(fmt.Sprintf("%v", e)))
	}
}


// 日志格式的配置
func LoggerFormat(params gin.LogFormatterParams) string {
	return fmt.Sprintf("lzy %v %d %s %s %s %v\n", params.TimeStamp, params.StatusCode, params.ClientIP, params.Method, params.Path, params.ErrorMessage)
}

