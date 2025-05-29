## OmniChart-Server
The server backend for OmniChart.

## Setup
Golang 1.24.3
Setup the .env file.

## Run
`go run cmd/server/main.go`

## Generate Documentation
`swag init -g ./cmd/server/main.go`
check swaggerUI on the `/swagger/index.html` endpoint
