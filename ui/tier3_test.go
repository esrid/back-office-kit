package ui

import (
	"strings"
	"testing"
)

func TestTier3Boundaries(t *testing.T) {
	if got := comboDelay(0); got != "250" {
		t.Fatalf("default combo delay = %q", got)
	}
	if got := jobProgress(-1); got != 0 {
		t.Fatalf("negative job progress = %d", got)
	}
	if got := jobProgress(120); got != 100 {
		t.Fatalf("overflow job progress = %d", got)
	}
	if got := wizardPosition(9, 3); got != 3 {
		t.Fatalf("wizard position = %d", got)
	}
}

func TestUnreadNotifications(t *testing.T) {
	items := []Notification{{Read: false}, {Read: true}, {Read: false}}
	if got := unreadNotifications(items); got != 2 {
		t.Fatalf("unread count = %d", got)
	}
}

func TestAsyncStateLoadingIsBusy(t *testing.T) {
	var out strings.Builder
	err := AsyncState(AsyncStateProps{ID: "users", Status: ResourceLoading}).Render(t.Context(), &out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `aria-busy="true"`) || !strings.Contains(out.String(), `role="status"`) {
		t.Fatalf("loading state lacks busy semantics: %q", out.String())
	}
}

func TestNotificationContentEscapes(t *testing.T) {
	var out strings.Builder
	err := NotificationCenter([]Notification{{Title: `<script>x</script>`, Href: "/"}}, "/notifications", false).Render(t.Context(), &out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "<script>") {
		t.Fatalf("notification title was not escaped: %q", out.String())
	}
}
