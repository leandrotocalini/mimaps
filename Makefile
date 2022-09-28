coverage:
	go test -coverprofile=coverage.out ./... ;    go tool cover -func=coverage.out

test:
	go test ./... -race


lint: 
	golangci-lint run .