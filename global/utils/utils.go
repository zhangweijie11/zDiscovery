package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func HttpPost(url string, data interface{}) (string, error) {
	client := &http.Client{}
	js, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewReader(js))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Println("公共请求建立连接时出现异常", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("公共请求发起请求时出现异常", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("公共请求读取响应时出现异常", err)
		return "", err
	}
	return string(body), nil
}
