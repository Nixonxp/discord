package main

import (
	config "github.com/Nixonxp/discord/channel/configs"
	"github.com/Nixonxp/discord/channel/internal/app/server"
	"log"
	"os"
)

func main() {
	cfg := config.GetConfig()

	app := server.NewApplication(cfg)
	if err := app.Run(); err != nil {
		log.Printf("server err: %v", err)
		os.Exit(1)
	}
}
