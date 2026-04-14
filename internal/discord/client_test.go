package discord

import (
	"testing"
	"time"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
)

func TestNewRejectsEmptyApplicationID(t *testing.T) {
	t.Parallel()

	if _, err := New(""); err == nil {
		t.Fatal("expected error for empty Discord application id")
	}
}

func TestBuildActivityUsesRequestedCopyAndButton(t *testing.T) {
	t.Parallel()

	startedAt := "2026-04-14T18:30:00Z"
	status := model.PresenceStatus{
		Active:           true,
		GameName:         "The Minish Cap",
		GameLogoURL:      "https://loon.bzh/api/zelda/presence/logo/theminishcap",
		SessionStartedAt: startedAt,
		ProfileURL:       "https://loon.bzh/zelda/profile/p_demo",
	}

	activity := buildActivity(status)

	if activity.Details != "En train de jouer a The Minish Cap" {
		t.Fatalf("unexpected details: %q", activity.Details)
	}
	if activity.State != "via loon.bzh/zelda" {
		t.Fatalf("unexpected state: %q", activity.State)
	}
	if activity.LargeImage != status.GameLogoURL {
		t.Fatalf("unexpected large image: %q", activity.LargeImage)
	}
	if activity.LargeText != status.GameName {
		t.Fatalf("unexpected large text: %q", activity.LargeText)
	}
	if activity.Timestamps == nil || activity.Timestamps.Start == nil {
		t.Fatal("expected start timestamp to be set")
	}
	if got := activity.Timestamps.Start.UTC().Format(time.RFC3339); got != startedAt {
		t.Fatalf("unexpected timestamp: %q", got)
	}
	if len(activity.Buttons) != 1 {
		t.Fatalf("expected one button, got %d", len(activity.Buttons))
	}
	if activity.Buttons[0].Label != "Voir le profil" {
		t.Fatalf("unexpected button label: %q", activity.Buttons[0].Label)
	}
	if activity.Buttons[0].Url != status.ProfileURL {
		t.Fatalf("unexpected button url: %q", activity.Buttons[0].Url)
	}
}

func TestBuildActivityFallsBackToNowOnInvalidTimestamp(t *testing.T) {
	t.Parallel()

	before := time.Now().Add(-2 * time.Second)
	activity := buildActivity(model.PresenceStatus{
		Active:           true,
		GameName:         "Majora's Mask",
		GameLogoURL:      "https://loon.bzh/api/zelda/presence/logo/majorasmask",
		SessionStartedAt: "not-a-date",
		ProfileURL:       "https://loon.bzh/zelda/profile/p_demo",
	})
	after := time.Now().Add(2 * time.Second)

	if activity.Timestamps == nil || activity.Timestamps.Start == nil {
		t.Fatal("expected start timestamp to be set")
	}

	if activity.Timestamps.Start.Before(before) || activity.Timestamps.Start.After(after) {
		t.Fatalf("unexpected fallback timestamp: %s", activity.Timestamps.Start.UTC().Format(time.RFC3339))
	}
}

func TestBuildActivityParsesPostgresStyleTimestamp(t *testing.T) {
	t.Parallel()

	status := model.PresenceStatus{
		Active:           true,
		GameName:         "The Minish Cap",
		GameLogoURL:      "https://www.loon.bzh/api/zelda/presence/logo/theminishcap",
		SessionStartedAt: "2026-04-14 18:30:00.123+00",
		ProfileURL:       "https://www.loon.bzh/zelda/profile/p_demo",
	}

	activity := buildActivity(status)

	if activity.Timestamps == nil || activity.Timestamps.Start == nil {
		t.Fatal("expected start timestamp to be set")
	}

	if got := activity.Timestamps.Start.UTC().Format(time.RFC3339Nano); got != "2026-04-14T18:30:00.123Z" {
		t.Fatalf("unexpected parsed timestamp: %q", got)
	}
}
