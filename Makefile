build:
	@go build ./...

run:
	@docker-compose up --build tax-calculator

test:
	@go test ./... -v -count=1 -tags=unit

testIT-local:
	@INTERVIEW_SERVER=http://localhost:5000 go test ./integration -tags=integration -v -count=1

testIT:
	@docker-compose up --build integration-test

