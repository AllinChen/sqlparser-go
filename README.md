# sqlparser-go

#1. how to build
GOOS=darwin && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-mac *.go

CGO_ENABLED=0 && GOOS=linux && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-linux *.go

CGO_ENABLED=0 && GOOS=windows && GOARCH=amd64 && go build -o ./bin/sql-fingerprint-windows *.go
