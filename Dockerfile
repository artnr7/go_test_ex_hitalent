FROM golang:tip-alpine3.23 AS builder
WORKDIR /chatapi
COPY ./chat_api/ .
RUN go mod download && go build -o chat_api ./cmd/chat-api/main.go

# --- #

FROM alpine:3.23
WORKDIR /backend
COPY --from=builder ./chatapi/chat_api .
COPY --from=builder ./chatapi/migrations ./migrations
RUN apk add --no-cache curl

ENTRYPOINT ["./chat_api"] 
