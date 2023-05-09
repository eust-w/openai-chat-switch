package gpt

import (
	"context"
	"fmt"
	"github.com/eust-w/openai-chat-switch/global"
	"net/http"
	"net/url"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const (
	DefaultTemperature = 0.6
	DefaultAiRole      = "AI"
	DefaultHumanRole   = "Human"
)

type ChatGpt struct {
	client        *openai.Client
	ctx           context.Context
	userId        string
	maxPromptLen  int           // 最大问题长度
	maxansewerLen int           // 最大答案长度
	maxText       int           // 最大文本 = 问题 + 回答, 接口限制
	timeOut       time.Duration // 超时时间, 0表示不超时
	cancel        func()

	ChatContext []string
}

func New(userId string) *ChatGpt {
	var ctx context.Context
	var cancel func()

	if global.App.Config.Timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(global.App.Config.Timeout))
	}
	timeOutChan := make(chan struct{}, 1)
	go func() {
		<-ctx.Done()
		timeOutChan <- struct{}{} // 发送超时信号，或是提示结束，用于聊天机器人场景，配合GetTimeOutChan() 使用
	}()

	config := openai.DefaultConfig(global.App.Config.ApiKey)
	if global.App.Config.Proxy != "" {
		config.HTTPClient.Transport = &http.Transport{
			// 设置代理
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(global.App.Config.Proxy)
			}}
	}
	if global.App.Config.CustomUrl != "" {
		config.BaseURL = global.App.Config.CustomUrl + "/v1"
	}

	return &ChatGpt{
		client:        openai.NewClientWithConfig(config),
		ctx:           ctx,
		userId:        userId,
		maxPromptLen:  2048, // 最大问题长度
		maxansewerLen: 2048, // 最大答案长度
		maxText:       4096, // 最大文本 = 问题 + 回答, 接口限制
		timeOut:       time.Duration(global.App.Config.Timeout),
		cancel: func() {
			cancel()
		},
		ChatContext: global.App.Db.GetContext(userId),
	}
}

func (c *ChatGpt) Close() {
	c.cancel()
}

func (c ChatGpt) getModel() string {
	model := global.App.Db.GetModel(c.userId)
	if model != "" {
		return model
	}
	return global.App.Config.Model
}

func (c ChatGpt) getTemperature() float32 {
	temperature := global.App.Db.GetTemperature(c.userId)
	if temperature != 0 {
		return temperature
	}
	return DefaultTemperature
}

func (c *ChatGpt) ChatWithContext(prompt string) (answer string, err error) {
	prompt = prompt + "."
	promptTable := c.ChatContext
	subSetPromptTable := GetMaxSubset(promptTable, c.maxPromptLen)
	realPrompt := fmt.Sprint(strings.Join(subSetPromptTable, "\n")) + "\n\n" + DefaultHumanRole + ":" + prompt + "\n" + DefaultAiRole + ":"
	global.App.Log.Info("realPrompt is:", len(realPrompt))
	model := c.getModel()
	temperature := c.getTemperature()
	if model == openai.GPT3Dot5Turbo0301 ||
		model == openai.GPT3Dot5Turbo ||
		model == openai.GPT4 || model == openai.GPT40314 ||
		model == openai.GPT432K || model == openai.GPT432K0314 {
		global.App.Log.Info("gpt model:", model)
		req := openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: realPrompt,
				},
			},
			MaxTokens:   3072,
			Temperature: temperature,
			User:        c.userId,
		}
		resp, err := c.client.CreateChatCompletion(c.ctx, req)
		if err != nil {
			global.App.Log.Error("CreateChatCompletion err", err)
			return "", err
		}
		resp.Choices[0].Message.Content = formatAnswer(resp.Choices[0].Message.Content)
		newPromptTable := append(subSetPromptTable, "\n"+DefaultHumanRole+":"+prompt, "\n"+DefaultAiRole+":"+resp.Choices[0].Message.Content)
		global.App.Db.SetContext(c.userId, newPromptTable)
		return resp.Choices[0].Message.Content, nil
	} else {
		req := openai.CompletionRequest{
			Model:       model,
			MaxTokens:   c.maxansewerLen,
			Prompt:      realPrompt,
			Temperature: temperature,
			User:        c.userId,
			Stop:        []string{DefaultAiRole + ":", DefaultHumanRole + ":"},
		}
		resp, err := c.client.CreateCompletion(c.ctx, req)
		if err != nil {
			global.App.Log.Error("CreateCompletion err", err)
			return "", err
		}
		resp.Choices[0].Text = formatAnswer(resp.Choices[0].Text)
		newPromptTable := append(subSetPromptTable, "\n"+DefaultHumanRole+":"+prompt, "\n"+DefaultAiRole+":"+resp.Choices[0].Text)
		global.App.Db.SetContext(c.userId, newPromptTable)
		return resp.Choices[0].Text, nil
	}
}
