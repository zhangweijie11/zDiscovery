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

func RegisterHandler(c *gin.Context) {
	log.Println("request api/register...")
	var req schemas.RequestRegister
	if e := c.ShouldBindJSON(&req); e != nil {
		log.Println("error:", e)
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//bind instance
	instance := registry.NewInstance(&req)
	if instance.Status != global.StatusReceive && instance.Status != global.StatusNotReceive {
		log.Println("register params status invalid")
		err := utils.ParamError
		c.JSON(http.StatusOK, gin.H{
			"code":    err.Code(),
			"message": err.Error(),
		})
		return
	}
	//dirtytime
	if req.DirtyTimestamp > 0 {
		instance.DirtyTimestamp = req.DirtyTimestamp
	}
	global.Discovery.Registry.Register(instance, req.LatestTimestamp)
	//default do replicate. if request come from other server, req.Replication is true, ignore replicate.
	if !req.Replication {
		global.Discovery.Nodes.Load().(*registry.Nodes).Replicate(global.Register, instance)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    global.StatusOK,
		"message": "",
		"data":    "",
	})
}
