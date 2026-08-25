package ui

import (
	"strings"
	"testing"
)

type row struct {
	Name string
	Age  int
}

func TestTextColumnInfersAndEscapes(t *testing.T) {
	c := Text("Nom", "name", func(r row) string { return r.Name })
	var sb strings.Builder
	if err := c.Cell(row{Name: `<script>x</script>`}).Render(t.Context(), &sb); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(sb.String(), "<script>") {
		t.Fatalf("cell did not escape: %q", sb.String())
	}

	n := Text("Âge", "", func(r row) int { return r.Age })
	sb.Reset()
	if err := n.Cell(row{Age: 42}).Render(t.Context(), &sb); err != nil {
		t.Fatal(err)
	}
	if sb.String() != "42" {
		t.Fatalf("int cell = %q", sb.String())
	}
}
