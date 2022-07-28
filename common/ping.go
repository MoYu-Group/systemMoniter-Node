package common

import (
	"net"
	"sync"
	"systemMoniter-Node/models"
	"time"
)

var PingValue struct {
	Ping10010   float64
	Ping10086   float64
	Ping189     float64
	Time10010   uint64
	Time10086   uint64
	Time189     uint64
	Status10010 int
	Status10086 int
	Status189   int
	IpStatus    bool
}

type Ping struct {
	stop chan struct{}
	mtx  sync.Mutex
}

func NewPing() *Ping {
	PingValue.Ping10010 = 0
	PingValue.Ping10086 = 0
	PingValue.Ping189 = 0
	PingValue.Time10010 = 0
	PingValue.Time10086 = 0
	PingValue.Time189 = 0
	PingValue.Status10010 = 0
	PingValue.Status10086 = 0
	PingValue.Status189 = 0
	PingValue.IpStatus = false
	return &Ping{
		stop: make(chan struct{}),
	}
}

func (ping *Ping) Stop() {
	if ping != nil {
		close(ping.stop)
	}
}

func (ping *Ping) RunCU() {
	go func() {
		t1 := time.Duration(models.Interval) * time.Second
		t := time.NewTicker(t1)

		var lostPacket = 0
		var allPacket = 0
		var lostConnect = false
		var status = 0
		var lostRate = float64(0.0)
		var pingValue = uint64(0)
		startTime := time.Now()
		defaulttimeout := 5 * time.Second
		for {
			select {
			case <-ping.stop:
				t.Stop()
				return
			case <-t.C:
				ping.mtx.Lock()
				t := time.Now()
				url := models.Cu
				conn, err := net.DialTimeout("tcp", url, defaulttimeout)
				if err != nil {
					//zap.L().Error("Error try to connect China Unicom :", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," [ping]Error try to connect China unicom :", err)
					lostConnect = true
					lostPacket += 1
				}
				tcpconn, ok := conn.(*net.TCPConn)
				if ok {
					tcpconn.SetLinger(0)
				}
				if conn != nil {
					conn.Close()
				}
				diffTime := time.Since(t)
				//TODO:三网延迟和丢包率算法存在问题
				//fmt.Println(diffTime)
				allPacket += 1
				if allPacket > 100 {
					lostRate = float64(lostPacket/allPacket) * 100
				}
				//fmt.Println("ALL     LOST    RATE")
				//fmt.Printf("%10d  %10d %10f\n",allPacket,lostPacket,pingValue.lostRate)
				if lostConnect {
					pingValue = 0
					status = 1
				} else {
					pingValue = uint64(diffTime / time.Millisecond)
					status = 0
				}
				lostConnect = false
				resetTime := uint64(time.Since(startTime) / time.Second)
				if resetTime > 3600 {
					lostPacket = 0
					allPacket = 0
					startTime = time.Now()
				}
				PingValue.Time10010 = pingValue
				PingValue.Ping10010 = lostRate
				PingValue.Status10010 = status
				if PingValue.Status10010+PingValue.Status10086+PingValue.Status189 >= 2 {
					PingValue.IpStatus = false
				} else {
					PingValue.IpStatus = true
				}
				ping.mtx.Unlock()
			}
		}
	}()
}

func (ping *Ping) RunCT() {
	go func() {
		t1 := time.Duration(models.Interval) * time.Second
		t := time.NewTicker(t1)
		var lostPacket = 0
		var allPacket = 0
		var lostConnect = false
		var status = 0
		var lostRate = 0.0
		var pingValue = uint64(0)
		startTime := time.Now()
		defaulttimeout := 5 * time.Second
		for {
			select {
			case <-ping.stop:
				t.Stop()
				return
			case <-t.C:
				ping.mtx.Lock()
				t := time.Now()
				url := models.Ct
				conn, err := net.DialTimeout("tcp", url, defaulttimeout)
				if err != nil {
					//zap.L().Error("Error try to connect China telecom :", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," [ping]Error try to connect China unicom :", err)
					lostConnect = true
					lostPacket += 1
				}
				tcpconn, ok := conn.(*net.TCPConn)
				if ok {
					tcpconn.SetLinger(0)
				}
				if conn != nil {
					conn.Close()
				}
				diffTime := time.Since(t)
				//TODO:三网延迟和丢包率算法存在问题
				//fmt.Println(diffTime)
				allPacket += 1
				if allPacket > 100 {
					lostRate = float64(lostPacket/allPacket) * 100
				}
				//fmt.Println("ALL     LOST    RATE")
				//fmt.Printf("%10d  %10d %10f\n",allPacket,lostPacket,pingValue.lostRate)
				if lostConnect {
					pingValue = 0
					status = 1
				} else {
					pingValue = uint64(diffTime / time.Millisecond)
					status = 0
				}
				lostConnect = false
				resetTime := uint64(time.Since(startTime) / time.Second)
				if resetTime > 3600 {
					lostPacket = 0
					allPacket = 0
					startTime = time.Now()
				}
				PingValue.Time189 = pingValue
				PingValue.Ping189 = lostRate
				PingValue.Status189 = status
				if PingValue.Status10010+PingValue.Status10086+PingValue.Status189 >= 2 {
					PingValue.IpStatus = false
				} else {
					PingValue.IpStatus = true
				}
				ping.mtx.Unlock()
			}
		}
	}()
}

func (ping *Ping) RunCM() {
	go func() {
		t1 := time.Duration(models.Interval) * time.Second
		t := time.NewTicker(t1)
		var lostPacket = 0
		var allPacket = 0
		var lostConnect = false
		var status = 0
		var lostRate = 0.0
		var pingValue = uint64(0)
		startTime := time.Now()
		defaulttimeout := 5 * time.Second
		for {
			select {
			case <-ping.stop:
				t.Stop()
				return
			case <-t.C:
				ping.mtx.Lock()
				t := time.Now()
				url := models.Cm
				conn, err := net.DialTimeout("tcp", url, defaulttimeout)
				if err != nil {
					//zap.L().Error("Error try to connect China mobile :", zap.Error(err))
					//fmt.Println(time.Now().Format("2006-01-02 15:04:05")," [ping]Error try to connect China unicom :", err)
					lostConnect = true
					lostPacket += 1
				}
				tcpconn, ok := conn.(*net.TCPConn)
				if ok {
					tcpconn.SetLinger(0)
				}
				if conn != nil {
					conn.Close()
				}
				diffTime := time.Since(t)
				//TODO:三网延迟和丢包率算法存在问题
				//fmt.Println(diffTime)
				allPacket += 1
				if allPacket > 100 {
					lostRate = float64(lostPacket/allPacket) * 100
				}
				//fmt.Println("ALL     LOST    RATE")
				//fmt.Printf("%10d  %10d %10f\n",allPacket,lostPacket,pingValue.lostRate)
				if lostConnect {
					pingValue = 0
					status = 1
				} else {
					pingValue = uint64(diffTime / time.Millisecond)
					status = 0
				}
				lostConnect = false
				resetTime := uint64(time.Since(startTime) / time.Second)
				if resetTime > 3600 {
					lostPacket = 0
					allPacket = 0
					startTime = time.Now()
				}
				PingValue.Time10086 = pingValue
				PingValue.Ping10086 = lostRate
				PingValue.Status10086 = status
				if PingValue.Status10010+PingValue.Status10086+PingValue.Status189 >= 2 {
					PingValue.IpStatus = false
				} else {
					PingValue.IpStatus = true
				}
				ping.mtx.Unlock()
			}
		}
	}()
}
