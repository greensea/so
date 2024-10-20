package main

import (
	"fmt"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/greensea/so/config"
	"github.com/greensea/so/message"
)

var translateID atomic.Int64

func requestGPT(messages []message.GPTRequestMessage, model string, authorization string) ([]byte, error) {
	req := message.GPTRequest{
		Stream:    false,
		MaxTokens: 2048,
		Model:     model,
	}
	req.Messages = messages

	// 2. 发送请求
	url := config.Get().CompatiableOpenAIEndpoint
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authorization).
		SetBody(req).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("Unable to request <%s>: %v", url, err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("Request <%s> failed, got HTTP status %d.\n%s", url, resp.StatusCode(), string(resp.Body()))
	}

	return resp.Body(), nil
}
