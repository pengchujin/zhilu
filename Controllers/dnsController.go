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

	a := Services.GetDnsRes("1")
	log.Println(a)

	c.JSON(http.StatusOK, a)
}
