package main

import "example.com/components/ui"

func demoAgentPlanSteps() []ui.AgentPlanStep {
	return []ui.AgentPlanStep{
		{Title: "Identifier les comptes concernés", Description: "Rechercher les comptes dont la facture est échue depuis plus de 30 jours.", Tool: "billing.search", Target: "comptes actifs", Risk: ui.AgentRiskRead},
		{Title: "Suspendre 12 comptes", Description: "Bloquer les nouvelles connexions tout en conservant les données.", Tool: "accounts.suspend", Target: "12 comptes", Risk: ui.AgentRiskWrite},
		{Title: "Informer les propriétaires", Description: "Envoyer un e-mail individuel avec le motif et le lien de régularisation.", Tool: "mail.send", Target: "12 destinataires", Risk: ui.AgentRiskSensitive},
	}
}

func demoAgentImpacts() []ui.AgentImpactItem {
	return []ui.AgentImpactItem{
		{Label: "Données consultées", Value: "200 comptes et factures", Tone: "info"},
		{Label: "Modifications", Value: "12 suspensions réversibles", Tone: "warning"},
		{Label: "Effet externe", Value: "12 e-mails envoyés", Tone: "error"},
		{Label: "Coût estimé", Value: "0,08 €", Tone: "neutral"},
	}
}

func demoAgentToolRuns() []ui.AgentToolRun {
	return []ui.AgentToolRun{
		{ID: "run_1", Tool: "billing.search", Summary: "Recherche des factures échues", Status: ui.AgentToolSuccess, Started: "10:14:02", Duration: "840 ms", Detail: "12 comptes correspondent aux règles du backend."},
		{ID: "run_2", Tool: "accounts.suspend", Summary: "Suspension des comptes", Status: ui.AgentToolSuccess, Started: "10:14:03", Duration: "1,2 s", Detail: "12 modifications confirmées par la base."},
		{ID: "run_3", Tool: "mail.send", Summary: "Envoi des notifications", Status: ui.AgentToolRunning, Started: "10:14:05", Detail: "8 e-mails envoyés sur 12."},
	}
}

func demoAgentReceipt() ui.AgentReceiptProps {
	return ui.AgentReceiptProps{
		ReceiptID: "receipt_91f2", PlanID: "plan_7b31", Version: "3", Title: "Suspension terminée", Summary: "12 comptes ont été suspendus et leurs propriétaires informés.", Completed: "à 10:14:09", Status: ui.AgentToolSuccess,
		Changes:   []ui.AgentReceiptChange{{Label: "Comptes actifs", Before: "1 284", After: "1 272"}, {Label: "Comptes suspendus", Before: "31", After: "43"}},
		Artifacts: []ui.AgentReceiptLink{{Label: "Télécharger le rapport CSV", Href: "/receipts/91f2/report.csv"}},
		Sources:   []ui.AgentReceiptLink{{Label: "Factures échues", Href: "/billing?status=overdue"}, {Label: "Politique de suspension", Href: "/policies/suspension"}},
	}
}

