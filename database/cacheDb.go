package database

import (
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	contextFlag     = "contextFlag/"
	temperatureFlag = "temperature/"
	topPFlag        = "topp/"
	modelFlag       = "model/"
)

type ChatDb struct {
	db *cache.Cache
}

func NewChatDb() *ChatDb {
	return &ChatDb{db: cache.New(time.Hour*24, time.Hour*24*7)}
}

// SetContext 设置会话上下文
func (s *ChatDb) SetContext(userId string, content []string) {
	s.db.Set(contextFlag+userId, content, cache.DefaultExpiration)
}

// GetContext 获取会话上下文
func (s *ChatDb) GetContext(userId string) []string {
	context, ok := s.db.Get(contextFlag + userId)
	if !ok {
		return nil
	}
	return context.([]string)
}

// DeleteContext 删除上下文
func (s *ChatDb) DeleteContext(userId string) {
	s.db.Delete(contextFlag + userId)
}

// SetTemperature 设置会话Temperature
func (s *ChatDb) SetTemperature(userId string, temperature float32) {
	s.db.Set(temperatureFlag+userId, temperature, cache.DefaultExpiration)
}

// GetTemperature 获取会话Temperature
func (s *ChatDb) GetTemperature(userId string) float32 {
	temperature, ok := s.db.Get(temperatureFlag + userId)
	if !ok {
		return 0
	}
	return temperature.(float32)
}

// SetTopP 设置会话TopP
func (s *ChatDb) SetTopP(userId string, topP float32) {
	s.db.Set(topPFlag+userId, topP, cache.DefaultExpiration)
}

// GetTopP  获取会话TopP
func (s *ChatDb) GetTopP(userId string) float32 {
	topP, ok := s.db.Get(topPFlag + userId)
	if !ok {
		return 0
	}
	return topP.(float32)
}

// SetModel 设置会话Model
func (s *ChatDb) SetModel(userId string, model string) {
	s.db.Set(modelFlag+userId, model, cache.DefaultExpiration)
}

// GetModel 获取会话Model
func (s *ChatDb) GetModel(userId string) string {
	model, ok := s.db.Get(modelFlag + userId)
	if !ok {
		return ""
	}
	return model.(string)
}
