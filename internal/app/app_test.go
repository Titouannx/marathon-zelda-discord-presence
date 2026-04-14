package app

import "testing"

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
