package main

import (
	"server/log"
	"server/web"
)

func main() {
	log.NewZapLogger(log.ZapConf{})
	gin := web.NewGinListener()
	gin.Start()
}
