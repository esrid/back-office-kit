package marketing

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func render(t *testing.T, c templ.Component) string {
	t.Helper()
	var out strings.Builder
	if err := c.Render(context.Background(), &out); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

// One <h1> per page, and it belongs to Hero. Everything else opens at <h2>,
// or the heading outline lies to search engines and screen readers.
func TestHeadingHierarchy(t *testing.T) {
	hero := render(t, Hero(HeroProps{Title: "Pilotez vos opérations"}))
	if strings.Count(hero, "<h1") != 1 {
		t.Errorf("Hero must render exactly one h1: %s", hero)
	}

	others := map[string]templ.Component{
		"FeatureGrid":     FeatureGrid("Fonctionnalités", "", []Feature{{Title: "A"}}),
		"PricingTable":    PricingTable("Tarifs", "", []Plan{{Name: "Pro"}}),
		"TestimonialGrid": TestimonialGrid("Ils l'utilisent", []Quote{{Text: "Bien", Author: "X"}}),
		"FAQ":             FAQ("Questions", []QA{{Question: "Q", Answer: "R"}}),
		"CTASection":      CTASection("Commencer", "", nil, ""),
	}
	for name, c := range others {
		html := render(t, c)
		if strings.Contains(html, "<h1") {
			t.Errorf("%s must not emit an h1: %s", name, html)
		}
		if !strings.Contains(html, "<h2") {
			t.Errorf("%s should open at h2: %s", name, html)
		}
	}
}

// An image with no alt text is dropped rather than shipped: it would fail for
// screen readers and say nothing to a crawler.
func TestPictureRequiresAltAndDimensions(t *testing.T) {
	withAlt := render(t, Picture(Image{Src: "/a.png", Alt: "Tableau de bord", Width: 1200, Height: 800}, ""))
	for _, want := range []string{`alt="Tableau de bord"`, `width="1200"`, `height="800"`, `loading="lazy"`} {
		if !strings.Contains(withAlt, want) {
			t.Errorf("missing %q: %s", want, withAlt)
		}
	}

	if got := render(t, Picture(Image{Src: "/a.png"}, "")); got != "" {
		t.Errorf("an image without alt must not render: %q", got)
	}

	// A zero dimension must be omitted, never written as width="0".
	noDims := render(t, Picture(Image{Src: "/a.png", Alt: "x"}, ""))
	if strings.Contains(noDims, `width="0"`) || strings.Contains(noDims, `height="0"`) {
		t.Errorf("zero dimensions must be omitted: %s", noDims)
	}
}

// The FAQ must work with JavaScript disabled.
func TestFAQIsScriptFree(t *testing.T) {
	html := render(t, FAQ("Questions", []QA{{Question: "Combien ?", Answer: "29 €"}}))
	if strings.Contains(html, "<script") || strings.Contains(html, "onclick") {
		t.Errorf("the FAQ must not need JavaScript: %s", html)
	}
	if !strings.Contains(html, "<details") || !strings.Contains(html, "<summary") {
		t.Errorf("expected a details/summary accordion: %s", html)
	}
	if !strings.Contains(html, `name="faq"`) {
		t.Errorf("a shared name makes the accordion exclusive where supported: %s", html)
	}
}

// A quote is data: the attribution must be machine-readable, not styled text.
func TestTestimonialUsesSemanticMarkup(t *testing.T) {
	html := render(t, TestimonialGrid("", []Quote{{Text: "Excellent", Author: "Amélie", Role: "Ops"}}))
	for _, want := range []string{"<blockquote", "<figcaption", "<cite"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s: %s", want, html)
		}
	}
}

// SiteHeader must not steal the page's only h1: it sits above the Hero, and a
// second h1 breaks the outline without changing anything visible.
func TestSiteHeaderCarriesNoHeading(t *testing.T) {
	html := render(t, SiteHeader(SiteHeaderProps{
		Product: "Acme",
		Links:   []Action{{Label: "Tarifs", Href: "/pricing"}},
		Actions: []Action{{Label: "Essayer", Href: "/signup", Primary: true}},
		Current: "/pricing",
	}))
	for _, tag := range []string{"<h1", "<h2", "<h3"} {
		if strings.Contains(html, tag) {
			t.Errorf("SiteHeader must emit no heading, found %s: %s", tag, html)
		}
	}
	if !strings.Contains(html, `aria-current="page"`) {
		t.Errorf("the current section must be marked, not merely coloured: %s", html)
	}
}
