package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type EmotionResult struct {
	Tag        string `json:"tag"`        // normal / anxious / complaint
	Confidence int    `json:"confidence"` // 0-100
}

func (p *openAICompatProvider) DetectEmotion(ctx context.Context, userMessage string) (*EmotionResult, error) {
	resp, err := p.Chat(ctx, ChatRequest{
		Messages: []Message{
			{Role: "system", Content: `你是一个情绪分析助手。分析用户的情绪并返回 JSON。
格式: {"tag":"normal|anxious|complaint","confidence":0-100}
用户是医院的患者或家属，焦虑可能源于健康担忧，投诉可能源于服务不满。
只返回 JSON，不要其他内容。`},
			{Role: "user", Content: userMessage},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response")
	}
	content := extractJSON(resp.Choices[0].Message.Content)
	var er EmotionResult
	if err := json.Unmarshal([]byte(content), &er); err != nil {
		return &EmotionResult{Tag: "normal", Confidence: 0}, nil
	}
	return &er, nil
}

func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if start := strings.Index(s, "{"); start >= 0 {
		if end := strings.LastIndex(s, "}"); end > start {
			return s[start : end+1]
		}
	}
	return s
}
