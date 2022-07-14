package common

import (
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	nnet "github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
)

var GetBasic struct {
	MemoryTotal uint64
	MemoryUsed  uint64
	CPU         float64
	Uptime      uint64
	SwapTotal   uint64
	SwapUsed    uint64
	Load1       float64
	Load5       float64
	Load15      float64
	NetworkIn   uint64
	NetworkOut  uint64
	HddUsed     uint64
	HddTotal    uint64
	TCP         uint64
	UDP         uint64
	Process     uint64
	Thread      uint64
}

type Basic struct {
	stop chan struct{}
	mtx  sync.Mutex
}

func NewBasic() *Basic {
	GetBasic.MemoryTotal = 0
	GetBasic.MemoryUsed = 0
	GetBasic.CPU = 0.0
	GetBasic.Uptime = 0
	GetBasic.SwapTotal = 0
	GetBasic.SwapUsed = 0
	GetBasic.Load1 = 0.0
	GetBasic.Load5 = 0.0
	GetBasic.Load15 = 0.0
	GetBasic.NetworkIn = 0
	GetBasic.NetworkOut = 0
	GetBasic.HddUsed = 0
	GetBasic.HddTotal = 0
	GetBasic.TCP = 0
	GetBasic.UDP = 0
	GetBasic.Process = 0
	GetBasic.Thread = 0
	return &Basic{
		stop: make(chan struct{}),
	}
}

func (basic *Basic) Stop() {
	close(basic.stop)
}

func (basic *Basic) Start() {
	go func() {
		t1 := time.Duration(1) * time.Second
		t := time.NewTicker(t1)
		for {
			select {
			case <-basic.stop:
				t.Stop()
				return
			case <-t.C:
				basic.mtx.Lock()
				memInfo, err := mem.VirtualMemory()
				if err != nil {
					//logger.Errorf("Get memory usage error:", err)
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get memory usage error:",err)
					GetBasic.MemoryTotal = 0
					GetBasic.MemoryUsed = 0
				} else {
					GetBasic.MemoryTotal = memInfo.Total / 1024 // 需要转单位
					GetBasic.MemoryUsed = memInfo.Used / 1024   // 需要转单位
				}

				totalPercent, err := cpu.Percent(time.Second, false)
				if err != nil {
					zap.L().Error("Get cpu usage error:", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," [GetInfo]Get cpu usage error:",err)
					GetBasic.CPU = 0.0
				} else {
					if totalPercent != nil {
						GetBasic.CPU = totalPercent[0]
					} else {
						zap.L().Error("Get cpu usage error:", zap.Error(err))
						//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get cpu usage error:",err)
					}
				}
				hInfo, err := host.Info()
				if err != nil {
					zap.L().Error("Get Uptime error:", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," get Uptime error",err)
					GetBasic.Uptime = 0
				} else {
					GetBasic.Uptime = hInfo.Uptime
				}
				//swap 没有造好的轮子，自己加的
				swapMemory, err := mem.SwapMemory()
				if err != nil {
					zap.L().Error("Get swap memory error:", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get swap memory error:",err)
					GetBasic.SwapTotal = 0
					GetBasic.SwapUsed = 0
				} else {
					GetBasic.SwapTotal = swapMemory.Total / 1024 // 需要转单位
					GetBasic.SwapUsed = swapMemory.Used / 1024   // 需要转单位
				}
				getLoad()
				trafficCount()
				spaceCount()
				tupd()
				basic.mtx.Unlock()
			}
		}
	}()
}

func trafficCount() {
	netInfo, err := nnet.IOCounters(true)
	if err != nil {
		zap.L().Error("Getting traffic count error:", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Getting traffic count error:",err)
	}
	var bytesSent uint64 = 0
	var bytesRecv uint64 = 0
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
	GetBasic.NetworkIn = bytesRecv
	GetBasic.NetworkOut = bytesSent
}

func spaceCount() {
	// golang 没有类似于在 python 的 dict 或 tuple 的 in 查找关键字，自己写多重判断实现
	diskList, _ := disk.Partitions(false)
	var total uint64 = 0
	var used uint64 = 0
	for _, d := range diskList {
		fsType := strings.ToLower(d.Fstype)
		//fmt.Println(d.Fstype)
		if !strings.Contains(fsType, "ext4") &&
			!strings.Contains(fsType, "ext2") &&
			!strings.Contains(fsType, "reiserfs") &&
			!strings.Contains(fsType, "jfs") &&
			!strings.Contains(fsType, "btrfs") &&
			!strings.Contains(fsType, "fuseblk") &&
			!strings.Contains(fsType, "zfs") &&
			!strings.Contains(fsType, "simfs") &&
			!strings.Contains(fsType, "ntfs") &&
			!strings.Contains(fsType, "fat32") &&
			!strings.Contains(fsType, "exfat") &&
			!strings.Contains(fsType, "xfs") {
		} else {
			if strings.Contains(d.Device, "Z:") { //特殊盘符自己写处理
				continue
			} else {
				diskUsageOf, _ := disk.Usage(d.Mountpoint)
				path := diskUsageOf.Path
				//不统计K8s的虚拟挂载点，see here：https://github.com/shirou/gopsutil/issues/1007
				if !strings.Contains(path, "/var/lib/kubelet") {
					used += diskUsageOf.Used
					total += diskUsageOf.Total
				}
			}
		}
	}
	GetBasic.HddUsed = used / 1024.0 / 1024.0
	GetBasic.HddTotal = total / 1024.0 / 1024.0
}

func getLoad() {
	// linux or freebsd only
	hInfo, err := host.Info()
	if err != nil {
		zap.L().Error("Get load info error", zap.Error(err))
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," get load info error",err)
		GetBasic.Load1 = 0.0
		GetBasic.Load5 = 0.0
		GetBasic.Load15 = 0.0
	} else {
		if hInfo.OS == "linux" || hInfo.OS == "freebsd" {
			l, err := load.Avg()
			if err != nil {
				zap.L().Error("Get CPU loads failed:", zap.Error(err))
				//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," Get CPU loads failed:",err)
				GetBasic.Load1 = 0.0
				GetBasic.Load5 = 0.0
				GetBasic.Load15 = 0.0
			} else {
				GetBasic.Load1 = l.Load1
				GetBasic.Load5 = l.Load5
				GetBasic.Load15 = l.Load15
			}
		} else {
			GetBasic.Load1 = 0.0
			GetBasic.Load5 = 0.0
			GetBasic.Load15 = 0.0
		}
	}
}
