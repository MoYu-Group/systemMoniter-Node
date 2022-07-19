package models

type Status struct {
	NodeId       string  `json:"node_id" `
	Type         string  `json:"type"`
	Location     string  `json:"location"`
	Disabled     bool    `json:"disabled"`
	Custom       string  `json:"custom"`
	Name         string  `json:"name"  binding:"required"`
	Host         string  `json:"host"  binding:"required"`
	Load1        float64 `json:"load_1"`
	Load5        float64 `json:"load_5"`
	Load15       float64 `json:"load_15"`
	IpStatus     bool    `json:"ip_status"`
	ThreadCount  uint64  `json:"thread_count"`
	ProcessCount uint64  `json:"process_count"`
	NetworkTx    uint64  `json:"network_tx"`
	NetworkRx    uint64  `json:"network_rx"`
	NetworkIn    uint64  `json:"network_in"`
	NetworkOut   uint64  `json:"network_out"`
	Ping10010    float64 `json:"ping_10010"`
	Ping10086    float64 `json:"ping_10086"`
	Ping189      float64 `json:"ping_189"`
	Time10010    uint64  `json:"time_10010"`
	Time10086    uint64  `json:"time_10086"`
	Time189      uint64  `json:"time_189"`
	TCPCount     uint64  `json:"tcp_count"`
	UDPCount     uint64  `json:"udp_count"`
	CPU          float64 `json:"cpu"`
	MemoryTotal  uint64  `json:"memory_total"`
	MemoryUsed   uint64  `json:"memory_used"`
	SwapTotal    uint64  `json:"swap_total"`
	SwapUsed     uint64  `json:"swap_used"`
	Uptime       uint64  `json:"uptime"`
	HddTotal     uint64  `json:"hdd_total"`
	HddUsed      uint64  `json:"hdd_used"`
	Online4      bool    `json:"online4"`
	Online6      bool    `json:"online6"`
}

func NewDefaultStatus() Status {
	return Status{
		NodeId:       "",
		Name:         "",
		Host:         "",
		Load1:        0.0,
		Load5:        0.0,
		Load15:       0.0,
		IpStatus:     false,
		ThreadCount:  0,
		ProcessCount: 0,
		NetworkTx:    0,
		NetworkRx:    0,
		NetworkIn:    0,
		NetworkOut:   0,
		Ping10010:    0.0,
		Ping10086:    0.0,
		Ping189:      0.0,
		Time10010:    0,
		Time10086:    0,
		Time189:      0,
		TCPCount:     0,
		UDPCount:     0,
		CPU:          0.0,
		MemoryTotal:  0,
		MemoryUsed:   0,
		SwapTotal:    0,
		SwapUsed:     0,
		Uptime:       0,
		HddTotal:     0,
		HddUsed:      0,
		Online4:      false,
		Online6:      false,
	}
}
