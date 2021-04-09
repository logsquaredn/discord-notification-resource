package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	discordgo "github.com/bwmarrin/discordgo"
	resource "github.com/logsquaredn/discord-notification-resource"
)

// DiscordNotificationResource struct which has the Check, In, and Out methods on it which comprise
// the three scripts needed to implement a Concourse Resource Type
type DiscordNotificationResource struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

// NewDiscordNotificationResource creates a new DiscordNotificationResource struct
func NewDiscordNotificationResource(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
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

func (r *DiscordNotificationResource) getMetadata(msg *discordgo.Message) ([]resource.Metadata, error) {
	if msg != nil {
		timestamp, err := msg.Timestamp.Parse()
		if err != nil {
			return nil, err
		}

		return []resource.Metadata{
			{Name: "channelId", Value: msg.ChannelID},
			{Name: "content", Value: msg.Content},
			{Name: "guildId", Value: msg.GuildID},
			{Name: "webhookId", Value: msg.WebhookID},
			{Name: "timestamp", Value: timestamp.String()},
		}, nil
	}

	return nil, fmt.Errorf("returned message was nil")
}

func (r *DiscordNotificationResource) getSrc() (string, error) {
	if len(r.args) < 2 {
		return "", fmt.Errorf("destination path not specified")
	}
	return r.args[1], nil
}

func (r *DiscordNotificationResource) getContent(p *resource.PutParams) (string, error) {
	if p.ContentFile != "" {
		src, err := r.getSrc()
		if err != nil {
			return "", err
		}

		b, err := ioutil.ReadFile(filepath.Join(src, p.ContentFile))
		if err != nil {
			return "", err
		}

		return r.expandEnv(string(b)), nil
	} else if p.Content != "" {
		return r.expandEnv(p.Content), nil
	}

	return "", nil
}

func (r *DiscordNotificationResource) getAvatarURL(p *resource.PutParams) (string, error) {
	if p.AvatarURLFile != "" {
		src, err := r.getSrc()
		if err != nil {
			return "", err
		}

		b, err := ioutil.ReadFile(filepath.Join(src, p.AvatarURLFile))
		if err != nil {
			return "", err
		}

		return r.expandEnv(string(b)), nil
	} else if p.AvatarURL != "" {
		return r.expandEnv(p.AvatarURL), nil
	}

	return "", nil
}

func (r *DiscordNotificationResource) getEmbeds(p *resource.PutParams) ([]*discordgo.MessageEmbed, error) {
	for _, e := range p.Embeds {
		e.URL = r.expandEnv(e.URL)
		e.Title = r.expandEnv(e.Title)
		e.Description = r.expandEnv(e.Description)
		e.Image.URL = r.expandEnv(e.Image.URL)
		e.Image.ProxyURL = r.expandEnv(e.Image.ProxyURL)
		e.Thumbnail.URL = r.expandEnv(e.Thumbnail.URL)
		e.Thumbnail.ProxyURL = r.expandEnv(e.Thumbnail.ProxyURL)
		e.Footer.Text = r.expandEnv(e.Footer.Text)
		e.Footer.IconURL = r.expandEnv(e.Footer.IconURL)
		e.Footer.ProxyIconURL = r.expandEnv(e.Footer.ProxyIconURL)
		e.Provider.Name = r.expandEnv(e.Provider.Name)
		e.Provider.URL = r.expandEnv(e.Provider.URL)
		e.Author.IconURL = r.expandEnv(e.Author.IconURL)
		e.Author.URL = r.expandEnv(e.Author.URL)
		e.Author.Name = r.expandEnv(e.Author.Name)
		e.Author.ProxyIconURL = r.expandEnv(e.Author.ProxyIconURL)
		for _, f := range e.Fields {
			f.Name = r.expandEnv(f.Name)
			f.Value = r.expandEnv(f.Value)
		}
	}

	return p.Embeds, nil
}

func (r *DiscordNotificationResource) getUsername(p *resource.PutParams) (string, error) {
	if p.UsernameFile != "" {
		src, err := r.getSrc()
		if err != nil {
			return "", err
		}

		b, err := ioutil.ReadFile(filepath.Join(src, p.UsernameFile))
		if err != nil {
			return "", err
		}

		return r.expandEnv(string(b)), nil
	} else if p.Username != "" {
		return r.expandEnv(p.Username), nil
	}

	return "", nil
}

func (r *DiscordNotificationResource) writeMetadata(mds []resource.Metadata) error {
	src, err := r.getSrc()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(src, ".metadata"), 0755)
	if err != nil {
		return fmt.Errorf("unable to make directory %s", src)
	}

	for _, md := range mds {
		err = ioutil.WriteFile(filepath.Join(src, ".metadata", md.Name), []byte(md.Value), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *DiscordNotificationResource) expandEnv(s string) string {
	if s == "" {
		return s
	}

	return os.Expand(s, func(v string) string {
		switch v {
		case "BUILD_ID", "BUILD_URL", "BUILD_NAME", "BUILD_JOB_NAME", "BUILD_PIPELINE_NAME", "BUILD_TEAM_NAME", "ATC_EXTERNAL_URL":
			return os.Getenv(v)
		}
		return "$" + v
	})
}
