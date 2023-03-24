package openai_chat_switch

import (
	"github.com/eust-w/openai-chat-switch/global"
	"time"
)

func init() {
	AnserFuncSlice["重置"] = reset
	AnserFuncSlice["报时"] = nowTime
}

func reset(prompt, userId string) string {
	global.App.Db.DeleteContext(userId)
	return "已经帮您清除上下文"
}

func nowTime(prompt, userId string) string {
	return "当前时间为:" + time.Now().Format("2006-01-02 15:04:05")
}
