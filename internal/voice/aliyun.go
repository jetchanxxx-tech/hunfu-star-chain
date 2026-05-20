package voice

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type AliyunVoiceProvider struct {
	accessKeyID     string
	accessKeySecret string
	ttsAppKey       string
	endpoint        string
	httpClient      *http.Client
}

func NewAliyunVoiceProvider(accessKeyID, accessKeySecret, ttsAppKey string) *AliyunVoiceProvider {
	return &AliyunVoiceProvider{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		ttsAppKey:       ttsAppKey,
		endpoint:        "https://dyvmsapi.aliyuncs.com",
		httpClient:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *AliyunVoiceProvider) Name() string { return "aliyun" }

func (p *AliyunVoiceProvider) MakeCall(phone, dialogueText, templateID string, maxRetries int) (*CallResult, error) {
	params := map[string]string{
		"Action":           "SingleCallByTts",
		"Version":          "2017-05-25",
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   uuid.New().String(),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"AccessKeyId":      p.accessKeyID,
		"CalledNumber":     phone,
		"CalledShowNumber": "95187",
		"TtsCode":          templateID,
	}
	if dialogueText != "" {
		params["TtsParam"] = dialogueText
	}

	queryStr := p.buildQuery(params)
	signature := p.sign("GET&%2F&" + url.QueryEscape(queryStr))
	queryStr += "&Signature=" + url.QueryEscape(signature)

	reqURL := p.endpoint + "/?" + queryStr
	resp, err := p.httpClient.Get(reqURL)
	if err != nil {
		return &CallResult{Status: "failed", FailReason: err.Error()}, nil
	}
	defer resp.Body.Close()

	var result struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		CallID  string `json:"CallId"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Code != "OK" {
		return &CallResult{Status: "failed", FailReason: result.Message}, nil
	}
	return &CallResult{CallID: result.CallID, Status: "ringing"}, nil
}

func (p *AliyunVoiceProvider) QueryCallStatus(callID string) (*CallResult, error) {
	params := map[string]string{
		"Action":           "QueryCallDetailByCallId",
		"Version":          "2017-05-25",
		"Format":           "JSON",
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureVersion": "1.0",
		"SignatureNonce":   uuid.New().String(),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"AccessKeyId":      p.accessKeyID,
		"CallId":           callID,
	}

	queryStr := p.buildQuery(params)
	signature := p.sign("GET&%2F&" + url.QueryEscape(queryStr))
	queryStr += "&Signature=" + url.QueryEscape(signature)

	resp, err := p.httpClient.Get(p.endpoint + "/?" + queryStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code     string `json:"Code"`
		Duration int    `json:"Duration"`
		Status   string `json:"Status"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	callStatus := "completed"
	if result.Status == "UNANSWERED" {
		callStatus = "no_answer"
	} else if result.Status == "REJECT" {
		callStatus = "rejected"
	}
	return &CallResult{CallID: callID, Status: callStatus, DurationSeconds: result.Duration}, nil
}

func (p *AliyunVoiceProvider) CallbackHandler(rawBody []byte) (*CallResult, error) {
	var cb struct {
		CallID   string `json:"call_id"`
		Status   string `json:"status"`
		Duration int    `json:"duration"`
	}
	json.Unmarshal(rawBody, &cb)
	return &CallResult{CallID: cb.CallID, Status: cb.Status, DurationSeconds: cb.Duration}, nil
}

func (p *AliyunVoiceProvider) buildQuery(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(params[k]))
	}
	return strings.Join(parts, "&")
}

func (p *AliyunVoiceProvider) sign(stringToSign string) string {
	mac := hmac.New(sha1.New, []byte(p.accessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
