build:
	go build ./...

test:
	go test ./... -v -count=1 -tags=unit

testIT:
	go test ./... -v -count=1 -tags=integration