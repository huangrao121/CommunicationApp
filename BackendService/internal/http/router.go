package http

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(SecureHeaders())
	router.Use(SetLogger())

	if gin.Mode() == gin.ReleaseMode {
		router.Use(CSRF())
	}

	// router.GET("/health", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "OK"})
	// })

	return router
}
