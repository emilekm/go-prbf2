test:
	go run gotest.tools/gotestsum@v1.11.0 -- -count=1 ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
