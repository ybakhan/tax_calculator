# Start from the latest golang base image
FROM golang:latest as builder

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /tax_calculator ./cmd

# Build the service image
FROM alpine:latest as service
COPY --from=builder /tax_calculator /tax_calculator
EXPOSE 8080
CMD ["/tax_calculator"]

# Build the integration test image
FROM builder as integration-test
WORKDIR /app
CMD ["go", "test", "./integration", "-v", "-tags=integration"]