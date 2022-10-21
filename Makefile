vet:
	go vet ./...

test:
	go test ./...

testc:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download