func tier4Sections() []Section {
	return []Section{
		{Group: "Tier 4 — Agentic UI", ID: "agent-workspace", Title: "AgentWorkspace", Purpose: "Composer conversation et état d'action dans un espace responsive.", Uses: []string{"conversation", "sidebar", "responsive"}, Demo: demoAgentWorkspace(), Snippet: `@ui.AgentWorkspace("Assistant opérations", subtitle, planSidebar()) { @ui.AgentThread(...) }`, Note: "Sur écran étroit, le plan descend sous la conversation ; aucune information d'approbation n'est cachée dans une drawer transitoire."},
		{Group: "Tier 4 — Agentic UI", ID: "agent-thread", Title: "AgentThread · AgentMessage", Purpose: "Conversation annoncée comme journal vivant, avec auteur et horodatage explicites.", Uses: []string{"chat", "role=log", "up-poll"}, Demo: demoAgentThread(), Snippet: `@ui.AgentThread(props) { @ui.AgentMessage(ui.AgentAssistant, "Assistant", "10:12") { ... } }`, Note: "Le thread peut être pollé pendant une réponse, mais le texte conversationnel ne devient jamais une autorisation implicite."},
		{Group: "Tier 4 — Agentic UI", ID: "agent-composer", Title: "AgentComposer", Purpose: "Demande libre avec contexte joint visible et retirable.", Uses: []string{"textarea", "context chips", "form"}, Demo: demoAgentComposer(), Snippet: `@ui.AgentComposer(ui.AgentComposerProps{Action: "/agent/messages", Contexts: context})`, Note: "Le contexte envoyé est visible avant soumission. Un identifiant caché ou un état client opaque ne doit pas élargir silencieusement la portée."},
		{Group: "Tier 4 — Agentic UI", ID: "agent-suggestions", Title: "AgentSuggestionSet", Purpose: "Proposer un ensemble borné de prochaines actions avec leur risque.", Uses: []string{"cards", "risk badges"}, Demo: demoAgentSuggestions(), Snippet: `@ui.AgentSuggestionSet("Que souhaitez-vous faire ?", suggestions, "#agent")`, Note: "Les suggestions formulent une intention ; le backend produit ensuite le plan exact et reste libre de relever son niveau de risque."},
		{Group: "Tier 4 — Agentic UI", ID: "action-plan", Title: "ActionPlan", Purpose: "Prévisualiser étapes, outils, cibles et risques avant approbation.", Uses: []string{"ordered plan", "data-plan-version"}, Demo: demoActionPlan(), Snippet: `@ui.ActionPlan(ui.ActionPlanProps{PlanID: plan.ID, Version: plan.Version, Steps: steps})`, Note: "Un plan est informatif. Seul ApprovalCard envoie une autorisation au backend."},
		{Group: "Tier 4 — Agentic UI", ID: "agent-impact", Title: "AgentImpactSummary", Purpose: "Quantifier lectures, modifications, effets externes et coût.", Uses: []string{"scope", "side effects"}, Demo: demoAgentImpact(), Snippet: `@ui.AgentImpactSummary(impactItems)`, Note: "Les nombres concrets sont plus utiles que des qualificatifs vagues comme « quelques comptes » ou « faible coût »."},
		{Group: "Tier 4 — Agentic UI", ID: "approval-card", Title: "ApprovalCard", Purpose: "Approbation globale liée à une version immuable du plan.", Uses: []string{"plan digest", "confirmation", "risk gate"}, Demo: demoApprovalCard(), Snippet: `@ui.ApprovalCard(approval, impact)`, Note: "plan_id, plan_version et plan_digest sont soumis ensemble. Le backend doit refuser toute approbation expirée ou ne correspondant plus au plan courant."},
		{Group: "Tier 4 — Agentic UI", ID: "tool-runs", Title: "AgentToolRunTimeline", Purpose: "Montrer l'exécution observable des outils sans exposer le raisonnement privé du modèle.", Uses: []string{"timeline", "up-poll", "tool status"}, Demo: demoToolTimeline(), Snippet: `@ui.AgentToolRunTimeline("runs", runs, true, "/agent/runs")`, Note: "La timeline expose outils, cibles, résultats et erreurs utiles à l'opérateur — pas une chaîne de pensée."},
		{Group: "Tier 4 — Agentic UI", ID: "result-receipt", Title: "AgentResultReceipt", Purpose: "Conserver la preuve durable de ce qui a réellement été exécuté.", Uses: []string{"receipt", "diff", "sources"}, Demo: demoResultReceipt(), Snippet: `@ui.AgentResultReceipt(receipt)`, Note: "Le reçu référence la version du plan, résume les changements avant/après et relie sources et livrables."},
		{Group: "Tier 4 — Agentic UI", ID: "undo-action", Title: "AgentUndoAction", Purpose: "Annuler une exécution réversible à partir de son reçu.", Uses: []string{"receipt_id", "expiry"}, Demo: demoUndoAction(), Snippet: `@ui.AgentUndoAction(ui.AgentUndoProps{ReceiptID: receipt.ID, Available: true, ...})`, Note: "L'annulation vise le reçu d'exécution, jamais le dernier message de la conversation."},
		{Group: "Tier 4 — Agentic UI", ID: "human-handoff", Title: "AgentHumanHandoff", Purpose: "Arrêter l'automatisation et transmettre le contexte auditable à une équipe humaine.", Uses: []string{"handoff", "conversation_id"}, Demo: demoHumanHandoff(), Snippet: `@ui.AgentHumanHandoff(ui.AgentHandoffProps{ConversationID: conversation.ID, Queue: "Sécurité"})`, Note: "Le transfert est une issue normale, pas une erreur : l'agent doit savoir s'arrêter lorsque la politique backend le demande."},
	}
}
