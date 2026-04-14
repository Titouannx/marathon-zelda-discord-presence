package app

import (
	"testing"

	"github.com/Titouannx/marathon-zelda-discord-presence/internal/model"
)

func TestResolveProfileTargetUsesProfileWhenAvailable(t *testing.T) {
	target := resolveProfileTarget("https://www.loon.bzh/zelda/profile/p_demo")

	if target != "https://www.loon.bzh/zelda/profile/p_demo" {
		t.Fatalf("expected profile url to be kept, got %q", target)
	}
}

func TestResolveProfileTargetFallsBackToSite(t *testing.T) {
	target := resolveProfileTarget("")

	if target != zeldaSiteURL {
		t.Fatalf("expected fallback site url %q, got %q", zeldaSiteURL, target)
	}
}

func TestShouldUpdatePresenceSkipsIdenticalPayload(t *testing.T) {
	current := model.PresenceStatus{
		Active:           true,
		GameID:           13,
		GameName:         "The Minish Cap",
		GameLogoURL:      "https://www.loon.bzh/api/zelda/presence/logo/theminishcap",
		SessionStartedAt: "2026-04-14T18:30:00.123Z",
		ProfileURL:       "https://www.loon.bzh/zelda/profile/p_demo",
	}

	if shouldUpdatePresence(&current, current) {
		t.Fatal("expected identical presence payload to be skipped")
	}
}

func TestShouldUpdatePresenceRunsOnChangedPayload(t *testing.T) {
	current := model.PresenceStatus{
		Active:           true,
		GameID:           13,
		GameName:         "The Minish Cap",
		GameLogoURL:      "https://www.loon.bzh/api/zelda/presence/logo/theminishcap",
		SessionStartedAt: "2026-04-14T18:30:00.123Z",
		ProfileURL:       "https://www.loon.bzh/zelda/profile/p_demo",
	}
	next := current
	next.GameName = "Ocarina of Time"

	if !shouldUpdatePresence(&current, next) {
		t.Fatal("expected changed presence payload to trigger an update")
	}
}

func TestVisibleTrayMenuEntriesAreMinimal(t *testing.T) {
	labels := visibleTrayMenuLabels()

	if len(labels) != 2 {
		t.Fatalf("expected two visible tray actions, got %d", len(labels))
	}
	if labels[0] != "Desinstaller" || labels[1] != "Quitter" {
		t.Fatalf("unexpected tray labels: %#v", labels)
	}
}
