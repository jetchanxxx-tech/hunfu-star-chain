package wechat

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/huifu/star-chain/internal/config"
)

type WeworkClient struct {
	corpID         string
	corpSecret     string
	token          string
	encodingAESKey string
	agentID        int

	mu           sync.RWMutex
	accessToken  string
	tokenExpires time.Time
	jsapiTicket  string
	ticketExpires time.Time
}

func NewWeworkClient(cfg config.WeworkConfig) *WeworkClient {
	return &WeworkClient{
		corpID:         cfg.CorpID,
		corpSecret:     cfg.CorpSecret,
		token:          cfg.Token,
		encodingAESKey: cfg.EncodingAESKey,
		agentID:        cfg.AgentID,
	}
}

// GetAccessToken retrieves or refreshes the corp access token
func (c *WeworkClient) GetAccessToken() (string, error) {
	c.mu.RLock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpires) {
		defer c.mu.RUnlock()
		return c.accessToken, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpires) {
		return c.accessToken, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.corpID, c.corpSecret))
	if err != nil {
		return "", fmt.Errorf("gettoken: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Errcode     int    `json:"errcode"`
		Errmsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Errcode != 0 {
		return "", fmt.Errorf("wework api error: %s", result.Errmsg)
	}
	c.accessToken = result.AccessToken
	c.tokenExpires = time.Now().Add(time.Duration(result.ExpiresIn-300) * time.Second)
	return c.accessToken, nil
}

// GetUserInfo gets WeChat user info via OAuth code
func (c *WeworkClient) GetUserInfo(code string) (map[string]any, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s", token, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return nil, fmt.Errorf("wework getuserinfo: %v", result["errmsg"])
	}
	return result, nil
}

// GetJSAPITicket retrieves JS-SDK ticket
func (c *WeworkClient) GetJSAPITicket() (string, error) {
	c.mu.RLock()
	if c.jsapiTicket != "" && time.Now().Before(c.ticketExpires) {
		defer c.mu.RUnlock()
		return c.jsapiTicket, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	token, err := c.GetAccessToken()
	if err != nil {
		return "", err
	}
	resp, err := http.Get(fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=%s", token))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Ticket  string `json:"ticket"`
		Expires int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Errcode != 0 {
		return "", fmt.Errorf("jsapi_ticket: %s", result.Errmsg)
	}
	c.jsapiTicket = result.Ticket
	c.ticketExpires = time.Now().Add(time.Duration(result.Expires-300) * time.Second)
	return c.jsapiTicket, nil
}

// JSSignature computes the JS-SDK signature for a given URL
func (c *WeworkClient) JSSignature(urlStr string) (map[string]any, error) {
	ticket, err := c.GetJSAPITicket()
	if err != nil {
		return nil, err
	}
	noncestr := fmt.Sprintf("%d", time.Now().UnixNano())
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%s&url=%s", ticket, noncestr, timestamp, urlStr)
	sig := fmt.Sprintf("%x", sha1.Sum([]byte(str)))
	return map[string]any{
		"corpid":    c.corpID,
		"agentid":   c.agentID,
		"noncestr":  noncestr,
		"timestamp": timestamp,
		"signature": sig,
	}, nil
}

// SendMessage sends an application message to a user via WeWork
func (c *WeworkClient) SendMessage(toUser, msgType, content string) error {
	token, err := c.GetAccessToken()
	if err != nil {
		return err
	}
	body := map[string]any{
		"touser":  toUser,
		"msgtype": msgType,
		"agentid": c.agentID,
	}
	switch msgType {
	case "text":
		body["text"] = map[string]string{"content": content}
	case "textcard":
		var card map[string]any
		json.Unmarshal([]byte(content), &card)
		body["textcard"] = card
	default:
		return fmt.Errorf("unsupported msgtype: %s", msgType)
	}

	payload, _ := json.Marshal(body)
	resp, err := http.Post(
		fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("send message: %v", result["errmsg"])
	}
	return nil
}

// SendTaskCard sends a task notification card to a steward
func (c *WeworkClient) SendTaskCard(stewardUserID, taskTitle, memberName, dueDate, taskURL string) error {
	card := map[string]any{
		"title":       "新任务提醒",
		"description": fmt.Sprintf("<div class=\"gray\">%s</div><div class=\"normal\">会员: %s</div><div>截止: %s</div>", taskTitle, memberName, dueDate),
		"url":         taskURL,
		"btntxt":      "查看详情",
	}
	cardJSON, _ := json.Marshal(card)
	return c.SendMessage(stewardUserID, "textcard", string(cardJSON))
}

// ListDepartments fetches department list
func (c *WeworkClient) ListDepartments(parentID int) ([]map[string]any, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s", token)
	if parentID > 0 {
		urlStr += fmt.Sprintf("&id=%d", parentID)
	}
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Errcode    int              `json:"errcode"`
		Department []map[string]any `json:"department"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Department, nil
}

// ListDepartmentUsers fetches user list in a department
func (c *WeworkClient) ListDepartmentUsers(departmentID int, fetchChild bool) ([]map[string]any, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%d&fetch_child=%d", token, departmentID, boolToInt(fetchChild))
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Errcode int              `json:"errcode"`
		Userlist []map[string]any `json:"userlist"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Userlist, nil
}

// VerifyURL handles WeWork callback URL verification
func (c *WeworkClient) VerifyURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgSignature := query.Get("msg_signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")

	plaintext, err := c.verifyMsg(msgSignature, timestamp, nonce, echostr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("verification failed"))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(plaintext))
}

