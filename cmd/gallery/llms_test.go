package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// llms.txt is generated, not written: a hand-maintained list of 90 components
// drifts on the first rename, and a signature that lies is worse than none.
func TestLLMsCoversEverySection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "llms.txt")
	// sections() reads the ui sources by the paths the generator uses, and it
	// runs from the repo root.
	t.Chdir("../..")
	secs := sections()
	if err := writeLLMs(path, secs); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out := string(body)

	if n := strings.Count(out, "\n- **"); n != len(secs) {
		t.Errorf("%d entries for %d sections", n, len(secs))
	}
	// The invariants matter more than the list: an agent that misses them
	// writes code that compiles and does nothing.
	for _, want := range []string{
		"Aucun changement d'état ne passe par un GET",
		"up-history",
		"en toutes lettres",
		"go get github.com/esrid/back-office-kit",
		"go list -m -f",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("llms.txt lost %q", want)
		}
	}
	// Signatures come from the source, so a missing one means the parser broke.
	if !strings.Contains(out, "func Pagination(current int, totalItems int, perPage int, q url.Values, target string) templ.Component") {
		t.Error("signatures are missing from llms.txt")
	}
}
