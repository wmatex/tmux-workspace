package main

import (
	"fmt"
	"log"

	"github.com/wmatex/automux/internal/config"
)

func main() {
	var config config.Config
	err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %s\n", err)
	}
	fmt.Println(config)
	fmt.Println("Hello world")
}
