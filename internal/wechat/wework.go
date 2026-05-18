package wechat

import (
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/config"
)

type WeworkClient struct {
	corpID         string
	token          string
	encodingAESKey string
}

func NewWeworkClient(cfg config.WeworkConfig) *WeworkClient {
	return &WeworkClient{
		corpID:         cfg.CorpID,
		token:          cfg.Token,
		encodingAESKey: cfg.EncodingAESKey,
	}
}

// VerifyURL validates the WeWork callback URL registration
func (c *WeworkClient) VerifyURL(w http.ResponseWriter, r *http.Request) {
	// WeWork URL verification: echostr verification via signature
	// P0 stub: return verification success
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

// HandleCallback receives WeWork message callbacks
func (c *WeworkClient) HandleCallback(w http.ResponseWriter, r *http.Request) {
	var msg map[string]any
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// P0 stub: log and ack
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"errcode": 0, "errmsg": "ok"})
}

func (c *WeworkClient) IsConfigured() bool {
	return c.corpID != "" && c.token != ""
}
