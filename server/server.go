package main

import "log"
import "github.com/mundanelizard/koyi/server/auth/routers"

func main() {
	server := routers.SetUpServer()
	err := server.Run()
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}
