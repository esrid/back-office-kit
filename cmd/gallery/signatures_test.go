package main

import (
	"strings"
	"testing"
)

func TestSignaturesAreReadFromSource(t *testing.T) {
	sigs := signatures("../../ui")
	if len(sigs) < 50 {
		t.Fatalf("only %d signatures found; the ui package should expose far more", len(sigs))
	}

	cases := map[string]string{
		"Shell":     "func Shell(title string, nav []NavItem, current string, user string, flashes []Flash, headerExtras ...templ.Component) templ.Component",
		"Timestamp": "func Timestamp(p TimestampProps) templ.Component",
	}
	for name, want := range cases {
		if got := sigs[name]; got != want {
			t.Errorf("%s:\n got %q\nwant %q", name, got, want)
		}
	}

	// Generics must survive: a signature without [T any] would mislead.
	if got := sigs["DataTable"]; !strings.Contains(got, "[T any]") {
		t.Errorf("DataTable lost its type parameter: %q", got)
	}
}

func TestSignatureForMatchesCompoundTitles(t *testing.T) {
	sigs := signatures("../../ui")
	if got := signatureFor(sigs, "Field · SelectField"); !strings.HasPrefix(got, "func Field(") {
		t.Errorf("compound title should resolve to its first component, got %q", got)
	}
	if got := signatureFor(sigs, "Rien de connu"); got != "" {
		t.Errorf("unknown title should yield nothing, got %q", got)
	}
}
