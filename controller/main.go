package controller

import (
	"github.com/eust-w/openai-chat-switch/gpt"
	"strings"
)

var AnserFuncSlice map[string]func(a, b string) string = map[string]func(a, b string) string{}

func Answer(prompt string, userId string) (out string) {
	var ok bool
	if out, ok = checkUser(userId); ok {
		return answer(prompt, userId)
	}
	return out
}

//检查用户是否拥有权限，返回结果和权限
func checkUser(userId string) (out string, permissions bool) {
	return "", true
}

// 返回结果，会在控制层对命令进行过滤
func answer(prompt string, userId string) (out string) {
	prompt = prunePrompt(prompt)
	if f, ok := AnserFuncSlice[prompt]; ok {
		return f(prompt, userId)
	}
	return answerByGpt(prompt, userId)
}

func prunePrompt(prompt string) string {
	return strings.TrimSpace(prompt)
}

// 通过gpt回复结果
func answerByGpt(prompt string, userId string) (out string) {
	prompt = processPrompt(prompt)
	return gpt.Answer(prompt, userId)
}

// 用来对prompt进行处理，对新的prompt进行加工，例如生成模板什么的
func processPrompt(prompt string) (processedPrompt string) {
	//
	prompt = strings.TrimSpace(prompt)
	return prompt
}
