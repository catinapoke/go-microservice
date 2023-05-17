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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckFileDestination(path string) {
	if b, _ := exists(path); !b {
		os.MkdirAll(path, 0774)
	}
}

func RunFileServer() {
	path := GetPath()
	CheckFileDestination(path)
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
