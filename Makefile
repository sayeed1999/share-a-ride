fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run

vet:
	go vet ./...

test:
	go test ./...

coverage:
	go test -cover ./...

sec:
	gosec ./...

vuln:
	govulncheck ./...

doc:
	godoc -http:.6060

# mocks:
#   mockgen -source=./pkg/repository/repository.go -destination=./pkg/repository/mock_repository.go -package=repository
#   mockgen -source=./pkg/service/service.go -destination=./pkg/service/mock_service.go -package=service

all: fmt lint vet test coverage sec vuln doc
