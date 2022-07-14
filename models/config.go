package models

import (
	"github.com/spf13/viper"
)

var (
	Interval  = 1
	Cu        = "120.52.99.224"
	Ct        = "183.78.182.66"
	Cm        = "211.139.145.129"
	Porbeport = 80
	IsOpen    = true
)

type Config struct {
	interval  int
	cu        string
	ct        string
	cm        string
	isOpen    bool
	porbeport int
}

// func NewConfig() Config {
// 	return Config{
// 		interval:  Interval,
// 		cu:        Cu,
// 		ct:        Ct,
// 		cm:        Cm,
// 		isOpen:    IsOpen,
// 		porbeport: Porbeport,
// 	}
// }

func SetConfig() {
	Interval = viper.GetInt("network.interval")
	IsOpen = false
	if Interval > 0 {
		IsOpen = true
	} else {
		IsOpen = false
	}
	Cu = viper.GetString("network.cu") + ":" + viper.GetString("network.port")
	Ct = viper.GetString("network.ct") + ":" + viper.GetString("network.port")
	Cm = viper.GetString("network.cm") + ":" + viper.GetString("network.port")

}
