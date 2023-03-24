package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once
var env = "config.json"

func InitializeConfig() {
	fmt.Println("正在使用配置：", env)

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(env)
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := v.Unmarshal(&App.Config); err != nil {
			fmt.Println(err)
		}
	})
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&App.Config); err != nil {
		fmt.Println(err)
	}

}

func OnceInitializeConfig(config string) {
	if config != "" {
		env = config
	}
	once.Do(InitializeConfig)
}
