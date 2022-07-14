package main

import (
	"fmt"
	"systemMoniter-Node/common"
	"systemMoniter-Node/logger"
	"systemMoniter-Node/models"
	"systemMoniter-Node/settings"
	"time"

	"go.uber.org/zap"
)

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
	Start()
	zap.L().Info("Successful exiting")
}

func Start() {
	basic := common.NewBasic()
	basic.Start()
	defer basic.Stop()
	netSpeed := common.NewNetSpeed()
	netSpeed.Start()
	defer netSpeed.Stop()
	if models.IsOpen == true {
		common.PingValue.IpStatus = false
		p10086 := common.NewPing()
		defer p10086.Stop()
		p10086.RunCM()
		p10010 := common.NewPing()
		defer p10010.Stop()
		p10010.RunCU()
		p189 := common.NewPing()
		defer p189.Stop()
		p189.RunCT()
	}
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
	status.CPUCount = common.GetBasic.CPU
	status.MemoryTotal = common.GetBasic.MemoryTotal
	status.MemoryUsed = common.GetBasic.MemoryUsed
	status.SwapTotal = common.GetBasic.SwapTotal
	status.SwapUsed = common.GetBasic.SwapUsed
	status.Uptime = common.GetBasic.Uptime
	status.HddTotal = common.GetBasic.HddTotal
	status.HddUsed = common.GetBasic.HddUsed
	status.IpStatus = common.PingValue.IpStatus
}
