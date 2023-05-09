package openai_chat_switch

import (
	"fmt"
	"testing"
)

func TestChangeModel(t *testing.T) {
	prompt := "模型切换到3.5"
	fmt.Println(changeModel(prompt, "uu"))
}

func TestChangeTemperature(t *testing.T) {
	prompt := "Temperature切换到0.5"
	fmt.Println(changeTemperature(prompt, "uu"))
}
