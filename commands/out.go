package commands

import (
	"crypto/tls"
	"fmt"
	"net/http"

	godiscord "github.com/bwmarrin/discordgo"

	resource "github.com/logsquaredn/discord-notification-resource"
)

// Out ...
func (r *DiscordNotificationResource) Out() error {
	var (
		req resource.OutRequest
		resp resource.OutResponse
	)

	err := r.readInput(&req)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	s, err := godiscord.New(req.Source.Token)
	if err != nil {
		return err
	}

	s.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify : true,
			},
		},
	}

	msg, err := s.WebhookExecute(
		req.Source.WebhookID,
		req.Source.Token,
		req.GetParams.Wait,
		&req.Params,
	)
	if err != nil {
		return err
	}

	if req.GetParams.Wait {
		resp.Version.Message = msg.ID
		resp.Metadata = r.getMetadata(msg)
	}

	r.writeOutput(&resp)
	if err != nil {
		return err
	}

	return nil
}
