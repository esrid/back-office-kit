package ui

import (
	"strings"
	"testing"
)

func TestFileFieldClampsProgress(t *testing.T) {
	if got := (FileFieldProps{Progress: -4}).progress(); got != 0 {
		t.Fatalf("negative progress = %d", got)
	}
	if got := (FileFieldProps{Progress: 140}).progress(); got != 100 {
		t.Fatalf("overflow progress = %d", got)
	}
}

func TestBulkActionBarHiddenWithoutSelection(t *testing.T) {
	var out strings.Builder
	if err := BulkActionBar(0, "élément sélectionné", "éléments sélectionnés").Render(t.Context(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("zero-selection bulk bar rendered %q", out.String())
	}
}

func TestPlural(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{{0, "élément"}, {1, "élément"}, {2, "éléments"}, {142, "éléments"}}
	for _, c := range cases {
		if got := plural(c.n, "élément", "éléments"); got != c.want {
			t.Errorf("plural(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestSlideOverAcceptedEscapesTarget(t *testing.T) {
	got := slideOverAccepted(`#users"bad`)
	if !strings.Contains(got, `\"`) {
		t.Fatalf("handler target was not quoted: %q", got)
	}
}
