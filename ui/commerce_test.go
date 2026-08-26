package ui

import (
	"strings"
	"testing"
)

func TestFormatAddressByCountry(t *testing.T) {
	base := Address{Name: "Amélie Rousseau", Line1: "12 rue des Lilas", City: "Lyon", PostalCode: "69003"}

	cases := []struct {
		iso  string
		want []string
	}{
		// Continental Europe: postal code before the city.
		{"FR", []string{"Amélie Rousseau", "12 rue des Lilas", "69003 Lyon"}},
		// Anglo layout: "city, region postal" -- with no region, the postal code
		// still follows the city rather than vanishing.
		{"US", []string{"Amélie Rousseau", "12 rue des Lilas", "Lyon, 69003"}},
		// An unknown country must not lose the address.
		{"", []string{"Amélie Rousseau", "12 rue des Lilas", "69003 Lyon"}},
	}
	for _, c := range cases {
		a := base
		a.CountryISO = c.iso
		got := formatAddress(a)
		if strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("%q: got %v, want %v", c.iso, got, c.want)
		}
	}

	us := formatAddress(Address{Name: "Jane Doe", Line1: "1 Market St", City: "San Francisco",
		Region: "CA", PostalCode: "94105", CountryISO: "US", Country: "United States"})
	if strings.Join(us, "|") != "Jane Doe|1 Market St|San Francisco, CA 94105|United States" {
		t.Errorf("US layout wrong: %v", us)
	}

	// The UK puts the postcode on its own line, after the city.
	gb := formatAddress(Address{Name: "John Smith", Line1: "10 Downing St", City: "London",
		PostalCode: "SW1A 2AA", CountryISO: "GB", Country: "United Kingdom"})
	if strings.Join(gb, "|") != "John Smith|10 Downing St|London|SW1A 2AA|United Kingdom" {
		t.Errorf("GB layout wrong: %v", gb)
	}
}

// Empty parts must be dropped, never printed as blank lines on a label.
func TestFormatAddressDropsEmptyParts(t *testing.T) {
	got := formatAddress(Address{Name: "Acme", Line1: "1 rue A", CountryISO: "FR"})
	for _, line := range got {
		if strings.TrimSpace(line) == "" {
			t.Fatalf("blank line in %v", got)
		}
	}
	if len(got) != 2 {
		t.Errorf("got %v, want 2 lines", got)
	}
}

func TestRefundableCents(t *testing.T) {
	cases := []struct{ total, already, want int64 }{
		{10000, 0, 10000},
		{10000, 2500, 7500},
		{10000, 10000, 0},
		{10000, 12000, 0}, // the server over-reported; never offer a negative refund
		{0, 0, 0},
	}
	for _, c := range cases {
		if got := refundableCents(c.total, c.already); got != c.want {
			t.Errorf("refundableCents(%d,%d) = %d, want %d", c.total, c.already, got, c.want)
		}
	}
}

// A fully refunded order must not present a refund control at all.
func TestRefundFormClosesWhenNothingIsLeft(t *testing.T) {
	done := renderTier4(t, RefundForm(RefundFormProps{
		OrderID: "o1", Currency: "EUR", TotalCents: 10000, AlreadyRefundedCents: 10000,
		Action: "/refund",
	}))
	if strings.Contains(done, `name="amount"`) {
		t.Errorf("no amount field should be offered: %s", done)
	}
	if !strings.Contains(done, "intégralement remboursée") {
		t.Errorf("the reason should be stated: %s", done)
	}

	open := renderTier4(t, RefundForm(RefundFormProps{
		OrderID: "o1", Currency: "EUR", TotalCents: 10000, AlreadyRefundedCents: 2500,
		Action: "/refund",
	}))
	if !strings.Contains(open, `name="amount"`) {
		t.Errorf("a partially refunded order must still allow a refund: %s", open)
	}
	// The remaining amount, not the total, is what is proposed and stated.
	if !strings.Contains(open, `value="75.00"`) {
		t.Errorf("the field should default to what is left: %s", open)
	}
	if !strings.Contains(open, "Au maximum 75,00"+nbsp+"€") {
		t.Errorf("the cap should be spelled out: %s", open)
	}
}
