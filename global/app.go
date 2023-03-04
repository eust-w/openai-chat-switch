package global

import (
	"github.com/eust-w/openai-chat-switch/database"
	"go.uber.org/zap"
)

type Application struct {
	Config Configuration      // 项目配置
	Log    *zap.SugaredLogger // 日志系统
	Db     *database.ChatDb
}

var App = new(Application)
