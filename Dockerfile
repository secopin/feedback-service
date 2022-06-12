FROM golang:1.17.3-alpine AS build
RUN apk --no-cache add upx git
WORKDIR /app
COPY go.mod main.go /app/
RUN go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v ./... && go build -ldflags="-s -w" -o ./out/app /app/ && upx ./out/app

FROM alpine:3.11 AS runtime
COPY --from=build /app/out/app /usr/local/bin/app
CMD ["/usr/local/bin/app"]
