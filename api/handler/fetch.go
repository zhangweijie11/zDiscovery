package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangweijie11/zDiscovery/common"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"github.com/zhangweijie11/zDiscovery/schemas"
	"log"
	"net/http"
)

func FetchHandler(c *gin.Context) {
	log.Println("request api/fetch...")
	var req schemas.RequestFetch
	if e := c.ShouldBindJSON(&req); e != nil {
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}

	// 同步
	fetchData, err := common.Discovery.Registry.Fetch(req.Env, req.AppId, req.Status, 0)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"data":    "",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    fetchData,
		"message": "",
	})
}
