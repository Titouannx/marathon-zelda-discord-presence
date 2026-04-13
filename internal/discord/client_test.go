package discord

import "testing"

func TestNewRejectsEmptyApplicationID(t *testing.T) {
	t.Parallel()

	if _, err := New(""); err == nil {
		t.Fatal("expected error for empty Discord application id")
	}
}
