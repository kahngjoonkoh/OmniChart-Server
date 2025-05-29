# OmniChart-Server
The server backend for OmniChart.

## Setup
Golang 1.24.3
Setup the .env file.

## Run
`go run main.go`

## Generate Documentation
`swag init -g main.go`
check swaggerUI on the `/swagger/index.html` endpoint

## Deployment
Check if environment variables are set on `tsuru env get -a omnichart-server`

[https://omnichart-server.impaas.uk/swagger/index.html](https://omnichart-server.impaas.uk/swagger/index.html)
