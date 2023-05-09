package openai_chat_switch

import (
	"fmt"
	"testing"
)

func TestProcessPrompt(t *testing.T) {
	fmt.Print(processPrompt("报时 "))
}

func TestAnswer(t *testing.T) {

	NewGlobal("config.json")
	prompt := "模型切换到34"
	fmt.Print(Answer(prompt, "my"))
}
