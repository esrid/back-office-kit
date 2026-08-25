package ui

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		offset time.Duration
		want   string
	}{
		{0, "à l'instant"},
		{-30 * time.Second, "à l'instant"},
		{-60 * time.Second, "il y a 1 minute"},
		{-5 * time.Minute, "il y a 5 minutes"},
		{-90 * time.Minute, "il y a 2 heures"},
		{-5 * time.Hour, "il y a 5 heures"},
		{3 * time.Hour, "dans 3 heures"},
		{48 * time.Hour, "dans 2 jours"},
	}
	for _, c := range cases {
		if got := relativeTime(now.Add(c.offset), now); got != c.want {
			t.Errorf("relativeTime(%v) = %q, want %q", c.offset, got, c.want)
		}
	}
}

// The number and its unit must never disagree, at any distance.
func TestRelativeTimeAgreesWithItsNumber(t *testing.T) {
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	for seconds := 1; seconds < 400*24*3600; seconds = seconds*3/2 + 1 {
		phrase := relativeTime(now.Add(-time.Duration(seconds)*time.Second), now)
		if phrase == "à l'instant" {
			continue
		}
		fields := strings.Fields(phrase)
		unit := fields[len(fields)-1]
		n, err := strconv.Atoi(fields[len(fields)-2])
		if err != nil {
			t.Fatalf("unparseable phrase %q", phrase)
		}
		if n < 1 {
			t.Errorf("%q: non-positive count", phrase)
		}
		// "mois" is invariable, so it is exempt from the -s rule.
		if unit == "mois" {
			continue
		}
		if isPlural := strings.HasSuffix(unit, "s"); (n > 1) != isPlural {
			t.Errorf("%q: number and unit disagree", phrase)
		}
	}
}

func TestAbsoluteTimeAlwaysNamesTheZone(t *testing.T) {
	paris, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	at := time.Date(2026, 8, 25, 17, 13, 0, 0, time.UTC)

	utc, local := absoluteTime(at, time.UTC), absoluteTime(at, paris)
	if utc == local {
		t.Fatal("the same instant must not read identically in two zones")
	}
	if !strings.Contains(utc, "UTC") || !strings.Contains(local, "CEST") {
		t.Errorf("zone missing: %q / %q", utc, local)
	}
	if !strings.Contains(local, "19:13") {
		t.Errorf("Paris is UTC+2 in August, expected 19:13: %q", local)
	}
	if !strings.Contains(local, "25 août 2026") {
		t.Errorf("malformed date: %q", local)
	}
}

func TestTimestampDefaultsToUTCNotLocalMachine(t *testing.T) {
	at := time.Date(2026, 1, 5, 9, 0, 0, 0, time.UTC)
	html := renderTier4(t, Timestamp(TimestampProps{At: at, Now: at, Absolute: true}))
	if !strings.Contains(html, "UTC") {
		t.Errorf("a missing Location must fall back to UTC explicitly: %s", html)
	}
	if !strings.Contains(html, `datetime="2026-01-05T09:00:00Z"`) {
		t.Errorf("machine-readable datetime missing: %s", html)
	}
}
