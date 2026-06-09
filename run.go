package main

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/core"
)

//go:embed testdata/xrambday.json
var embeddedConfig []byte

func isHTTPSConfigSource(source string) bool {
	return strings.HasPrefix(strings.ToLower(source), "https://")
}

func fetchHTTPSConfig(source string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	if isTruthyEnv("CONFIG_TLS_INSECURE") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("unexpected HTTP status: ", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func isTruthyEnv(name string) bool {
	switch strings.ToLower(os.Getenv(name)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func startXray() (core.Server, error) {
	configSource := os.Getenv("CONFIG")
	if configSource == "" {
		log.Println("Using embedded config")
		return newXrayServer(embeddedConfig, "embedded")
	}
	log.Println("Using config from CONFIG")

	var configBytes []byte
	var err error
	if isHTTPSConfigSource(configSource) {
		configBytes, err = fetchHTTPSConfig(configSource)
		if err != nil {
			return nil, errors.New("failed to fetch remote config from CONFIG").Base(err)
		}
	} else if isJSONConfigSource(configSource) {
		configBytes = []byte(configSource)
	} else {
		configBytes, err = os.ReadFile(configSource)
		if err != nil {
			return nil, errors.New("failed to read config from CONFIG").Base(err)
		}
	}
	return newXrayServer(configBytes, "CONFIG")
}

func isJSONConfigSource(source string) bool {
	return strings.HasPrefix(strings.TrimSpace(source), "{")
}

func newXrayServer(configBytes []byte, sourceLabel string) (core.Server, error) {
	c, err := core.LoadConfig("json", bytes.NewReader(configBytes))
	if err != nil {
		return nil, errors.New("failed to load config from ", sourceLabel).Base(err)
	}

	server, err := core.New(c)
	if err != nil {
		return nil, errors.New("failed to create server").Base(err)
	}

	return server, nil
}
