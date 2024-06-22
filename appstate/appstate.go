package appstate

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type GetCurrIPMethod string

const (
	GET_CURR_IP_METHOD_NSLOOKUP GetCurrIPMethod = "NSLOOKUP"
	GET_CURR_IP_METHOD_CF       GetCurrIPMethod = "CF"
)

type AppState struct {
	cfApiKey  string
	domain    string
	subdomain string
	proxy     bool

	sleepInterval   time.Duration
	getCurrIPMethod GetCurrIPMethod

	logLevel slog.Level

	request *http.Request
}

func NewAppState() (*AppState, error) {
	if os.Getenv("API_KEY") == "" {
		return nil, fmt.Errorf("NewConfig: API_KEY is not set")
	}
	if os.Getenv("DOMAIN") == "" {
		return nil, fmt.Errorf("NewConfig: DOMAIN is not set")
	}

	getRealIPMethod := GET_CURR_IP_METHOD_NSLOOKUP
	if strings.ToLower(os.Getenv("GET_CURR_IP_METHOD")) == "cf" {
		getRealIPMethod = GET_CURR_IP_METHOD_CF
	}

	sleepInterval := time.Minute * 5
	var err error
	if os.Getenv("SLEEP_INTERVAL") != "" {
		sleepInterval, err = time.ParseDuration(os.Getenv("SLEEP_INTERVAL"))
		if err != nil {
			return nil, fmt.Errorf("NewConfig: SLEEP_INTERVAL is not a valid duration: %s", err)
		}
	}

	var logLevel slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	return &AppState{
		cfApiKey:  os.Getenv("API_KEY"),
		domain:    os.Getenv("DOMAIN"),
		subdomain: os.Getenv("SUBDOMAIN"),
		proxy:     strings.ToLower(os.Getenv("PROXY")) == "true",

		getCurrIPMethod: getRealIPMethod,
		sleepInterval:   sleepInterval,

		logLevel: logLevel,
	}, nil
}

func (as *AppState) GetApiKey() string {
	return as.cfApiKey
}
func (as *AppState) GetDomain() string {
	return as.domain
}
func (as *AppState) GetDDNSDomain() string {
	domain := as.domain
	if as.subdomain != "" {
		domain = as.subdomain + "." + domain
	}
	return domain
}
func (as *AppState) GetProxy() bool {
	return as.proxy
}

func (as *AppState) GetCurrIPMethod() GetCurrIPMethod {
	return as.getCurrIPMethod
}
func (as *AppState) GetSleepInterval() time.Duration {
	return as.sleepInterval
}

func (as *AppState) GetLogLevel() slog.Level {
	return as.logLevel
}
