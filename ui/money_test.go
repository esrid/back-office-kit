package ui

import (
	"math"
	"strings"
	"testing"
)

// nbsp is U+00A0: FormatMoney uses it both between thousands and before
// the currency symbol, so "€" can never wrap onto a line of its own.
const nbsp = "\u00a0"

func TestFormatMoney(t *testing.T) {
	cases := []struct {
		cents    int64
		currency string
		want     string
	}{
		{0, "EUR", "0,00" + nbsp + "€"},
		{5, "EUR", "0,05" + nbsp + "€"},
		{50, "EUR", "0,50" + nbsp + "€"},
		{100, "EUR", "1,00" + nbsp + "€"},
		{1250, "EUR", "12,50" + nbsp + "€"},
		{128450, "EUR", "1" + nbsp + "284,50" + nbsp + "€"},
		{123456789, "EUR", "1" + nbsp + "234" + nbsp + "567,89" + nbsp + "€"},
		{-4500, "EUR", "−45,00" + nbsp + "€"},
		{-5, "EUR", "−0,05" + nbsp + "€"},
		{1250, "USD", "12,50" + nbsp + "$"},
		// A currency with no minor unit must not gain fake decimals.
		{1000, "JPY", "1" + nbsp + "000" + nbsp + "¥"},
		// Three-decimal currencies exist.
		{1234, "KWD", "1,234" + nbsp + "KWD"},
		// Unknown currency falls back to its code rather than guessing.
		{999, "XYZ", "9,99" + nbsp + "XYZ"},
	}
	for _, c := range cases {
		if got := FormatMoney(c.cents, c.currency); got != c.want {
			t.Errorf("FormatMoney(%d, %q) = %q, want %q", c.cents, c.currency, got, c.want)
		}
	}
}

// The function must be total: no input may panic or produce a malformed amount.
func TestFormatMoneyIsTotal(t *testing.T) {
	values := []int64{math.MinInt64, math.MaxInt64, -1, 0, 1, 9, 10, 99, 100}
	for _, v := range values {
		for _, cur := range []string{"EUR", "JPY", "KWD", ""} {
			got := FormatMoney(v, cur)
			if got == "" {
				t.Errorf("FormatMoney(%d, %q) returned nothing", v, cur)
			}
			if strings.Contains(got, "-") {
				t.Errorf("FormatMoney(%d, %q) = %q: uses a hyphen instead of U+2212", v, cur, got)
			}
		}
	}
}

// Grouping must never leave a separator at either edge.
func TestGroupThousandsEdges(t *testing.T) {
	cases := map[string]string{
		"1": "1", "12": "12", "123": "123",
		"1234": "1" + nbsp + "234", "1234567": "1" + nbsp + "234" + nbsp + "567",
	}
	for in, want := range cases {
		got := groupThousands(in)
		if got != want {
			t.Errorf("groupThousands(%q) = %q, want %q", in, got, want)
		}
		if strings.HasPrefix(got, nbsp) || strings.HasSuffix(got, nbsp) {
			t.Errorf("groupThousands(%q) = %q: separator on an edge", in, got)
		}
	}
}

func TestMoneySignOnlyOnPositive(t *testing.T) {
	plus := renderTier4(t, Money(MoneyProps{Cents: 1250, Currency: "EUR", Sign: true}))
	if !strings.Contains(plus, "+12,50") {
		t.Errorf("expected a leading +: %s", plus)
	}
	minus := renderTier4(t, Money(MoneyProps{Cents: -1250, Currency: "EUR", Sign: true}))
	if strings.Contains(minus, "+") {
		t.Errorf("a negative amount must never gain a +: %s", minus)
	}
	zero := renderTier4(t, Money(MoneyProps{Cents: 0, Currency: "EUR", Sign: true}))
	if strings.Contains(zero, "+") {
		t.Errorf("zero is not positive: %s", zero)
	}
}

func TestLineItemsSum(t *testing.T) {
	p := LineItemsProps{
		Currency: "EUR",
		Items: []LineItem{
			{Label: "Carnet", Quantity: 2, UnitCents: 1200, TotalCents: 2400},
			{Label: "Stylo", Quantity: 1, UnitCents: 4500, TotalCents: 4500},
		},
		Adjustments: []Adjustment{
			{Label: "Remise", Cents: -690},
			{Label: "Livraison", Cents: 590},
		},
	}
	if got := lineItemsSum(p); got != 6800 {
		t.Fatalf("lineItemsSum = %d, want 6800", got)
	}
}

// A total that contradicts its own lines must be shown, never smoothed over.
func TestLineItemsShowsTheDiscrepancy(t *testing.T) {
	p := LineItemsProps{
		Currency:   "EUR",
		Items:      []LineItem{{Label: "Carnet", Quantity: 1, UnitCents: 1200, TotalCents: 1200}},
		TotalCents: 9900, // the server disagrees with its own lines
	}
	html := renderTier4(t, LineItemsTable(p))
	if !strings.Contains(html, "alert-warning") {
		t.Errorf("a mismatched total must raise a visible warning: %s", html)
	}
	if !strings.Contains(html, "Écart") {
		t.Errorf("the discrepancy should be named: %s", html)
	}
	// The server's total still wins on screen; the component never invents one.
	if !strings.Contains(html, "99,00"+nbsp+"€") {
		t.Errorf("the displayed total must remain the server's: %s", html)
	}

	p.TotalCents = 1200
	if ok := renderTier4(t, LineItemsTable(p)); strings.Contains(ok, "alert-warning") {
		t.Errorf("a consistent total must not warn: %s", ok)
	}
}
