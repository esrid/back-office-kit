package ui

import (
	"strings"
	"testing"
)

var importFields = []ImportField{
	{Name: "email", Label: "E-mail", Required: true},
	{Name: "name", Label: "Nom", Required: true},
	{Name: "role", Label: "Rôle"},
}

func names(fields []ImportField) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		out[i] = f.Name
	}
	return out
}

func TestUnmappedRequired(t *testing.T) {
	cases := []struct {
		name    string
		columns []ImportColumn
		want    []string
	}{
		{"rien associé", nil, []string{"email", "name"}},
		{"tout associé", []ImportColumn{
			{Source: "Email address", Target: "email"},
			{Source: "Full name", Target: "name"},
		}, nil},
		{"optionnel seul ne suffit pas", []ImportColumn{
			{Source: "Role", Target: "role"},
		}, []string{"email", "name"}},
		{"colonne ignorée ne compte pas", []ImportColumn{
			{Source: "Email address", Target: ""},
			{Source: "Full name", Target: "name"},
		}, []string{"email"}},
	}
	for _, c := range cases {
		got := names(unmappedRequired(importFields, c.columns))
		if strings.Join(got, ",") != strings.Join(c.want, ",") {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestDuplicateTargets(t *testing.T) {
	cols := []ImportColumn{
		{Source: "Email", Target: "email"},
		{Source: "E-mail pro", Target: "email"},
		{Source: "Nom", Target: "name"},
		{Source: "Ignorée", Target: ""},
		{Source: "Aussi ignorée", Target: ""},
	}
	if got := names(duplicateTargets(importFields, cols)); strings.Join(got, ",") != "email" {
		t.Errorf("got %v, want [email]", got)
	}
	// Two ignored columns are not a duplicate of each other.
	if got := duplicateTargets(importFields, []ImportColumn{{Target: ""}, {Target: ""}}); len(got) != 0 {
		t.Errorf("ignored columns must not collide: %v", names(got))
	}
}

// The submit button must be unavailable while a required field is unmapped:
// the warning alone is too easy to scroll past.
func TestImportCannotStartWhileRequiredFieldsAreUnmapped(t *testing.T) {
	blocked := renderTier4(t, ImportMapper(ImportMapperProps{
		ID: "imp", Action: "/import", Fields: importFields,
		Columns: []ImportColumn{{Source: "Nom", Target: "name"}},
	}))
	if !strings.Contains(blocked, "disabled") {
		t.Errorf("import should be blocked while E-mail is unmapped: %s", blocked)
	}
	if !strings.Contains(blocked, "obligatoire") {
		t.Errorf("the reason should be stated, not just the block: %s", blocked)
	}

	ok := renderTier4(t, ImportMapper(ImportMapperProps{
		ID: "imp", Action: "/import", Fields: importFields,
		Columns: []ImportColumn{{Source: "Nom", Target: "name"}, {Source: "Mail", Target: "email"}},
	}))
	i := strings.Index(ok, "Lancer l'import")
	start := strings.LastIndex(ok[:i], "<button")
	if strings.Contains(ok[start:i], "disabled") {
		t.Errorf("import should be available once every required field is mapped: %s", ok[start:i])
	}
}

func TestErrorPageQuotesItsRequestID(t *testing.T) {
	html := renderTier4(t, ErrorPage(ErrorPageProps{
		Code: "500", Title: "Quelque chose a cassé", RequestID: "req_8f31c2",
	}))
	if !strings.Contains(html, "req_8f31c2") {
		t.Errorf("the operator needs something precise to quote: %s", html)
	}
	if !strings.Contains(html, toneOf("error").Badge) {
		t.Errorf("a 500 should read as an error, got: %s", html)
	}
	if !strings.Contains(renderTier4(t, ErrorPage(ErrorPageProps{Code: "404"})), toneOf("neutral").Badge) {
		t.Error("a 404 is not a system failure and should stay neutral")
	}
}

func TestBannersRenderNothingWhenInactive(t *testing.T) {
	if html := renderTier4(t, ImpersonationBanner(ImpersonationProps{})); html != "" {
		t.Errorf("no impersonation, no banner, got %q", html)
	}
	if html := renderTier4(t, UndoSnackbar(UndoProps{})); html != "" {
		t.Errorf("nothing to undo, no snackbar, got %q", html)
	}
}
