package ui

import (
	"strings"
	"testing"
)

func TestApprovalInboxPreservesBackendIdentity(t *testing.T) {
	html := renderTier4(t, ApprovalInbox(ApprovalInboxProps{
		ID: "approvals", Busy: true, PollSource: "/approvals",
		Items: []ApprovalRequest{{
			ID: "approval_1", PlanID: "plan_42", PlanVersion: "7",
			Title: "Suspendre les comptes", Status: ApprovalPending,
		}},
	}))

	for _, want := range []string{
		`data-approval-id="approval_1"`,
		`data-plan-id="plan_42"`,
		`data-plan-version="7"`,
		`up-source="/approvals"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("approval inbox missing %q: %s", want, html)
		}
	}
}

func TestChangeDiffEscapesBackendValues(t *testing.T) {
	html := renderTier4(t, ChangeDiff(ChangeDiffProps{
		ID: "diff", Title: "Changements",
		Items: []ChangeDiffItem{{Label: "Nom", Before: "Acme", After: `<script>x</script>`, Kind: ChangeUpdated}},
	}))
	if strings.Contains(html, "<script>") || !strings.Contains(html, "<del") || !strings.Contains(html, "<ins") {
		t.Fatalf("change diff must escape values and preserve diff semantics: %s", html)
	}
}

func TestConflictResolverBindsExactVersions(t *testing.T) {
	html := renderTier4(t, ConflictResolver(ConflictResolverProps{
		ID: "conflict", ConflictID: "conflict_9", SubmittedVersion: "12", CurrentVersion: "14",
		Title: "Conflit", Action: "/conflicts/9", AllowMerge: true,
	}, nil))

	for _, want := range []string{
		`name="conflict_id" value="conflict_9"`,
		`name="submitted_version" value="12"`,
		`name="current_version" value="14"`,
		`name="strategy" value="merge"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("conflict resolver missing %q: %s", want, html)
		}
	}
	if strings.Contains(html, `value="overwrite"`) {
		t.Fatalf("overwrite must stay hidden when the backend disallows it: %s", html)
	}
}

func TestAccessNoticeExposesBackendDecision(t *testing.T) {
	html := renderTier4(t, AccessNotice(AccessNoticeProps{
		ID: "access", DecisionID: "decision_17", Decision: AccessDenied,
		Reason: "Rôle insuffisant", Policy: "billing.write",
	}))
	for _, want := range []string{
		`data-access-decision="denied"`,
		`data-decision-id="decision_17"`,
		"Rôle insuffisant",
		"billing.write",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("access notice missing %q: %s", want, html)
		}
	}
}
