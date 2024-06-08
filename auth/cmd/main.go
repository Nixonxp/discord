package main

import (
	config "github.com/Nixonxp/discord/auth/configs"
	"github.com/Nixonxp/discord/auth/internal/app/server"
	"github.com/rs/zerolog/log"
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
