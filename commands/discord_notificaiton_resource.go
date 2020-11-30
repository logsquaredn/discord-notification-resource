package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	godiscord "github.com/bwmarrin/discordgo"
	resource "github.com/logsquaredn/discord-notification-resource"
)

// DiscordNotificationResource ...
type DiscordNotificationResource struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

// NewDiscordNotificationResource ...
func NewDiscordNotificationResource(
	stdin  io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args   []string,
) *DiscordNotificationResource {
	return &DiscordNotificationResource{
		stdin,
		stderr,
		stdout,
		args,
	}
}

func (r *DiscordNotificationResource) readInput(req *resource.OutRequest) error {
	decoder := json.NewDecoder(r.stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	return nil
}

func (r *DiscordNotificationResource) writeOutput(resp interface{}) error {
	err := json.NewEncoder(r.stdout).Encode(resp)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	return nil
}


func (r *DiscordNotificationResource) getMetadata(msg *godiscord.Message) []resource.Metadata {
	if msg != nil {
		return []resource.Metadata{

		}
	}

	return []resource.Metadata{}
}

func (r *DiscordNotificationResource) expandEnv(s string) string {
	return os.Expand(s, func(v string) string {
		switch v {
		case "BUILD_ID", "BUILD_NAME", "BUILD_JOB_NAME", "BUILD_PIPELINE_NAME", "BUILD_TEAM_NAME", "ATC_EXTERNAL_URL":
			return os.Getenv(v)
		}
		return "$" + v
	})
}
