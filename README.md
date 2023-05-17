# File service

It is a microservice with CRUD access via API  
*(Actually there are 2 microservices, but catfacts was just study one and I don't want to delete it)*

## Commands

There are next commands:
- `/get?id=x` - sends you file via http
- `/set` - load file from you and response with id
- `/delete?id=x` - delete file by id

## Enviroment variables

There are 2 variables:
- `FILESERVICE_PATH` - defines where files will be saved
- `FILESERVICE_PORT` - defines what port will be used

## Useful code

It is covered with tests in `fileservice/fileservice_test.go` and has controller to send requests to service in `fileservice/controller.go`

## Run

You can launch it with `go run main.go` and it will use default settings  
or using docker  
`docker build -t catinapoke/go-microservice .`  
`docker run -l fileservice -p 3001:3001 --rm --env FILESERVICE_PATH=/go/data --mount t type=bind,src=./bin,dst=/go/data catinapoke/go-microservice`