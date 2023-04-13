test:
	go generate ./...
	go test ./... -coverprofile ./coverage.out

cover: test
	go tool cover -func ./coverage.out
