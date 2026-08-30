package ui

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// renderHTML renders a component to a string for the assertions below.
func renderHTML(t *testing.T, component templ.Component) string {
	t.Helper()
	var out strings.Builder
	if err := component.Render(t.Context(), &out); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

// A <textarea value="..."> renders empty: the value belongs in the content.
func TestTextareaFieldCarriesValueInContent(t *testing.T) {
	html := renderHTML(t, TextareaField(FieldProps{Name: "note", Label: "Note", Value: "deux\nlignes"}))
	if strings.Contains(html, `value=`) {
		t.Errorf("value moved into an attribute: %s", html)
	}
	if !strings.Contains(html, ">deux\nlignes</textarea>") {
		t.Errorf("content lost: %s", html)
	}
}

func TestBoxFieldsSubmitAndCheck(t *testing.T) {
	cases := []struct {
		name  string
		props FieldProps
		want  []string
		deny  []string
	}{
		{"unchecked", FieldProps{Name: "tos", Label: "J'accepte"}, []string{`value="1"`}, []string{"checked"}},
		{"checked", FieldProps{Name: "tos", Label: "J'accepte", Checked: true}, []string{"checked", `value="1"`}, nil},
		{"custom value", FieldProps{Name: "tos", Value: "yes", Checked: true}, []string{`value="yes"`}, []string{`value="1"`}},
		{"error", FieldProps{Name: "tos", Error: "Obligatoire"}, []string{"checkbox-error", `aria-invalid="true"`, "Obligatoire"}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			html := renderHTML(t, CheckboxField(c.props))
			for _, want := range c.want {
				if !strings.Contains(html, want) {
					t.Errorf("missing %q: %s", want, html)
				}
			}
			for _, deny := range c.deny {
				if strings.Contains(html, deny) {
					t.Errorf("unexpected %q: %s", deny, html)
				}
			}
		})
	}
}

// The colour classes are written out in full, so a toggle must not borrow the
// checkbox ones -- and neither may be concatenated from the control name.
func TestToggleFieldKeepsItsOwnClasses(t *testing.T) {
	html := renderHTML(t, ToggleField(FieldProps{Name: "beta", Label: "Bêta", Error: "Non disponible"}))
	if !strings.Contains(html, `class="toggle toggle-error"`) || strings.Contains(html, `class="checkbox`) {
		t.Errorf("toggle rendered with the wrong control classes: %s", html)
	}
}

func TestRadioGroupChecksOnlyTheCurrentValue(t *testing.T) {
	options := []Option{{Value: "admin", Label: "Administrateur"}, {Value: "member", Label: "Membre"}}
	html := renderHTML(t, RadioGroup(FieldProps{Name: "role", Label: "Rôle", Value: "member", Validate: true}, options))

	if n := strings.Count(html, "checked"); n != 1 {
		t.Errorf("%d radios checked, want 1: %s", n, html)
	}
	if !strings.Contains(html, `id="f-role-member" name="role" type="radio" value="member" class="radio" checked`) {
		t.Errorf("the checked radio is not the current value: %s", html)
	}
	// Distinct ids keep [up-validate] resolving to fieldset:has(#f-role-admin).
	if !strings.Contains(html, `id="f-role-admin"`) || !strings.Contains(html, `for="f-role-member"`) {
		t.Errorf("per-option ids missing: %s", html)
	}
	if n := strings.Count(html, "<fieldset"); n != 1 {
		t.Errorf("the group is one fieldset, got %d: %s", n, html)
	}
}
