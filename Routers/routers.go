package Routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengchujin/zhilu/Middlewares"
)

func InitRouter() {

	router := gin.Default()
	router.Use(Middlewares.Cors())

	dns := router.Group("dns")
	{
		dns.GET("query", )
	}
	
	router.Run()

}
