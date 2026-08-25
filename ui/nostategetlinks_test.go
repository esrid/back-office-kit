package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// anchorHrefs returns every href that appears on an <a> element.
var anchorHref = regexp.MustCompile(`<a\b[^>]*\bhref="([^"]*)"`)

func anchorHrefs(html string) []string {
	var out []string
	for _, m := range anchorHref.FindAllStringSubmatch(html, -1) {
		out = append(out, m[1])
	}
	return out
}

// State changes must never be reachable through a link. http.CrossOriginProtection
// always allows GET, HEAD and OPTIONS, so a state-changing <a> gets no CSRF
// protection at all -- and is additionally exposed to prefetching and to
// [up-preload], which fires without showing [up-confirm].
func TestStateChangingActionsAreNotLinks(t *testing.T) {
	const destructive = "/DESTRUCTIVE"

	cases := map[string]templ.Component{
		"ApprovalCard.RejectHref": ApprovalCard(ApprovalCardProps{
			PlanID: "p1", Action: "/approve", RejectHref: destructive, RequireText: "OUI",
		}, nil),
		"AgentComposer.RemoveHref": AgentComposer(AgentComposerProps{
			ID: "c", Action: "/send", Name: "prompt",
			Contexts: []AgentContext{{Label: "Vue", Value: "x", RemoveHref: destructive}},
		}),
		"RepeaterField.RemoveHref": RepeaterField("rows", "Lignes",
			[]RepeaterRow{{ID: "r1", RemoveHref: destructive}}, "/add"),
		"RepeaterField.AddHref":   RepeaterField("rows", "Lignes", nil, destructive),
		"FormActions.destroyHref": FormActions("Enregistrer", destructive),
		"RowActionMenu.Method": RowActionMenu("m1", "Actions", []Action{
			{Label: "Supprimer", Href: destructive, Method: "delete", Confirm: "Sûr ?"},
		}),
	}

	for name, component := range cases {
		html := renderTier4(t, component)
		for _, href := range anchorHrefs(html) {
			if href == destructive {
				t.Errorf("%s: state change is an <a href> and rides on GET:\n%s", name, html)
			}
		}
		if !strings.Contains(html, destructive) {
			t.Errorf("%s: destructive target vanished from the output entirely", name)
		}
	}
}

// A read-only action stays a plain link: turning navigation into a form would
// break middle-click, bookmarking and the back button for no security gain.
func TestReadOnlyRowActionStaysALink(t *testing.T) {
	html := renderTier4(t, RowActionMenu("m1", "Actions", []Action{
		{Label: "Voir la fiche", Href: "/users/42"},
	}))
	if !strings.Contains(html, `<a href="/users/42"`) {
		t.Fatalf("navigation action should remain a link: %s", html)
	}
}

// Secondary submits must skip constraint validation, or a required field meant
// to gate the primary action would also block cancelling or rejecting.
func TestSecondarySubmitsSkipValidation(t *testing.T) {
	html := renderTier4(t, ApprovalCard(ApprovalCardProps{
		PlanID: "p1", Action: "/approve", RejectHref: "/reject", RequireText: "CONFIRMER",
	}, nil))
	if !strings.Contains(html, "required") {
		t.Fatal("test premise broken: RequireText should render a required input")
	}
	i := strings.Index(html, `formaction="/reject"`)
	if i == -1 {
		t.Fatalf("reject button missing: %s", html)
	}
	end := strings.Index(html[i:], ">") + i
	if !strings.Contains(html[i:end], "formnovalidate") {
		t.Errorf("reject would be blocked by the confirmation field: %s", html[i:end])
	}
}
