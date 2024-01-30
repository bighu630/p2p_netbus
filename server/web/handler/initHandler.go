package handler

import (
	"net/http"
	"server/log"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContextHandler interface {

	// Handle 处理
	Handle(ctx *gin.Context)
}

var handlerMap = make(map[string]ContextHandler)

func InitWebHandler() {
	handlerMap["HeartBeat"] = NewHeartBeatHandler()

}

// Entrance 请求处理入口
func Entrance(ctx *gin.Context) {
	contextHandler := ParseUrl(ctx)
	if contextHandler == nil {
		log.Error("can not find context handler to handle request")
		// 返回错误信息
		jsonResponse(ctx, http.StatusNotImplemented, "can not find context handler to handle request")
		return
	}
	contextHandler.Handle(ctx)
}

// jsonResponse wrapper the response
func jsonResponse(ctx *gin.Context, httpStatus int, data interface{}) {
	ctx.JSON(httpStatus, data)
}

// ParseUrl load method of this request
func ParseUrl(ctx *gin.Context) ContextHandler {
	log.Info("Receive http request", "url", ctx.Request.URL.String())
	// 获取方法
	method := urlToMethod(ctx.Request.URL.String())

	if handler, exist := handlerMap[method]; exist {
		return handler
	} else {
		log.Error("not support the method")
		return nil
	}
}

// 提取url最后路径名
func urlToMethod(url string) string {
	if url == "" {
		return url
	}
	u := strings.Split(url, "/")
	uu := strings.Split(u[len(u)-1], "?")
	return uu[0]
}
