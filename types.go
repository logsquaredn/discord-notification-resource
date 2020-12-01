package resource

import (
	discordgo "github.com/bwmarrin/discordgo"
)

// Version is the JSON object that is passed to and from Concourse
type Version struct {
	Message string `json:"message"`
}

// OutRequest is the JSON object that Concourse passes to /opt/resource/out through stdin
type OutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

// PutParams ...
type PutParams struct {
	discordgo.WebhookParams
	Wait bool `json:"wait,omitempty"`
}

// OutResponse is the JSON object that we pass back to Concourse through stdout from /opt/resource/out
type OutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// Source is the JSON (yaml) object configured under the resources array in a Concourse pipeline
type Source struct {
	WebhookID string `json:"webhook_id"`
	Token     string `json:"token"`
}

// Metadata is the object which is passed in array form to Concourse through stdout from /opt/resource/out and /opt/resource/in
// to provide additional information about the Version
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
