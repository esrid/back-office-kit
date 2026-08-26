package main

import (
	"net/http"
	"net/url"

	"github.com/esrid/back-office-kit/ui"
)

type AgentMessageView struct {
	Role   ui.AgentRole
	Author string
	Time   string
	Text   string
}

type AgentView struct {
	Operator    string
	Busy        bool
	Messages    []AgentMessageView
	Suggestions []ui.AgentSuggestion
	Contexts    []ui.AgentContext
	Runs        []ui.AgentToolRun
}

// agentSuggestions varies with the chosen angle, so a click has a visible
// effect and the URL can be checked against it.
func agentSuggestions(angle string) []ui.AgentSuggestion {
	base := []ui.AgentSuggestion{
		{Label: "Voir les comptes concernés", Description: "Ouvrir la liste filtrée sur les factures échues.",
			Href: "/agent?angle=comptes", Risk: ui.AgentRiskRead},
		{Label: "Simuler la suspension", Description: "Estimer l'impact sans rien modifier.",
			Href: "/agent?angle=simulation", Risk: ui.AgentRiskRead},
		{Label: "Préparer un plan d'action", Description: "Rédiger les étapes à faire approuver.",
			Href: "/agent?angle=plan", Risk: ui.AgentRiskWrite},
	}
	if angle == "" {
		return base
	}
	return append([]ui.AgentSuggestion{
		{Label: "Revenir à la vue d'ensemble", Description: "Abandonner l'angle « " + angle + " ».",
			Href: "/agent", Risk: ui.AgentRiskRead},
	}, base...)
}

func agentMessages(angle string) []AgentMessageView {
	msgs := []AgentMessageView{
		{Role: ui.AgentUser, Author: "Vous", Time: "10:12",
			Text: "Identifie les comptes en retard, suspends-les et préviens leurs propriétaires."},
		{Role: ui.AgentAssistant, Author: "Assistant opérations", Time: "10:13",
			Text: "J'ai trouvé 12 comptes correspondant aux règles du backend. Le plan inclut une suspension réversible et l'envoi de 12 e-mails. Vérifiez-le avant de l'approuver."},
	}
	if angle != "" {
		msgs = append(msgs, AgentMessageView{
			Role: ui.AgentSystem, Author: "Système", Time: "10:14",
			Text: "Angle d'analyse : " + angle + ".",
		})
	}
	return msgs
}

func agentRuns(busy bool) []ui.AgentToolRun {
	runs := []ui.AgentToolRun{
		{ID: "r1", Tool: "billing.search", Summary: "12 comptes trouvés", Status: ui.AgentToolSuccess,
			Started: "10:12", Duration: "1,2 s"},
		{ID: "r2", Tool: "accounts.preview", Summary: "Impact estimé", Status: ui.AgentToolSuccess,
			Started: "10:13", Duration: "0,4 s"},
	}
	if busy {
		return append(runs, ui.AgentToolRun{
			ID: "r3", Tool: "mail.render", Summary: "Génération des messages", Status: ui.AgentToolRunning, Started: "10:14",
		})
	}
	return append(runs, ui.AgentToolRun{
		ID: "r3", Tool: "accounts.suspend", Summary: "En attente d'approbation", Status: ui.AgentToolPending,
	})
}

func agentViewFor(q url.Values) AgentView {
	angle := q.Get("angle")
	return AgentView{
		Busy:        q.Get("busy") == "1",
		Messages:    agentMessages(angle),
		Suggestions: agentSuggestions(angle),
		Contexts: []ui.AgentContext{
			{Label: "Vue", Value: "Factures en retard", RemoveHref: "/agent/context/vue/remove"},
			{Label: "Période", Value: "30 derniers jours"},
		},
		Runs: agentRuns(q.Get("busy") == "1"),
	}
}

func registerAgent(mux *http.ServeMux, operator string) {
	mux.HandleFunc("GET /agent", func(w http.ResponseWriter, r *http.Request) {
		v := agentViewFor(r.URL.Query())
		v.Operator = operator
		render(r.Context(), w, AgentPage(v))
	})
}
