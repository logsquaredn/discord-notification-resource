package commands

import (
	"crypto/tls"
	"fmt"
	"net/http"

	discordgo "github.com/bwmarrin/discordgo"

	resource "github.com/logsquaredn/discord-notification-resource"
)

// Out runs the in script which checks stdin for a JSON object of the form of an OutRequest,
// sends a new message and then, if wait was set to true, fetches and writes it as well as Metadata about it to stdout
func (r *DiscordNotificationResource) Out() error {
	var (
		req  		  resource.OutRequest
		resp 		  resource.OutResponse
	)

	err := r.readInput(&req)
	if err != nil {
		return fmt.Errorf("could not marshal JSON: %s", err)
	}

	s, err := discordgo.New(req.Source.Token)
	if err != nil {
		return err
	}

	s.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	msg, err := s.WebhookExecute(
		req.Source.WebhookID,
		req.Source.Token,
		req.Params.Wait,
		&discordgo.WebhookParams{
			Content: req.Params.Content,
			Username: req.Params.Username,
			AvatarURL: req.Params.AvatarURL,
			TTS: req.Params.TTS,
			File: req.Params.File,
			Embeds: req.Params.Embeds,
			AllowedMentions: req.Params.AllowedMentions,
		},
	)
	if err != nil {
		return err
	}

	if req.Params.Wait {
		resp.Version.Message = msg.ID
		resp.Metadata, err = r.getMetadata(msg)
		if err != nil {
			return err
		}
	}

	r.writeOutput(&resp)
	if err != nil {
		return err
	}

	return nil
}
