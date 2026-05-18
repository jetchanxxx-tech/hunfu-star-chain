package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
	Wechat   WechatConfig   `yaml:"wechat"`
	AI       AIConfig       `yaml:"ai"`
	Storage  StorageConfig  `yaml:"storage"`
	VoiceCall ProviderConfig `yaml:"voice_call"`
	SMS      ProviderConfig `yaml:"sms"`
	Push     ProviderConfig `yaml:"push"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"`
	DSN          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type RabbitMQConfig struct {
	URL string `yaml:"url"`
}

type WechatConfig struct {
	Miniprogram MiniprogramConfig `yaml:"miniprogram"`
	Wework      WeworkConfig      `yaml:"wework"`
}

type MiniprogramConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type WeworkConfig struct {
	CorpID         string `yaml:"corp_id"`
	Token          string `yaml:"token"`
	EncodingAESKey string `yaml:"encoding_aes_key"`
}

type AIConfig struct {
	DefaultProvider  string              `yaml:"default_provider"`
	Providers        map[string]AIProvider `yaml:"providers"`
	EmergencyKeywords []string           `yaml:"emergency_keywords"`
	EmotionDetection EmotionConfig       `yaml:"emotion_detection"`
}

type AIProvider struct {
	Endpoint   string `yaml:"endpoint"`
	APIKey     string `yaml:"api_key"`
	Model      string `yaml:"model"`
	Timeout    string `yaml:"timeout"`
	MaxRetries int    `yaml:"max_retries"`
}

type EmotionConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Provider string `yaml:"provider"`
}

type StorageConfig struct {
	Default string           `yaml:"default"`
	Local   LocalStorageConfig `yaml:"local"`
	COS     COSConfig        `yaml:"cos"`
}

type LocalStorageConfig struct {
	Path string `yaml:"path"`
}

type COSConfig struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	SecretID  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}

type ProviderConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Provider string `yaml:"provider"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	data = []byte(os.ExpandEnv(string(data)))
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	cfg.setDefaults()
	return &cfg, nil
}

func (c *Config) setDefaults() {
	c.Database.DSN = expandEnv(c.Database.DSN)
	c.Redis.Addr = expandEnv(c.Redis.Addr)
	c.Redis.Password = expandEnv(c.Redis.Password)
	c.RabbitMQ.URL = expandEnv(c.RabbitMQ.URL)
	for k, p := range c.AI.Providers {
		p.APIKey = expandEnv(p.APIKey)
		c.AI.Providers[k] = p
	}
}

func expandEnv(s string) string {
	if strings.HasPrefix(s, "${") && strings.HasSuffix(s, "}") {
		inner := s[2 : len(s)-1]
		parts := strings.SplitN(inner, ":-", 2)
		name := parts[0]
		if val, ok := os.LookupEnv(name); ok {
			return val
		}
		if len(parts) == 2 {
			return parts[1]
		}
		return s
	}
	return os.ExpandEnv(s)
}
