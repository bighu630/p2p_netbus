package web

import (
	"context"
	"net/http"
	"server/log"
	"server/web/handler"

	"github.com/gin-gonic/gin"
)

type WebListener struct {
	server *http.Server
}

func NewGinListener() *WebListener {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	server := &WebListener{
		&http.Server{
			Addr:    ":7709",
			Handler: router,
		},
	}
	return server
}

func (wl *WebListener) Start() error {
	// 初始化消息处理器
	// 启动Http服务
	// service connections
	if TLS {
		if err := wl.server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Error("Web Server Listen", "err", err)
		}
	} else {
		if err := wl.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Web Server Listen", "err", err)
		}
	}
	return nil
}

func (wl *WebListener) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), WebTimeOut)
	defer cancel()
	if err := wl.server.Shutdown(ctx); err != nil {
		log.Error("Web Server Shutdown", "err", err)
	}
	log.Info("Module web-listener stopped")
	return nil
}

func initRouter(router *gin.Engine) {
	group := router.Group(ApiTag)
	initApi(group)
}

func initApi(group *gin.RouterGroup) {
	group.POST("HeartBeat", handler.Entrance)
}
