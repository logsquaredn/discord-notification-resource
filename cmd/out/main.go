package main

import (
	"log"
	"os"

	"github.com/logsquaredn/discord-notification-resource/commands"
)

func main() {
	command := commands.NewDiscordNotificationResource(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	err := command.Out()
	if err != nil {
		log.Fatal(err)
	}
}
