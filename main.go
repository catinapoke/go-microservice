package main

import (
	"fmt"
	"log"
	"os"

	"github.com/catinapoke/go-microservice/fileservice"
)

func GetPort() string {
	port := os.Getenv("FILESERVICE_PORT")
	if port == "" {
		port = "3001"
	}

	return port
}

func GetPath() string {
	path := os.Getenv("FILESERVICE_PATH")
	if path == "" {
		path = "./bin"
	}

	return path
}

func RunFileServer() {
	path := GetPath()
	service := fileservice.CreateFileService(path)
	service = fileservice.CreateFileServiceLogger(service)

	apiServer := fileservice.CreateAPIServer(service)

	port := GetPort()
	fmt.Printf("Starting File service at port %s and storage at '%s'\n", port, path)
	log.Fatal(apiServer.Start(":" + port))
}

func main() {
	RunFileServer()
}
