package Controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/pengchujin/zhilu/Services"
)

func QueryDns(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	log.Println(name)

	dnsType := c.DefaultQuery("type", "1")
	log.Println(dnsType)

	subNetAddr := c.DefaultQuery("edns_client_subnet", c.ClientIP())

	log.Println("Client IP: ", subNetAddr)

	a := Services.GetDnsRes(name, dnsType, subNetAddr)
	log.Println(a)

	c.JSON(http.StatusOK, a)
}
