package main

import (
	"log"
)

func main() {
	service := NewCatFactService("https://catfact.ninja/fact")
	service = NewLoggingService(service)

	apiServer := NewApiServer(service)

	log.Fatal(apiServer.Start(":3000"))
}
