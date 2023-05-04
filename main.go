package main

import (
	"fmt"
	"log"

	"github.com/catinapoke/go-microservice/catfacts"
	"github.com/catinapoke/go-microservice/fileservice"
)

func RunCatFactsServer() {
	service := catfacts.NewCatFactService("https://catfact.ninja/fact")
	service = catfacts.NewLoggingService(service)

	apiServer := catfacts.NewApiServer(service)

	fmt.Println("Starting Cat fact service at port 3001")
	log.Fatal(apiServer.Start(":3001"))
}

func RunFileServer() {
	service := fileservice.CreateFileService()
	service = fileservice.CreateFileServiceLogger(service)

	apiServer := fileservice.CreateAPIServer(service)

	fmt.Println("Starting File service at port 3001")
	log.Fatal(apiServer.Start(":3001"))
}

func main() {
	RunFileServer()
}
