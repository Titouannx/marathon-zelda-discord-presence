package discord

import (
	"errors"
	"strings"
	"time"

	rich "github.com/jrh3k5/rich-go/client"
	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
)

type Client struct {
	appID     string
	connected bool
}

func New(appID string) (*Client, error) {
	if strings.TrimSpace(appID) == "" {
		return nil, errors.New("missing Discord application id")
	}

	return &Client{appID: appID}, nil
}

func (c *Client) ensureConnected() error {
	if c.connected {
		return nil
	}

	if err := rich.Login(c.appID); err != nil {
		return err
	}

	c.connected = true
	return nil
}

func (c *Client) Set(status model.PresenceStatus) error {
	if err := c.ensureConnected(); err != nil {
		return err
	}

	start, err := time.Parse(time.RFC3339, status.SessionStartedAt)
	if err != nil {
		start = time.Now()
	}

	return rich.SetActivity(rich.Activity{
		Details:    status.GameName,
		State:      "Marathon en cours - via loon.bzh/zelda",
		LargeImage: status.GameLogoURL,
		LargeText:  status.GameName,
		Timestamps: &rich.Timestamps{Start: &start},
		Buttons: []*rich.Button{
			{
				Label: "Voir le profil",
				Url:   status.ProfileURL,
			},
		},
	})
}

func (c *Client) Clear() error {
	if !c.connected {
		return nil
	}

	return rich.SetActivity(rich.Activity{})
}

func (c *Client) Close() {
	if !c.connected {
		return
	}

	rich.Logout()
	c.connected = false
}
