package shared

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type IConfig struct {
	AppEnv             string
	BasePort           int
	BaseUrl            string
	BaseToken          string
	BaseMultiProcess   bool
	BaseAllowedOrigins []string
	TelegramBotToken   string
	TelegramChatID     string
}

const (
	EnvAppEnv             = "APP_ENV"
	EnvBasePort           = "BASE_PORT"
	EnvBaseUrl            = "BASE_URL"
	EnvBaseToken          = "BASE_TOKEN"
	EnvBaseMultiProcess   = "BASE_MULTI_PROCESS"
	EnvBaseAllowedOrigins = "BASE_ALLOWED_ORIGINS"
	EnvTelegramBotToken   = "TELEGRAM_BOT_TOKEN"
	EnvTelegramChatID     = "TELEGRAM_CHAT_ID"
)

var (
	Config IConfig
)

func (c *IConfig) Validate() error {
	if c.AppEnv != "development" && c.AppEnv != "production" {
		return fmt.Errorf("APP_ENV out of range: %s", c.AppEnv)
	}

	if c.BasePort <= 0 || c.BasePort > 65535 {
		return fmt.Errorf("BASE_PORT must be between 1 and 65535")
	}

	if c.BaseUrl == "" {
		return fmt.Errorf("BASE_URL is required")
	}

	if c.BaseToken == "" {
		return fmt.Errorf("BASE_TOKEN is required")
	}

	if c.BaseMultiProcess {
		c.BaseMultiProcess = true
	}

	if len(c.BaseAllowedOrigins) == 0 {
		return fmt.Errorf("BASE_ALLOWED_ORIGINS is required")
	}

	if c.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	if c.TelegramChatID == "" {
		return fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}

	return nil
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}

func Load() error {
	CapturePanic()

	if !runningInDocker() {
		_ = godotenv.Load()
	}

	var err error
	cfg := &IConfig{}

	cfg.AppEnv = GetEnv(EnvAppEnv, "development")
	if cfg.AppEnv == "" {
		return fmt.Errorf("environment variable %s is required", EnvAppEnv)
	}

	basePortStr := GetEnv(EnvBasePort, "3000")
	cfg.BasePort, err = strconv.Atoi(basePortStr)
	if err != nil {
		return fmt.Errorf("environment variable %s must be a valid number", EnvBasePort)
	}

	if cfg.BasePort <= 0 || cfg.BasePort > 65535 {
		return fmt.Errorf("environment variable %s must be between 1 and 65535", EnvBasePort)
	}

	cfg.BaseUrl = GetEnv(EnvBaseUrl, "http://localhost:3000")
	if cfg.BaseUrl == "" {
		return fmt.Errorf("environment variable %s is required", EnvBaseUrl)
	}

	cfg.BaseToken = GetEnv(EnvBaseToken, "")
	if cfg.BaseToken == "" {
		return fmt.Errorf("environment variable %s is required", EnvBaseToken)
	}

	cfg.BaseMultiProcess = GetEnv(EnvBaseMultiProcess, "false") == "true"
	if cfg.BaseMultiProcess {
		cfg.BaseMultiProcess = true
	}

	allowedOriginsStr := GetEnv(EnvBaseAllowedOrigins, "")
	if allowedOriginsStr == "" {
		return fmt.Errorf("environment variable %s is required", EnvBaseAllowedOrigins)
	}
	cfg.BaseAllowedOrigins = strings.Split(allowedOriginsStr, ",")
	for i, origin := range cfg.BaseAllowedOrigins {
		cfg.BaseAllowedOrigins[i] = strings.TrimSpace(origin)
		if cfg.BaseAllowedOrigins[i] == "" {
			return fmt.Errorf("environment variable %s contains an empty origin", EnvBaseAllowedOrigins)
		}
		// Detect IPv6 by presence of ':'
		if strings.Contains(cfg.BaseAllowedOrigins[i], ":") {
			return fmt.Errorf("environment variable %s contains an IPv6 origin, which is not allowed: %s", EnvBaseAllowedOrigins, cfg.BaseAllowedOrigins[i])
		}
	}

	cfg.TelegramBotToken = GetEnv(EnvTelegramBotToken, "")
	if cfg.TelegramBotToken == "" {
		return fmt.Errorf("environment variable %s is required", EnvTelegramBotToken)
	}

	cfg.TelegramChatID = GetEnv(EnvTelegramChatID, "")
	if cfg.TelegramChatID == "" {
		return fmt.Errorf("environment variable %s is required", EnvTelegramChatID)
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	Config = *cfg

	return nil
}

func runningInDocker() bool {
	return checkIfRunningInEnvironment("/proc/1/cgroup", "docker", "kubepods")
}

func checkIfRunningInEnvironment(path, kw1, kw2 string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, kw1) || strings.Contains(line, kw2) {
			return true
		}
	}

	return s.Err() == nil && false
}
