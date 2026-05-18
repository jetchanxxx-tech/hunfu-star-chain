package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/huifu/star-chain/internal/config"
)

type MiniprogramClient struct {
	appID     string
	appSecret string
	client    *http.Client
}

type SessionResult struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
}

func NewMiniprogramClient(cfg config.MiniprogramConfig) *MiniprogramClient {
	return &MiniprogramClient{
		appID:     cfg.AppID,
		appSecret: cfg.AppSecret,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *MiniprogramClient) Code2Session(ctx context.Context, code string) (*SessionResult, error) {
	params := url.Values{}
	params.Set("appid", c.appID)
	params.Set("secret", c.appSecret)
	params.Set("js_code", code)
	params.Set("grant_type", "authorization_code")

	apiURL := "https://api.weixin.qq.com/sns/jscode2session?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("code2session: %w", err)
	}
	defer resp.Body.Close()

	var sr SessionResult
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if sr.Errcode != 0 {
		return nil, fmt.Errorf("wechat error %d: %s", sr.Errcode, sr.Errmsg)
	}
	return &sr, nil
}

func (c *MiniprogramClient) GetAccessToken(ctx context.Context) (*AccessTokenResult, error) {
	params := url.Values{}
	params.Set("appid", c.appID)
	params.Set("secret", c.appSecret)
	params.Set("grant_type", "client_credential")

	apiURL := "https://api.weixin.qq.com/cgi-bin/token?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	defer resp.Body.Close()

	var at AccessTokenResult
	if err := json.NewDecoder(resp.Body).Decode(&at); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	if at.Errcode != 0 {
		return nil, fmt.Errorf("wechat error %d", at.Errcode)
	}
	return &at, nil
}
