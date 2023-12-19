package registry

import (
	"encoding/json"
	"fmt"
	"github.com/zhangweijie11/zDiscovery/config"
	"github.com/zhangweijie11/zDiscovery/global"
	"github.com/zhangweijie11/zDiscovery/global/utils"
	"log"
	"strconv"
	"time"
)

type Node struct {
	config      *config.Config // 节点配置
	addr        string         // 节点地址
	status      int            // 节点状态
	registerUrl string         // 注册地址
	cancelUrl   string         // 注销地址
	renewUrl    string         // 续约地址
}

func NewNode(config *config.Config, addr string) *Node {
	return &Node{
		config:      config,
		addr:        addr,
		status:      global.NodeStatusDown,
		registerUrl: fmt.Sprintf("http://%s%s", addr, global.RegisterURL),
		cancelUrl:   fmt.Sprintf("http://%s%s", addr, global.CancelURL),
		renewUrl:    fmt.Sprintf("http://%s%s", addr, global.RenewURL),
	}
}

func (node *Node) call(url string, action global.Action, instance *Instance, data interface{}) error {
	params := make(map[string]interface{})
	params["env"] = instance.Env
	params["appid"] = instance.AppID
	params["hostname"] = instance.Hostname
	params["replication"] = true //broadcast stop here
	switch action {
	case global.Register:
		params["addresses"] = instance.Addresses
		params["status"] = instance.Status
		params["version"] = instance.Version
		params["reg_timestamp"] = strconv.FormatInt(instance.RegTimestamp, 10)
		params["dirty_timestamp"] = strconv.FormatInt(instance.DirtyTimestamp, 10)
		params["latest_timestamp"] = strconv.FormatInt(instance.LatestTimestamp, 10)
	case global.Renew:
		params["dirty_timestamp"] = strconv.FormatInt(instance.DirtyTimestamp, 10)
		params["renew_timestamp"] = time.Now().UnixNano()
	case global.Cancel:
		params["latest_timestamp"] = strconv.FormatInt(instance.LatestTimestamp, 10)
	}
	// 请求其他节点
	resp, err := utils.HttpPost(url, params)
	if err != nil {
		log.Println(err)
		return err
	}
	response := Response{}
	err = json.Unmarshal([]byte(resp), &response)
	if err != nil {
		log.Println(err)
		return err
	}
	if response.Code != global.StatusOK { //code!=200
		log.Printf("uri is (%v),response code (%v)\n", url, response.Code)
		json.Unmarshal(response.Data, data)
		return utils.Conflict
	}
	return nil
}

func (node *Node) Register(instance *Instance) error {
	return node.call(node.registerUrl, global.Register, instance, nil)
}

func (node *Node) Cancel(instance *Instance) error {
	return node.call(node.cancelUrl, global.Cancel, instance, nil)
}

func (node *Node) Renew(instance *Instance) error {
	var res *Instance
	err := node.call(node.renewUrl, global.Renew, instance, &res)
	// 如果节点续约出现问题，则直接判定为节点下线
	if err == utils.ServerError {
		log.Printf("node call %s ! renew error %s \n", node.renewUrl, err)
		node.status = global.NodeStatusDown //node down
		return err
	}
	// 如果显示节点不存在，则注册节点
	if err == utils.NotFound { //register
		log.Printf("node call %s ! renew not found, register again \n", node.renewUrl)
		return node.call(node.registerUrl, global.Register, instance, nil)
	}
	// 如果网络冲突并且实例不为空则注册节点
	if err == utils.Conflict && res != nil {
		return node.call(node.registerUrl, global.Register, res, nil)
	}
	return err
}
