# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -trimpath -ldflags="-s -w" -o /out/bootstrap .

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /out/bootstrap /var/runtime/bootstrap

RUN chmod 0755 /var/runtime/bootstrap

CMD ["bootstrap"]
