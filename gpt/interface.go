package gpt

import (
	"github.com/avast/retry-go"
	"time"
)

func Answer(question, userId string) (answer string) {
	var gpt *ChatGpt
	var err error
	gpt = New(userId)
	retryStrategy := []retry.Option{
		retry.Delay(100 * time.Millisecond),
		retry.Attempts(3),
		retry.LastErrorOnly(true)}
	// 使用重试策略进行重试
	err = retry.Do(
		func() error {
			answer, err = gpt.ChatWithContext(question)
			if err != nil {
				return err
			}
			return nil
		},
		retryStrategy...)
	return
}
