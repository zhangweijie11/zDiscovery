package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangweijie11/zDiscovery/common"
	"log"
	"net/http"
)

func FetchAllHandler(c *gin.Context) {
	log.Println("request api/fetchall...")

	data := common.Discovery.Registry.FetchAll()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "",
		"data":    data,
	})
}
