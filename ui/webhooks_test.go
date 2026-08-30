package ui

import "testing"

// The tone decides whether an operator retries or goes fixing the other end,
// so the boundaries are the whole point: 4xx is the remote refusing on
// purpose -- replaying the same payload gets the same refusal -- while 3xx
// and 5xx are failures a replay can legitimately resolve.
func TestWebhookTone(t *testing.T) {
	cases := map[int]string{
		200: "success", 201: "success", 299: "success",
		300: "error", 301: "error",
		400: "warning", 404: "warning", 422: "warning", 499: "warning",
		500: "error", 502: "error",
		0: "error", 199: "error",
	}
	for status, want := range cases {
		if got := webhookTone(status); got != want {
			t.Errorf("webhookTone(%d) = %q, want %q", status, got, want)
		}
	}
}

// A delivery nothing ever answered must not print "0": it reads like a status
// code and is not one.
func TestWebhookStatusLabel(t *testing.T) {
	if got := webhookStatusLabel(0); got != "Sans réponse" {
		t.Errorf("webhookStatusLabel(0) = %q", got)
	}
	if got := webhookStatusLabel(503); got != "503" {
		t.Errorf("webhookStatusLabel(503) = %q", got)
	}
}

func TestCurrentOrgNameFallsBackToTheFirst(t *testing.T) {
	orgs := []Org{{ID: "acme", Name: "Acme"}, {ID: "globex", Name: "Globex"}}
	if got := currentOrgName(OrgSwitcherProps{Orgs: orgs, CurrentID: "globex"}); got != "Globex" {
		t.Errorf("got %q, want Globex", got)
	}
	// An unknown tenant must not render an empty trigger.
	if got := currentOrgName(OrgSwitcherProps{Orgs: orgs, CurrentID: "inconnu"}); got != "Acme" {
		t.Errorf("unknown id: got %q, want Acme", got)
	}
	if got := currentOrgName(OrgSwitcherProps{}); got != "" {
		t.Errorf("no org: got %q", got)
	}
}
