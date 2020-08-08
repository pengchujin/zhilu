package Controllers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func queryDns(c *gin.Context) {

	name := c.DefaultQuery("name", "")
	log.Println(name)

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"message": "success",
		"data": 1,
	})
}
