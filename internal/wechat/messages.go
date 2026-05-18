package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SubscribeMessage struct {
	ToUser     string                 `json:"touser"`
	TemplateID string                 `json:"template_id"`
	Page       string                 `json:"page,omitempty"`
	Data       map[string]MessageItem `json:"data"`
}

type MessageItem struct {
	Value string `json:"value"`
}

func SendSubscribeMessage(ctx context.Context, accessToken string, msg SubscribeMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	apiURL := "https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=" + accessToken
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Errcode != 0 {
		return fmt.Errorf("wechat error %d: %s", result.Errcode, result.Errmsg)
	}
	return nil
}
