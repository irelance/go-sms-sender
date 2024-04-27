package go_sms_sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ForceSMSClient struct {
	AppId string
	AppSecret string
	DeviceId string
	BaseUrl string
	template string
}

func (a *ForceSMSClient) sign(message []byte, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write(message)
	h.Write(key)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (a *ForceSMSClient) IsReceiveMessage(param map[string]string) (bool, error) {
	postBody, _ := json.Marshal(map[string]string{
		"appId": a.AppId,
		"deviceId": a.DeviceId,
		"from":  param["from"],
		"msg":  param["msg"],
		"ttl":  param["ttl"],
	})
	fmt.Println("postBody",string(postBody),a.AppSecret)

	// 将JSON数据转换为字节序列
	requestBody := bytes.NewBuffer(postBody)

	url:=a.BaseUrl+"/api/v1/gateway/received"
	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return false, err
	}

	// 设置请求头信息，指定内容类型为JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sign", a.sign(postBody,a.AppSecret))

	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return string(body)=="true" ,nil
}

func GetForceSMSClient(accessKeyID string, secretAccessKey string, template string, other []string) (*ForceSMSClient, error) {
	smsClient := &ForceSMSClient{
		AppId:accessKeyID,
		AppSecret :secretAccessKey,
		DeviceId :other[1],
		BaseUrl :other[0],
		template: template,
	}
	return smsClient, nil
}

func (a *ForceSMSClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	return fmt.Errorf("implement me")
}
