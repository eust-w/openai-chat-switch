package openai_chat_switch

import (
	"fmt"
	"github.com/eust-w/openai-chat-switch/global"
	"github.com/sashabaranov/go-openai"
	"strconv"
	"strings"
	"time"
)

func init() {
	ExtensionExactMatchReturnOutcome["重置"] = reset
	ExtensionExactMatchReturnOutcome["报时"] = nowTime
	ExtensionPrefixMatchReturnOutcome["模型切换到"] = changeModel
	ExtensionPrefixMatchReturnOutcome["Temperature切换到"] = changeTemperature
}

func reset(prompt, userId string) string {
	global.App.Db.DeleteContext(userId)
	return "已经帮您清除上下文"
}

func nowTime(prompt, userId string) string {
	return "当前时间为:" + time.Now().Format("2006-01-02 15:04:05")
}

func changeModel(prompt, userId string) string {
	errMessage := "不支持您的模型，支持的模型有"
	m := strings.Split(prompt, "模型切换到")
	if len(m) != 2 {
		return errMessage
	}
	model := m[1]
	realModel, err := checkModel(model)
	if err != nil {
		return errMessage
	}
	global.App.Db.SetModel(userId, realModel)
	return "已成功帮您将模型切换到" + realModel
}

func checkModel(model string) (string, error) {
	gpt3Dot5TurboWhitelist := map[string]struct{}{"chatgpt": {}, "gpt3": {}, "gpt3.5": {}, "gpt-3.5": {}, "3.5": {}, "gpt35": {}, "gpt-3.5-turbo": {}, "gpt-3.5turbo": {}, "gpt3.5turbo": {}, "gpt3.5-turbo": {}}
	gpt3Dot5Turbo0301Whitelist := map[string]struct{}{"gpt-3.5-turbo-0301": {}, "gpt3.5-turbo-0301": {}, "gpt3.5turbo-0301": {}, "gpt3.5turbo0301": {}, "gpt-3.5-turbo0301": {}, "gpt-3.5turbo0301": {}, "gpt-3.5turbo-0301": {}}
	gpt4Whitelist := map[string]struct{}{"gpt4": {}, "gpt-4": {}, "chatgpt4": {}}
	gpt40314Whitelist := map[string]struct{}{"gpt-4-0314": {}, "gpt4-0314": {}, "gpt40314": {}, "gpt-40314": {}}
	if _, ok := gpt3Dot5TurboWhitelist[model]; ok {
		return openai.GPT3Dot5Turbo, nil
	}
	if _, ok := gpt3Dot5Turbo0301Whitelist[model]; ok {
		return openai.GPT3Dot5Turbo0301, nil
	}
	if _, ok := gpt4Whitelist[model]; ok {
		return openai.GPT4, nil
	}
	if _, ok := gpt40314Whitelist[model]; ok {
		return openai.GPT40314, nil
	}
	return "", fmt.Errorf("unsupported model")
}

func changeTemperature(prompt, userId string) string {
	errMessage := "不支持的性格设置"
	m := strings.Split(prompt, "Temperature切换到")
	if len(m) != 2 {
		return errMessage
	}
	temperature, err := strconv.ParseFloat(m[1], 32)
	if err != nil {
		return errMessage
	}
	if 0 < temperature && temperature < 1 {
		global.App.Db.SetTemperature(userId, float32(temperature))
	}
	return errMessage
}