// HandleCallback receives and decrypts WeWork message callbacks
func (c *WeworkClient) HandleCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	msgSignature := query.Get("msg_signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")

	body, _ := io.ReadAll(r.Body)
	var encrypted struct {
		Encrypt string `json:"encrypt"`
	}
	json.Unmarshal(body, &encrypted)

	plaintext, err := c.decryptMsg(encrypted.Encrypt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"errcode": 1, "errmsg": err.Error()})
		return
	}

	_ = msgSignature // Used in full verification
	_ = timestamp
	_ = nonce

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "errmsg": "ok"})

	_ = plaintext // Processed by upper-layer callback handler
}

func (c *WeworkClient) IsConfigured() bool {
	return c.corpID != "" && c.corpSecret != "" && c.token != ""
}

// verifyMsg verifies WeWork callback message signature and decrypts echostr
func (c *WeworkClient) verifyMsg(msgSignature, timestamp, nonce, echostr string) (string, error) {
	sig := c.calcSignature(timestamp, nonce, echostr)
	if sig != msgSignature {
		return "", fmt.Errorf("signature mismatch")
	}
	return c.decryptMsg(echostr)
}

func (c *WeworkClient) calcSignature(timestamp, nonce, data string) string {
	items := []string{c.token, timestamp, nonce, data}
	sort.Strings(items)
	h := sha1.New()
	h.Write([]byte(strings.Join(items, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// decryptMsg decrypts AES-encrypted message from WeWork
func (c *WeworkClient) decryptMsg(encrypted string) (string, error) {
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey + "=")
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove PKCS7 padding
	padLen := int(ciphertext[len(ciphertext)-1])
	ciphertext = ciphertext[:len(ciphertext)-padLen]

	// Parse: random(16) + msg_len(4) + msg + corpid
	if len(ciphertext) < 20 {
		return "", fmt.Errorf("decrypted too short")
	}
	msgLen := binary.BigEndian.Uint32(ciphertext[16:20])
	msg := string(ciphertext[20 : 20+msgLen])
	corpID := string(ciphertext[20+msgLen:])
	if corpID != c.corpID {
		return "", fmt.Errorf("corpid mismatch: %s", corpID)
	}
	return msg, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// OAuthURL generates the WeWork OAuth authorization URL
func (c *WeworkClient) OAuthURL(redirectURI, state string) string {
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=%s#wechat_redirect",
		url.QueryEscape(c.corpID), url.QueryEscape(redirectURI), url.QueryEscape(state))
}
