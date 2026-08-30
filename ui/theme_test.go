package ui

import (
	"os"
	"strings"
	"testing"
)

// daisyUI applique le thème sombre par trois chemins distincts : la préférence
// système, l'attribut [data-theme] et, pour ThemeToggle, un
// `:root:has(input.theme-controller[value=dark]:checked)` -- aucun attribut
// n'est posé sur la racine dans ce dernier cas. Le correctif de contraste sur
// --color-primary doit couvrir les trois : n'en couvrir que deux rendait le
// thème sombre sans le correctif dès qu'on basculait avec ThemeToggle,
// c'est-à-dire 4,13:1 au lieu du seuil AA de 4,5.
func TestDarkPrimaryFixCoversEveryDarkSelector(t *testing.T) {
	css, err := os.ReadFile("../assets/app.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, sel := range []string{
		`:root:not([data-theme="light"]):not(:has(input.theme-controller[value="light"]:checked))`,
		`:root[data-theme="dark"]`,
		`:root:has(input.theme-controller[value="dark"]:checked)`,
	} {
		if !strings.Contains(string(css), sel) {
			t.Errorf("assets/app.css ne corrige pas --color-primary pour %s", sel)
		}
	}
}
