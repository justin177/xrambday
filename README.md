# xrambday

AWS Lambda Function URL image for running an Xray server during an activation
window.

This is a standalone Lambda image. The Xray run code is copied from the `render`
project, and this project's `main.go` only wraps it with an AWS Lambda Function
URL handler. When the Function URL receives a request, the handler calls
`startXray()`, starts the returned server, waits for the activation window, then
closes the server.

AWS documents `public.ecr.aws/lambda/provided:al2023` as the OS-only base image
for compiled Go Lambda container images.

## Build

```sh
docker buildx build \
  --platform linux/amd64 \
  -t xrambday:latest \
  .
```

Use `--platform linux/arm64` for ARM Lambda functions.

## Configure Lambda

Create the Lambda from this container image and enable a Function URL.
Set `CONFIG` to the Xray config path, matching the `render` project behavior.

The Lambda timeout should be longer than 840 seconds. For a 900 second timeout,
the 840 second activation window leaves 60 seconds for cleanup.

If `CONFIG` is not set, the binary falls back to `config.json` in the current
working directory, then to stdin, matching the local run behavior. The sample
Lambda config lives at `testdata/xrambday.json`; deployment ZIPs only contain the
`bootstrap` binary.

## Local test

```sh
docker run --rm -p 9000:8080 -e CONFIG=/var/task/config.json xrambday:latest
```

Invoke through the Lambda Runtime Interface Emulator:

```sh
curl -sS 'http://localhost:9000/2015-03-31/functions/function/invocations' \
  -H 'content-type: application/json' \
  -d '{}'
```

Expected response:

```text
activation window ended, xray stopped
```

## Upstream merge

The upstream reverse portal merge fragment and helper script are under
`testdata/`:

```sh
testdata/merge-upstream-reverse.sh upstream.json
```
