package ui

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func renderTier4(t *testing.T, component templ.Component) string {
	t.Helper()
	var out strings.Builder
	if err := component.Render(t.Context(), &out); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestApprovalBindsExactPlan(t *testing.T) {
	html := renderTier4(t, ApprovalCard(ApprovalCardProps{
		PlanID: "plan_42", PlanVersion: "7", PlanDigest: "sha256:abc",
		Title: "Suspendre les comptes", Risk: AgentRiskSensitive,
		Action: "/plans/42/approve", ConfirmLabel: "Approuver", RequireText: "CONFIRMER",
	}, nil))

	for _, want := range []string{
		`name="plan_id" value="plan_42"`,
		`name="plan_version" value="7"`,
		`name="plan_digest" value="sha256:abc"`,
		`name="confirmation" autocomplete="off" required`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("approval missing %q: %s", want, html)
		}
	}
}

func TestBusyAgentThreadPolls(t *testing.T) {
	html := renderTier4(t, AgentThread(AgentThreadProps{
		ID: "agent", Title: "Assistant", Busy: true, PollSource: "/agent/42",
	}))
	if !strings.Contains(html, `aria-busy="true"`) || !strings.Contains(html, `up-source="/agent/42"`) {
		t.Fatalf("busy thread lacks polling semantics: %s", html)
	}
	if !strings.Contains(html, `role="log"`) || strings.Contains(html, `<ol class="flex flex-col gap-3 p-4" role="log"`) {
		t.Fatalf("thread must preserve list semantics inside its live log: %s", html)
	}
}

func TestAgentComposerUsesProvidedID(t *testing.T) {
	html := renderTier4(t, AgentComposer(AgentComposerProps{
		ID: "invoice-agent-prompt", Action: "/agent/messages", Name: "prompt",
	}))
	if !strings.Contains(html, `for="invoice-agent-prompt"`) || !strings.Contains(html, `id="invoice-agent-prompt"`) {
		t.Fatalf("composer label is not bound to its configured field: %s", html)
	}
}

func TestUndoReferencesReceipt(t *testing.T) {
	html := renderTier4(t, AgentUndoAction(AgentUndoProps{
		ReceiptID: "receipt_99", Action: "/receipts/99/undo", Available: true,
	}))
	if !strings.Contains(html, `name="receipt_id" value="receipt_99"`) {
		t.Fatalf("undo is not bound to its receipt: %s", html)
	}
}

func TestAgentReceiptEscapesBackendText(t *testing.T) {
	html := renderTier4(t, AgentResultReceipt(AgentReceiptProps{
		ReceiptID: "r1", PlanID: "p1", Version: "1", Title: `<script>x</script>`, Status: AgentToolSuccess,
	}))
	if strings.Contains(html, "<script>") {
		t.Fatalf("receipt title was not escaped: %s", html)
	}
}

func TestRiskLevelsRemainDistinct(t *testing.T) {
	if agentRiskTone(AgentRiskRead) == agentRiskTone(AgentRiskWrite) || agentRiskTone(AgentRiskWrite) == agentRiskTone(AgentRiskSensitive) {
		t.Fatal("agent risk levels must have distinct semantic tones")
	}
}

// [up-layer] accepts three stacking values and nothing else. Folding the
// overlay mode in — "new modal", "new drawer" — parses as an unknown value and
// the overlay silently never opens as asked.
func TestLayerStackingRejectsFoldedModes(t *testing.T) {
	for _, valid := range []string{"new", "swap", "shatter"} {
		if got := layerStacking(valid); got != valid {
			t.Errorf("layerStacking(%q) = %q", valid, got)
		}
	}
	for _, invalid := range []string{"new modal", "new drawer", "modal", "", "popup", "NEW"} {
		if got := layerStacking(invalid); got != "" {
			t.Errorf("layerStacking(%q) = %q, want empty", invalid, got)
		}
	}

	html := renderTier4(t, AccessNotice(AccessNoticeProps{
		Decision: AccessDenied, Title: "Accès refusé",
		ActionHref: "/request", ActionLayer: "new modal",
	}))
	if strings.Contains(html, "new modal") {
		t.Errorf("a folded mode must never reach the markup: %s", html)
	}
}

// No component may emit a folded [up-layer] value. Comment lines are skipped:
// the fix in slideover.templ is documented by quoting the broken form, and a
// guard that fires on its own explanation is a guard nobody keeps.
func TestNoComponentFoldsTheLayerMode(t *testing.T) {
	files, err := filepath.Glob("*.templ")
	if err != nil {
		t.Fatal(err)
	}
	folded := regexp.MustCompile(`up-layer="[a-z]+ `)

	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for n, line := range strings.Split(string(body), "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "//") {
				continue
			}
			if folded.MatchString(line) {
				t.Errorf("%s:%d folds the overlay mode into [up-layer]; it belongs in [up-mode]:\n\t%s",
					file, n+1, strings.TrimSpace(line))
			}
		}
	}
}
