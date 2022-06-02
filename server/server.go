package main

import "log"

func main() {
	server := SetUpServer()
	err := server.Run()
	if err != nil {
		log.Fatalf("Finally, there's an error %s", err)
	}
}
