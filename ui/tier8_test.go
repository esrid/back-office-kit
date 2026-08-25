package ui

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestPermissionMatrixPreservesLockedGrantedPermissions(t *testing.T) {
	html := renderTier4(t, PermissionMatrix(PermissionMatrixProps{
		ID: "permissions", Action: "/roles/ops", RoleID: "role_ops", RoleVersion: "17",
		Columns: []PermissionColumn{{Key: "read", Label: "Lire"}},
		Rows: []PermissionRow{{Resource: "accounts", Label: "Comptes", Cells: []PermissionCell{{
			Permission: "accounts:read", Granted: true, Inherited: true, Origin: "Membre",
		}}}},
	}))

	if !strings.Contains(html, `method="post"`) || !strings.Contains(html, `name="role_version" value="17"`) {
		t.Fatalf("permission changes must be versioned POSTs: %s", html)
	}
	if strings.Count(html, `value="accounts:read"`) != 2 {
		t.Fatalf("a disabled checked permission also needs a hidden successful control: %s", html)
	}
	if !strings.Contains(html, `disabled`) || !strings.Contains(html, "Héritée") {
		t.Fatalf("an inherited permission must be visibly locked: %s", html)
	}
}

func TestTier8MutationsAreNeverLinks(t *testing.T) {
	const destructive = "/DESTRUCTIVE"
	cases := map[string]templ.Component{
		"grants":                 AccessGrantList(AccessGrantListProps{ID: "grants", Action: destructive, Version: "1", Items: []AccessGrant{{ID: "g1", Subject: "A", Revocable: true}}}),
		"sessions":               SessionManager(SessionManagerProps{ID: "sessions", Action: destructive, Version: "1", Sessions: []SecuritySession{{ID: "s1", Device: "Phone", Revocable: true}}}),
		"keys":                   APIKeyManager(APIKeyManagerProps{ID: "keys", Action: destructive, Version: "1", Keys: []APIKey{{ID: "k1", Name: "Key", Status: APIKeyActive, Revocable: true}}}),
		"secret acknowledgement": SecretReveal(SecretRevealProps{ID: "secret", Secret: "secret-value", SecretID: "x1", AcknowledgeAction: destructive}),
		"step-up":                StepUpAuthCard(StepUpAuthProps{ID: "step", Action: destructive, ChallengeID: "c1", Method: StepUpMethod{Value: "totp", Label: "TOTP"}}),
	}

	for name, component := range cases {
		html := renderTier4(t, component)
		if !strings.Contains(html, `method="post"`) {
			t.Errorf("%s: mutation must use POST: %s", name, html)
		}
		for _, href := range anchorHrefs(html) {
			if href == destructive {
				t.Errorf("%s: mutation leaked into an href: %s", name, html)
			}
		}
	}
}

func TestCurrentSessionCannotBeRevokedIndividually(t *testing.T) {
	html := renderTier4(t, SessionManager(SessionManagerProps{
		ID: "sessions", Action: "/sessions", Sessions: []SecuritySession{
			{ID: "current", Device: "This browser", Current: true, Revocable: true},
			{ID: "other", Device: "Phone", Revocable: true},
		},
	}))
	if strings.Contains(html, `value="revoke:current"`) {
		t.Fatalf("the current session must not expose an individual revoke control: %s", html)
	}
	if !strings.Contains(html, `value="revoke:other"`) {
		t.Fatalf("another revocable session should expose its exact ID: %s", html)
	}
}

func TestAPIKeyManagerOnlyRendersMaskedPrefixes(t *testing.T) {
	html := renderTier4(t, APIKeyManager(APIKeyManagerProps{
		ID: "keys", Action: "/keys", Keys: []APIKey{{ID: "k1", Name: "Prod", Prefix: "bok_live_91f2", Revocable: true}},
	}))
	if !strings.Contains(html, "bok_live_91f2••••••••") {
		t.Fatalf("the recognizable masked prefix should be visible: %s", html)
	}
	if strings.Contains(html, "plaintext-secret") {
		t.Fatal("test premise broken: APIKeyManager must have no field for a plaintext secret")
	}
}

func TestSecretRevealIsOneTimeAndAcknowledgedByID(t *testing.T) {
	if html := renderTier4(t, SecretReveal(SecretRevealProps{ID: "secret"})); html != "" {
		t.Fatalf("no secret should render no handoff surface, got %q", html)
	}

	html := renderTier4(t, SecretReveal(SecretRevealProps{
		ID: "secret", Secret: "bok_live_secret_91f2", SecretID: "secret_91f2", AcknowledgeAction: "/ack",
	}))
	if !strings.Contains(html, `value="bok_live_secret_91f2"`) || !strings.Contains(html, `name="secret_id" value="secret_91f2"`) {
		t.Fatalf("the one-time value and its opaque acknowledgement ID are required: %s", html)
	}
	if strings.Contains(html, `data-secret`) {
		t.Fatalf("plaintext secrets must not be duplicated into technical attributes: %s", html)
	}
}

func TestStepUpCarriesBackendChallengeIdentity(t *testing.T) {
	html := renderTier4(t, StepUpAuthCard(StepUpAuthProps{
		ID: "step", Action: "/verify", ChallengeID: "challenge_42", ChallengeVersion: "3",
		Method: StepUpMethod{Value: "totp", Label: "Application", Autocomplete: "one-time-code"},
	}))
	for _, want := range []string{`name="challenge_id" value="challenge_42"`, `name="challenge_version" value="3"`, `name="method" value="totp"`, `autocomplete="one-time-code"`} {
		if !strings.Contains(html, want) {
			t.Errorf("step-up challenge missing %q: %s", want, html)
		}
	}
}

func TestPolicySimulatorIsReadOnlyGET(t *testing.T) {
	html := renderTier4(t, PolicySimulator(PolicySimulatorProps{
		ID: "policy", Action: "/simulate", Actors: []Option{{Value: "u1", Label: "User"}},
		Resources: []Option{{Value: "r1", Label: "Record"}}, Actions: []Option{{Value: "delete", Label: "Delete"}},
		Simulation: &PolicySimulation{DecisionID: "d1", Decision: AccessDenied},
	}))
	if !strings.Contains(html, `method="get"`) || strings.Contains(html, `method="post"`) {
		t.Fatalf("policy evaluation must stay a read-only GET: %s", html)
	}
	if !strings.Contains(html, "aucun accès n’a été accordé") || !strings.Contains(html, `data-decision-id="d1"`) {
		t.Fatalf("the result must state its non-authoritative nature and decision ID: %s", html)
	}
}

func TestSecurityEventFeedKeepsAuditReferences(t *testing.T) {
	html := renderTier4(t, SecurityEventFeed(SecurityEventFeedProps{
		ID: "events", Events: []SecurityEvent{{ID: "evt_1", Title: "Login", RequestID: "req_1", Href: "/events/1"}},
	}))
	for _, want := range []string{`data-event-id="evt_1"`, "req_1", `<a href="/events/1"`} {
		if !strings.Contains(html, want) {
			t.Errorf("security event missing %q: %s", want, html)
		}
	}
}
