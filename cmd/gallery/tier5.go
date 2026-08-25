package main

import "example.com/components/ui"

func demoApprovalRequests() []ui.ApprovalRequest {
	return []ui.ApprovalRequest{
		{
			ID: "approval_81", PlanID: "plan_7b31", PlanVersion: "3",
			Title: "Suspendre 12 comptes et envoyer les notifications",
			Summary: "La politique backend exige une approbation humaine pour les effets externes.",
			Requester: "Assistant opérations", RequestedAt: "il y a 4 min", RequestedDatetime: "2026-08-25T18:42:00-04:00",
			ExpiresAt: "dans 6 min", ExpiresDatetime: "2026-08-25T18:52:00-04:00",
			Risk: ui.AgentRiskSensitive, Status: ui.ApprovalPending, ReviewHref: "/approvals/81",
		},
		{
			ID: "approval_79", PlanID: "plan_a204", PlanVersion: "1",
			Title: "Corriger les catégories de 43 factures",
			Summary: "Modification réversible approuvée par Nora Martin.",
			Requester: "Assistant facturation", RequestedAt: "aujourd’hui à 16:08", RequestedDatetime: "2026-08-25T16:08:00-04:00",
			Risk: ui.AgentRiskWrite, Status: ui.ApprovalApproved, ReviewHref: "/approvals/79",
		},
		{
			ID: "approval_72", PlanID: "plan_119c", PlanVersion: "4",
			Title: "Exporter les coordonnées des prospects",
			Summary: "La demande n’a pas été examinée avant son expiration.",
			Requester: "Assistant commercial", RequestedAt: "hier à 11:20", RequestedDatetime: "2026-08-24T11:20:00-04:00",
			Risk: ui.AgentRiskSensitive, Status: ui.ApprovalExpired, ReviewHref: "/approvals/72",
		},
	}
}

func demoChangeItems() []ui.ChangeDiffItem {
	return []ui.ChangeDiffItem{
		{Label: "Statut", Before: "Actif", After: "Suspendu", Kind: ui.ChangeUpdated},
		{Label: "Limite mensuelle", Before: "20 000 €", After: "10 000 €", Kind: ui.ChangeUpdated},
		{Label: "Motif interne", Before: "—", After: "Facture échue depuis 37 jours", Kind: ui.ChangeAdded},
		{Label: "Exception temporaire", Before: "Valable jusqu’au 31 août", After: "—", Kind: ui.ChangeRemoved},
	}
}

func tier5Sections() []Section {
	return []Section{
		{
			Group: "Tier 5 — Gouvernance", ID: "approval-inbox", Title: "ApprovalInbox",
			Purpose: "Centraliser les demandes d’approbation durables et leur état backend.",
			Uses: []string{"approval_id", "plan version", "up-poll"}, Demo: demoApprovalInbox(),
			Snippet: `@ui.ApprovalInbox(ui.ApprovalInboxProps{ID: "approvals", Items: requests, Toolbar: filters})`,
			Note: "Chaque ligne transporte l’identifiant d’approbation et la version exacte du plan. Examiner ouvre la surface dédiée ; la liste n’approuve rien elle-même.",
		},
		{
			Group: "Tier 5 — Gouvernance", ID: "change-diff", Title: "ChangeDiff",
			Purpose: "Comparer les valeurs avant/après dans un plan, un conflit ou un reçu.",
			Uses: []string{"table", "del/ins", "redacted values"}, Demo: demoChangeDiff(),
			Snippet: `@ui.ChangeDiff(ui.ChangeDiffProps{ID: "changes", Title: "Modifications", Items: changes})`,
			Note: "Le backend filtre et masque les champs sensibles avant rendu. Le composant ne doit jamais recevoir une valeur que l’utilisateur n’est pas autorisé à consulter.",
		},
		{
			Group: "Tier 5 — Gouvernance", ID: "conflict-resolver", Title: "ConflictResolver",
			Purpose: "Résoudre une écriture devenue obsolète sans écraser silencieusement une modification concurrente.",
			Uses: []string{"conflict_id", "versions", "strategy"}, Demo: demoConflictResolver(),
			Snippet: `@ui.ConflictResolver(props, ui.ChangeDiff(diff))`,
			Note: "Le backend décide quelles stratégies sont permises puis vérifie de nouveau conflict_id et les deux versions au moment du POST.",
		},
		{
			Group: "Tier 5 — Gouvernance", ID: "access-notice", Title: "AccessNotice",
			Purpose: "Expliquer une autorisation, un refus ou une vérification supplémentaire décidés côté serveur.",
			Uses: []string{"decision_id", "policy", "step-up"}, Demo: demoAccessNotices(),
			Snippet: `@ui.AccessNotice(ui.AccessNoticeProps{Decision: ui.AccessDenied, DecisionID: decision.ID, ...})`,
			Note: "AccessNotice explique une décision ; il ne protège aucune ressource. La route et le backend doivent rester inaccessibles quand la décision est refusée.",
		},
	}
}
