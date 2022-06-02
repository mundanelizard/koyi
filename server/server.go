package main

import (
	"github.com/mundanelizard/koyi/server/routers"
	"log"
)

func main() {
	server := routers.SetUpServer()
	err := server.Run()
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}
