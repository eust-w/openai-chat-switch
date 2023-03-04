package global

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var once sync.Once

func InitializeConfig() {
	// 配置所有环境
	// * 添加、修改环境配置文件请在这里进行
	envs := map[string]struct {
		_tip  string // 提示信息
		_path string // 配置文件路径
	}{
		"dev":  {"正在使用开发环境配置", "config.json"},
		"prod": {"正在使用生产环境配置", "config.prod.toml"},
	}

	// 检查环境变量
	var goEnv string
	if goEnv = os.Getenv("GO_ENV"); goEnv == "" {
		goEnv = "dev" // 默认为开发环境
	}
	env := envs[goEnv] // 取出对应环境
	fmt.Println(env._tip)

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(env._path)
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

func OnceInitializeConfig() {
	once.Do(InitializeConfig)
}
