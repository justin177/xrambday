package main

import (
	"log"
	"os"
)

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
}
