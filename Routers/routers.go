package Routers

import (
	"github.com/gin-gonic/gin"
	"github.com/pengchujin/zhilu/Middlewares"
	"github.com/pengchujin/zhilu/Controllers"
)

func InitRouter() {

	router := gin.Default()
	router.Use(Middlewares.Cors())

	dns := router.Group("dns", Controllers.QueryDns)
	{
		dns.GET("query", )
	}
	
	router.Run()

}
