package openai_chat_switch

import (
	"fmt"
	"github.com/eust-w/openai-chat-switch/database"
	"github.com/eust-w/openai-chat-switch/global"
	"github.com/eust-w/openai-chat-switch/gpt"
	"strings"
	"sync"
)

type RwMap struct {
	globalMap map[string]struct{}
	sync.RWMutex
}

func (r *RwMap) Get(name string) bool {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.globalMap[name]
	return ok
}

func (r *RwMap) Set(name string) bool {
	if r.Get(name) {
		return false
	}
	r.Lock()
	defer r.Unlock()
	r.globalMap[name] = struct{}{}
	return true
}

func (r *RwMap) Del(name string) bool {
	if !r.Get(name) {
		return false
	}
	r.Lock()
	defer r.Unlock()
	delete(r.globalMap, name)
	return true
}

var globalRwMap = RwMap{globalMap: map[string]struct{}{}}

type extension struct {
	funcMap map[string]func(a, b string) string
}

func (e extension) GetFuncs() map[string]func(a, b string) string {
	return e.funcMap
}

func (e extension) GetFunc(name string) (func(a, b string) string, bool) {
	f, ok := e.funcMap[name]
	return f, ok
}

func (e extension) AddFunc(name string, function func(a, b string) string) error {
	if globalRwMap.Set(name) {
		e.funcMap[name] = function
		return nil
	}
	return fmt.Errorf("add func error")
}

func (e extension) DelFunc(name string) {
	if !globalRwMap.Del(name) {
		global.App.Log.Warnf("del global func name error")
	}
	delete(e.funcMap, name)
}

//var ExtensionExactMatchReturnOutcome map[string]func(a, b string) string = map[string]func(a, b string) string{}
//var ExtensionPrefixMatchReturnOutcome map[string]func(a, b string) string = map[string]func(a, b string) string{}
//var ExtensionExactMatchReturnPrompt map[string]func(a, b string) string = map[string]func(a, b string) string{}
//var ExtensionPrefixMatchReturnPrompt map[string]func(a, b string) string = map[string]func(a, b string) string{}

var ExtensionExactMatchReturnOutcome extension = extension{funcMap: map[string]func(a string, b string) string{}}
var ExtensionPrefixMatchReturnOutcome = extension{funcMap: map[string]func(a string, b string) string{}}
var ExtensionExactMatchReturnPrompt = extension{funcMap: map[string]func(a string, b string) string{}}
var ExtensionPrefixMatchReturnPrompt = extension{funcMap: map[string]func(a string, b string) string{}}

func NewGlobal(config string) *global.Application {
	global.OnceInitializeConfig(config)
	global.App.Log = global.InitializeLog()
	global.App.Db = database.NewChatDb()
	return global.App
}
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
	//完全匹配，直接返回
	if f, ok := ExtensionExactMatchReturnOutcome.GetFunc(prompt); ok {
		return f(prompt, userId)
	}
	//前缀匹配，直接返回
	if f, ok := getFuncFromPrefixMatchMap(ExtensionPrefixMatchReturnOutcome.GetFuncs(), prompt); ok {
		return f(prompt, userId)
	}
	//完全匹配，修饰prompt
	if f, ok := ExtensionExactMatchReturnPrompt.GetFunc(prompt); ok {
		newPrompt := f(prompt, userId)
		return answerByGpt(newPrompt, userId)
	}
	//前缀匹配，修饰prompt
	if f, ok := getFuncFromPrefixMatchMap(ExtensionPrefixMatchReturnPrompt.GetFuncs(), prompt); ok {
		newPrompt := f(prompt, userId)
		return answerByGpt(newPrompt, userId)
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

//匹配前缀是否相同
func getFuncFromPrefixMatchMap(m map[string]func(a, b string) string, source string) (func(a, b string) string, bool) {
	for k, f := range m {
		keyLen := len(k)
		if len(source) < keyLen {
			continue
		}
		if source[:keyLen] == k {
			return f, true
		}
	}
	return nil, false
}
