package apiUser

import (
	"CTFgo/logs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	HEX_BOT_URL     = "http://xxxxx"
	CONTENT_TYPE    = "application/json"
	SERVER_CHAN_URL = "https://sctapi.ftqq.com/<secret_key>.send?title=%s&desp=%s"
)

// HexBotRequest 定义对hex酱请求后的返回结构体。
type HexBotRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HexBotMsgRequest 定义对hex酱请求后的返回结构体。
type HexBotMsgRequest struct {
	Msg   string `json:"msg"`
	QQNum int    `json:"num"`
}

func PostHexBot(data *HexBotMsgRequest) error {
	var title = "Hex Bot Error"
	client := &http.Client{Timeout: 1 * time.Second} // timeout: 1s
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(HEX_BOT_URL, CONTENT_TYPE, bytes.NewBuffer(jsonStr))
	if err != nil {
		logs.WARNING("Post msg to Hex Bot error:", err)
		go GetServerChan(title, data.Msg)
		return err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	var request HexBotRequest
	if err = json.Unmarshal(result, &request); err != nil {
		logs.WARNING("Hex Bot msg unmarshal error:", err)
		go GetServerChan(title, data.Msg)
		return err
	}
	if request.Code != 200 {
		logs.WARNING("Hex Bot msg code error:", err)
		go GetServerChan(title, data.Msg)
		return err
	}
	return nil
}

func GetServerChan(title, desp string) error {
	url := fmt.Sprintf(SERVER_CHAN_URL, title, desp)
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logs.WARNING("Get msg to ServerChan error:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
