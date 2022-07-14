package common

import (
	"strings"
	"sync"
	"time"

	nnet "github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
)

var GetNetSpeed struct {
	Netrx float64
	Nettx float64
	Clock float64
	Diff  float64
	Avgrx uint64
	Avgtx uint64
}

type NetSpeed struct {
	stop chan struct{}
	mtx  sync.Mutex
}

func NewNetSpeed() *NetSpeed {
	GetNetSpeed.Avgrx = 0
	GetNetSpeed.Avgtx = 0
	GetNetSpeed.Netrx = 0
	GetNetSpeed.Nettx = 0
	GetNetSpeed.Clock = 0
	GetNetSpeed.Diff = 0
	return &NetSpeed{
		stop: make(chan struct{}),
	}
}

func (netSpeed *NetSpeed) Start() {
	go func() {
		t1 := time.Duration(1) * time.Second
		t := time.NewTicker(t1)
		for {
			select {
			case <-netSpeed.stop:
				t.Stop()
				return
			case <-t.C:
				netSpeed.mtx.Lock()
				var bytesSent uint64 = 0
				var bytesRecv uint64 = 0
				netInfo, err := nnet.IOCounters(true)
				if err != nil {
					zap.L().Error("Get network speed error: ", zap.Error(err))
					//logger.Errorf("Get network speed error:", err)
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get network speed error:",err)
				}
				for _, v := range netInfo {
					if strings.Contains(v.Name, "lo") ||
						strings.Contains(v.Name, "tun") ||
						strings.Contains(v.Name, "docker") ||
						strings.Contains(v.Name, "veth") ||
						strings.Contains(v.Name, "br-") ||
						strings.Contains(v.Name, "vmbr") ||
						strings.Contains(v.Name, "vnet") ||
						strings.Contains(v.Name, "kube") {
						continue
					}
					bytesSent += v.BytesSent
					bytesRecv += v.BytesRecv
				}
				timeUnix := float64(time.Now().Unix())
				GetNetSpeed.Diff = timeUnix - GetNetSpeed.Clock
				GetNetSpeed.Clock = timeUnix
				GetNetSpeed.Netrx = float64(bytesRecv-GetNetSpeed.Avgrx) / GetNetSpeed.Diff
				GetNetSpeed.Nettx = float64(bytesSent-GetNetSpeed.Avgtx) / GetNetSpeed.Diff
				GetNetSpeed.Avgtx = bytesSent
				GetNetSpeed.Avgrx = bytesRecv
				netSpeed.mtx.Unlock()
			}
		}
	}()
}

func (netSpeed *NetSpeed) Stop() {
	close(netSpeed.stop)
}
