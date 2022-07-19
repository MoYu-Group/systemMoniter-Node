package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"systemMoniter-Node/common"
	"systemMoniter-Node/logger"
	"systemMoniter-Node/models"
	"systemMoniter-Node/settings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var errCode int
var errMsg string
var token string

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		zap.L().Info("Signal stop received,stop client program.")
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," [main] Ctrl+C pressed in Terminal,Stop client program")
		basic.Stop()
		netSpeed.Stop()
		if models.IsOpen != true {
			p10086.Stop()
			p10010.Stop()
			p189.Stop()
		}
		os.Exit(0)
	}()
}

var basic *common.Basic
var netSpeed *common.NetSpeed
var p10010 *common.Ping
var p189 *common.Ping
var p10086 *common.Ping

func main() {
	//设置东八中文时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
	//1.加载配置文件
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed err:%v\n", err)
		return
	}
	//2、初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	interval := viper.GetInt("sendInterval")
	if interval <= 0 {
		zap.L().Error("Send Status interval config error")
		return
	}
	SetupCloseHandler()
	basic = common.NewBasic()
	basic.Start()
	netSpeed = common.NewNetSpeed()
	netSpeed.Start()
	if models.IsOpen == true {
		common.PingValue.IpStatus = false
		p10086 = common.NewPing()
		p10086.RunCM()
		p10010 = common.NewPing()
		p10010.RunCU()
		p189 = common.NewPing()
		p189.RunCT()
	}
	Login()
	if errCode == 0 && errMsg == "" && token != "" {
		zap.L().Info("Sleep for waiting 5 seconds...")
		time.Sleep(time.Duration(5) * time.Second)
		SendStatus()
	} else {
		zap.L().Error("Get auth token error, exiting now")
		return
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
	StartHere:
		if (errCode == 20201 && errMsg == "ErrValidation") || (errCode == 20203 && errMsg == "ErrTokenExpired") {
			zap.L().Info("Validation token expired, need to re-login")
			Login()
			zap.L().Info("Sleep for waiting 5 seconds...")
			time.Sleep(time.Duration(5) * time.Second)
			goto StartHere
		} else if errCode != 0 && errMsg != "" {
			return
		}
		SendStatus()
	}
	zap.L().Info("Successful exiting")
}

func SendStatus() {
	//token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im5vZGUiLCJleHAiOjE2NTc4ODQ5MzgsImlhdCI6MTY1Nzg3NDEzOCwibmJmIjoxNjU3ODc0MTM4fQ.0zKRliyhLa8ZrXJ63-1UlMmRq80hI4scSoIJcMZn0ww"
	nodeData := models.NodeData{}
	err := models.SetNode(&nodeData)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	status := models.Status{}
	SetStatus(&status, &nodeData)
	body, err := json.Marshal(status)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	url := viper.GetString("url")
	saveStatusAPI := viper.GetString("api.saveStatus")
	if url == "" || saveStatusAPI == "" {
		zap.L().Error("Lost API URL config")
		return
	}
	req, err := http.NewRequest("POST", url+saveStatusAPI, bytes.NewBuffer(body))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		zap.L().Error(err2.Error())
		return
	}
	var saveStatusResponse interface{}
	err = json.Unmarshal(body, &saveStatusResponse)
	if err != nil {
		zap.L().Error(err2.Error())
		return
	}
	jsonData := saveStatusResponse.(map[string]interface{})
	errCode = int(jsonData["error"].(float64))
	errMsg = jsonData["error_msg"].(string)
	if (errCode == 20201 && errMsg == "ErrValidation") || (errCode == 20203 && errMsg == "ErrTokenExpired") {
		zap.L().Error("Token Validation Error: " + errMsg)
		token = ""
		return
	} else if errCode != 0 && errMsg != "" {
		zap.L().Error("Save node status failed, error: " + errMsg)
		return
	}
	data := jsonData["data"].(map[string]interface{})
	id := data["id"].(string)
	zap.L().Info("Save node status successful, status id: " + id)
}

func Login() {
	zap.L().Info("Start to login")
	loginUser := models.LoginUser{}
	err := models.SetLoginUser(&loginUser)
	if err != nil {
		return
	}
	body, err := json.Marshal(loginUser)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	url := viper.GetString("url")
	loginAPI := viper.GetString("api.login")
	if url == "" || loginAPI == "" {
		zap.L().Error("Lost API URL config")
		return
	}
	resp, err := http.Post(url+loginAPI, "application/json",
		bytes.NewBuffer(body))
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		zap.L().Error(err2.Error())
		return
	}
	var loginUserResponse interface{}
	err = json.Unmarshal(body, &loginUserResponse)
	if err != nil {
		zap.L().Error(err2.Error())
		return
	}
	jsonData := loginUserResponse.(map[string]interface{})
	errCode = int(jsonData["error"].(float64))
	errMsg = jsonData["error_msg"].(string)

	if errCode != 0 && errMsg != "" {
		zap.L().Error("Login failed, error: " + errMsg)
		token = ""
		return
	}
	data := jsonData["data"].(map[string]interface{})
	token = data["token"].(string)
	zap.L().Info("Login successful, token: " + token)
}

func SetStatus(status *models.Status, nodeData *models.NodeData) {
	status.Host = nodeData.Host
	status.Name = nodeData.Name
	status.Load1 = common.GetBasic.Load1
	status.Load5 = common.GetBasic.Load5
	status.Load15 = common.GetBasic.Load15
	status.ThreadCount = common.GetBasic.Thread
	status.ProcessCount = common.GetBasic.Process
	status.NetworkTx = common.GetNetSpeed.Avgtx
	status.NetworkRx = common.GetNetSpeed.Avgrx
	status.NetworkIn = uint64(common.GetNetSpeed.Nettx)
	status.NetworkOut = uint64(common.GetNetSpeed.Netrx)
	status.Ping10010 = common.PingValue.Ping10010
	status.Ping10086 = common.PingValue.Ping10086
	status.Ping189 = common.PingValue.Ping189
	status.Time10010 = common.PingValue.Time10010
	status.Time10086 = common.PingValue.Time10086
	status.Time189 = common.PingValue.Time189
	status.TCPCount = common.GetBasic.TCP
	status.UDPCount = common.GetBasic.UDP
	status.CPU = common.GetBasic.CPU
	status.MemoryTotal = common.GetBasic.MemoryTotal
	status.MemoryUsed = common.GetBasic.MemoryUsed
	status.SwapTotal = common.GetBasic.SwapTotal
	status.SwapUsed = common.GetBasic.SwapUsed
	status.Uptime = common.GetBasic.Uptime
	status.HddTotal = common.GetBasic.HddTotal
	status.HddUsed = common.GetBasic.HddUsed
	status.IpStatus = common.PingValue.IpStatus
}
