package discord

import (
	"errors"
	"strings"
	"time"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
	rich "github.com/jrh3k5/rich-go/client"
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

	if err := rich.SetActivity(buildActivity(status)); err != nil {
		var closed *rich.ErrClosedConnection
		if errors.As(err, &closed) {
			c.connected = false
		}
		return err
	}

	return nil
}

func (c *Client) Clear() error {
	if !c.connected {
		return nil
	}

	err := rich.Logout()
	c.connected = false
	return err
}

func (c *Client) Close() {
	_ = c.Clear()
}

func buildActivity(status model.PresenceStatus) rich.Activity {
	start, err := time.Parse(time.RFC3339, status.SessionStartedAt)
	if err != nil {
		start = time.Now()
	}

	return rich.Activity{
		Details:    "En train de jouer a " + status.GameName,
		State:      "via loon.bzh/zelda",
		LargeImage: status.GameLogoURL,
		LargeText:  status.GameName,
		Timestamps: &rich.Timestamps{Start: &start},
		Buttons: []*rich.Button{
			{
				Label: "Voir le profil",
				Url:   status.ProfileURL,
			},
		},
	}
}
