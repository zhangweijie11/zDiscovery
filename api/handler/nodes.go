package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"github.com/zhangweijie11/zDiscovery/schemas"
	"log"
	"net/http"
)

func NodesHandler(c *gin.Context) {
	log.Println("request api/nodes...")
	var req schemas.RequestNodes
	if e := c.ShouldBindJSON(&req); e != nil {
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}

	fetchData, err := global.Discovery.Registry.Fetch(req.Env, global.DiscoveryAppId, global.NodeStatusUp, 0)
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
		"message": "",
		"data":    fetchData,
	})
}
