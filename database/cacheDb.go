package database

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	context = "context/"
	mode    = "mode/"
)

type ChatDb struct {
	db *cache.Cache
}

func NewChatDb() *ChatDb {
	return &ChatDb{db: cache.New(time.Hour*1, time.Hour*24*7)}
}

// SetContext 设置会话上下文
func (s *ChatDb) SetContext(userId string, content []string) {
	s.db.Set(context+userId, content, cache.DefaultExpiration)
}

// GetContext 获取会话上下文
func (s *ChatDb) GetContext(userId string) []string {
	context, ok := s.db.Get(context + userId)
	if !ok {
		return nil
	}
	return context.([]string)
}

// DeleteContext 删除上下文
func (s *ChatDb) DeleteContext(userId string) {
	s.db.Delete(context + userId)
}

// GetMode 获取对话模式
func (s *ChatDb) GetMode(userId string) string {
	mode, ok := s.db.Get(mode + userId)
	if !ok {
		return ""
	}
	return mode.(string)
}

// SetMode 设置对话模式
func (s *ChatDb) SetMode(userId string, mode string) {
	s.db.Set(mode+userId, mode, cache.DefaultExpiration)
}
