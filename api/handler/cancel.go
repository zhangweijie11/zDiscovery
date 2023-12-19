package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"github.com/zhangweijie11/zDiscovery/schemas"
	"github.com/zhangweijie11/zDiscovery/services/registry"
	"log"
	"net/http"
)

func CancelHandler(c *gin.Context) {
	log.Println("request api/cancel...")
	var req schemas.RequestCancel
	if e := c.ShouldBindJSON(&req); e != nil {
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	instance, err := global.Discovery.Registry.Cancel(req.Env, req.AppId, req.Hostname, req.LatestTimestamp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//replication to other server
	if !req.Replication {
		global.Discovery.Nodes.Load().(*registry.Nodes).Replicate(global.Cancel, instance)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.StatusOK,
		"message": "",
	})
}
