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

func TestVisibleTrayMenuEntriesWhenNoSessionIncludeOpenPage(t *testing.T) {
	labels := visibleTrayMenuLabels(nil)

	if len(labels) != 3 {
		t.Fatalf("expected three visible tray actions without session, got %d", len(labels))
	}
	if labels[0] != "Ouvrir la page du marathon" || labels[1] != "Desinstaller" || labels[2] != "Quitter" {
		t.Fatalf("unexpected tray labels without session: %#v", labels)
	}
}

func TestVisibleTrayMenuEntriesWhenSessionActiveStayMinimal(t *testing.T) {
	current := &model.PresenceStatus{Active: true, GameName: "The Minish Cap"}

	labels := visibleTrayMenuLabels(current)

	if len(labels) != 2 {
		t.Fatalf("expected two visible tray actions with an active session, got %d", len(labels))
	}
	if labels[0] != "Desinstaller" || labels[1] != "Quitter" {
		t.Fatalf("unexpected tray labels with active session: %#v", labels)
	}
}
