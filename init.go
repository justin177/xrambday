package main

import (
	"bytes"
	_ "embed"
	"log"
	"os"

	"github.com/xtls/xray-core/main/confloader/external"
)

//go:embed init.conf
var configBytes []byte

func init() {
	if _, ok := os.LookupEnv("NO_LOG_NULL"); !ok {
		f, err := os.Open(os.DevNull)
		if err != nil {
			panic(err)
		}
		os.Stdout = f
		os.Stderr = f
		log.SetOutput(os.Stdout)
	}
	if _, ok := os.LookupEnv("USE_CONFIG_EMBED"); ok {
		r, w, err := os.Pipe()
		if err != nil {
			panic(err)
		}
		os.Stdin = r

		if v, ok := os.LookupEnv(string(configBytes)); ok && len(v) > 0 {
			configBytes = []byte(v)
		}

		if bytes.HasPrefix(configBytes, []byte("https://")) {
			configBytes, err = external.FetchHTTPContent(string(configBytes))
			if err != nil {
				panic(err)
			}
		}

		_, err = w.Write(configBytes)
		if err != nil {
			panic(err)
		}
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}
}
