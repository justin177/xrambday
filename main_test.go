package main

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda/messages"
)

func TestRPCInvokeClient(t *testing.T) {
	port := os.Getenv("_LAMBDA_SERVER_PORT")
	if port == "" {
		port = "8888"
	}

	client, err := rpc.Dial("tcp", "localhost:"+port)
	if err != nil {
		t.Fatalf("dial rpc localhost:%s: %v", port, err)
	}
	defer client.Close()

	payload, err := json.Marshal(events.LambdaFunctionURLRequest{
		Version: "2.0",
		RawPath: "/",
	})
	if err != nil {
		t.Fatal(err)
	}

	req := messages.InvokeRequest{
		RequestId: "local-test",
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: time.Now().Add(2 * time.Second).Unix(),
		},
		Payload: payload,
	}

	var resp messages.InvokeResponse
	if err := client.Call("Function.Invoke", &req, &resp); err != nil {
		t.Fatalf("invoke rpc: %v", err)
	}
	if resp.Error != nil {
		t.Fatalf("lambda error: %s", resp.Error.Message)
	}

	var out events.LambdaFunctionURLResponse
	if err := json.Unmarshal(resp.Payload, &out); err != nil {
		t.Fatalf("unmarshal response: %v; payload=%s", err, string(resp.Payload))
	}
	if out.StatusCode != 200 {
		t.Fatalf("status code = %d, body=%q", out.StatusCode, out.Body)
	}
	if out.Body != "lambda context done, xray stopped\n" {
		t.Fatalf("unexpected body: %q", out.Body)
	}
}

func TestRPCPortIsNumeric(t *testing.T) {
	port := os.Getenv("_LAMBDA_SERVER_PORT")
	if port == "" {
		port = "8888"
	}
	if _, err := strconv.Atoi(port); err != nil {
		t.Fatalf("_LAMBDA_SERVER_PORT must be numeric: %q", port)
	}
	fmt.Println("testing lambda rpc port", port)
}
