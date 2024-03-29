package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhangweijie11/zDiscovery/common"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"github.com/zhangweijie11/zDiscovery/schemas"
	"github.com/zhangweijie11/zDiscovery/services/registry"
	"log"
	"net/http"
)

func RenewHandler(c *gin.Context) {
	log.Println("request api/renew...")
	var req schemas.RequestRenew
	if e := c.ShouldBindJSON(&req); e != nil {
		log.Println("续约参数错误error:", e)
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}

	//registry global  discovery
	instance, err := common.Discovery.Registry.Renew(req.Env, req.AppId, req.Hostname)
	if err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}

	// 复制到其他节点
	if !req.Replication {
		common.Discovery.Nodes.Load().(*registry.Nodes).Replicate(global.Renew, instance)
	}

	//过期
	if req.DirtyTimestamp > instance.DirtyTimestamp {
		err = utils.NotFound
	} else if req.DirtyTimestamp < instance.DirtyTimestamp { //冲突
		err = utils.Conflict
	}
	c.JSON(http.StatusOK, gin.H{
		"code": global.StatusOK,
	})
}
