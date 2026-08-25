package ui

import (
	"net/url"
	"testing"
)

var demoFilters = []Filter{
	{Name: "q", Label: "Recherche"},
	{Name: "status", Label: "Statut", Options: []Option{{Value: "active", Label: "Actif"}}},
	{Name: "role", Label: "Rôle", Options: []Option{{Value: "admin", Label: "Admin"}}},
}

func TestHasFilters(t *testing.T) {
	cases := []struct {
		q    url.Values
		want bool
	}{
		{url.Values{}, false},
		{url.Values{"sort": {"email"}, "page": {"3"}}, false}, // sort is not a filter
		{url.Values{"q": {"bob"}}, true},
		{url.Values{"status": {"active"}}, true},
		{url.Values{"q": {""}}, false}, // present but empty is not filtering
	}
	for _, c := range cases {
		if got := HasFilters(c.q, demoFilters); got != c.want {
			t.Errorf("HasFilters(%v) = %v, want %v", c.q, got, c.want)
		}
	}
}

func TestResetHrefKeepsSortDropsFiltersAndPage(t *testing.T) {
	q := url.Values{
		"q": {"bob"}, "status": {"active"}, "role": {"admin"},
		"sort": {"email"}, "dir": {"desc"}, "page": {"4"},
	}
	if got, want := ResetHref(q, demoFilters), "?dir=desc&sort=email"; got != want {
		t.Errorf("ResetHref = %q, want %q", got, want)
	}
	if q.Get("q") != "bob" || q.Get("page") != "4" {
		t.Error("ResetHref mutated the caller's url.Values")
	}

	if got, want := ResetHref(url.Values{"q": {"bob"}}, demoFilters), "?"; got != want {
		t.Errorf("ResetHref with nothing left = %q, want %q", got, want)
	}
}

func TestResetHrefIsIdempotent(t *testing.T) {
	q := url.Values{"q": {"bob"}, "sort": {"email"}}
	first := ResetHref(q, demoFilters)
	parsed, err := url.ParseQuery(first[1:])
	if err != nil {
		t.Fatal(err)
	}
	if second := ResetHref(parsed, demoFilters); second != first {
		t.Errorf("ResetHref not idempotent: %q then %q", first, second)
	}
}

func TestDateRangeCountsAndResetsBothParameters(t *testing.T) {
	filters := []Filter{{
		Kind: FilterDateRange, Name: "from", EndName: "to", Label: "Période",
	}}
	q := url.Values{
		"from": {"2026-08-01"}, "to": {"2026-08-25"},
		"sort": {"created_at"}, "page": {"3"},
	}

	if !HasFilters(q, filters) {
		t.Fatal("date range was not detected as an active filter")
	}
	if got, want := ResetHref(q, filters), "?sort=created_at"; got != want {
		t.Fatalf("ResetHref(date range) = %q, want %q", got, want)
	}
	if q.Get("from") == "" || q.Get("to") == "" {
		t.Fatal("ResetHref mutated the caller's date range")
	}
}
