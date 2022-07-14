package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigFile("config.yaml") // 指定配置文件
	viper.AddConfigPath(".")           // 指定查找配置文件的路径
	err := viper.ReadInConfig()        // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig() Failed err:%v\n", err)
		return err
	}
	//监控并重新读取配置文件
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})
	return nil
}
