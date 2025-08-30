package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huangrao121/CommunicationApp/BackendService/internal/common/middleware"
)

func InitializeRouter(handlerInit *HandlerInit) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.SecureHeaders())

	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.CORS())
	}
	r.Use(gin.Recovery())

	// 设置信任代理，如果在生产环境中，使用负载均衡和反向代理
	// if gin.Mode() == gin.ReleaseMode {
	// 	r.SetTrustedProxies([]string{"127.0.0.1"})
	// }

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "message service is running",
			"status":  "running",
		})
	})

	return r

}
