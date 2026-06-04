package main

import (
	"context"
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/xtls/xray-core/main/distro/all"
)

const (
	defaultActivationSeconds = 840
	contentTypeTextPlainUTF8 = "text/plain; charset=utf-8"
)

func handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("starting xray")

	server, err := startXray()
	if err != nil {
		return text(500, "failed to start xray: "+err.Error()+"\n"), nil
	}

	if err := server.Start(); err != nil {
		server.Close()
		return text(500, "failed to start xray: "+err.Error()+"\n"), nil
	}

	runtime.GC()
	debug.FreeOSMemory()

	stopXray := func() {
		log.Println("stopping xray")
		if err := server.Close(); err != nil {
			log.Printf("xray stopped with error: %v", err)
			return
		}
		log.Println("xray stopped")
	}

	timer := time.NewTimer(defaultActivationSeconds * time.Second)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		stopXray()
		return text(200, "lambda context done, xray stopped\n"), nil

	case <-timer.C:
		stopXray()
		return text(200, "activation window ended, xray stopped\n"), nil
	}
}

func text(code int, body string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{
		StatusCode: code,
		Headers: map[string]string{
			"content-type": contentTypeTextPlainUTF8,
		},
		Body: body,
	}
}

func main() {
	lambda.Start(handler)
}
