package main

import (
	"net/http"
	"net/url"

	"github.com/esrid/back-office-kit/ui"
)

type PolicyView struct {
	Operator string
	Props    ui.PolicySimulatorProps
}

var policyActors = []ui.Option{
	{Value: "amelie", Label: "Amélie Rousseau (Administrateur)"},
	{Value: "bruno", Label: "Bruno Keller (Membre)"},
	{Value: "dario", Label: "Dario Esposito (Lecture seule)"},
}

var policyResources = []ui.Option{
	{Value: "invoices", Label: "Factures"},
	{Value: "accounts", Label: "Comptes"},
	{Value: "api_keys", Label: "Clés d'API"},
}

var policyActions = []ui.Option{
	{Value: "read", Label: "Lire"},
	{Value: "write", Label: "Modifier"},
	{Value: "delete", Label: "Supprimer"},
}

// evaluatePolicy is a deliberately small rule engine: enough to make the
// simulator answer differently for different inputs, which is what the screen
// exists to demonstrate.
func evaluatePolicy(actor, resource, action string) ui.PolicySimulation {
	role := map[string]string{"amelie": "admin", "bruno": "member", "dario": "viewer"}[actor]

	allowed := false
	reason := ""
	switch {
	case role == "admin":
		allowed = action != "delete" || resource != "api_keys"
		if !allowed {
			reason = "La suppression d'une clé d'API demande une élévation de session."
		} else {
			reason = "Le rôle Administrateur couvre cette action."
		}
	case role == "member":
		allowed = action == "read" || (action == "write" && resource == "invoices")
		if !allowed {
			reason = "Le rôle Membre ne peut modifier que les factures."
		} else {
			reason = "Le rôle Membre couvre cette action."
		}
	default:
		allowed = action == "read"
		if !allowed {
			reason = "Le rôle Lecture seule n'autorise aucune écriture."
		} else {
			reason = "La lecture est ouverte à tous les rôles."
		}
	}

	decision := ui.AccessDenied
	if allowed {
		decision = ui.AccessAllowed
	}
	return ui.PolicySimulation{
		DecisionID: "dec_" + actor + "_" + resource + "_" + action,
		Decision:   decision,
		Reason:     reason,
		Policy:     "billing.rbac.v3",
		EvaluatedAt: "à l'instant",
		Facts: []ui.PolicyFact{
			{Label: "Rôle effectif", Value: role},
			{Label: "Ressource", Value: resource},
			{Label: "Action", Value: action},
		},
	}
}

func policyPropsFor(q url.Values) ui.PolicySimulatorProps {
	p := ui.PolicySimulatorProps{
		ID: "policy", Title: "Simuler une décision",
		Description: "Choisissez qui fait quoi sur quoi. Rien n'est modifié.",
		Action:      "/policy", Target: "#policy",
		Actors: policyActors, Resources: policyResources, Actions: policyActions,
		Actor: q.Get("actor"), Resource: q.Get("resource"), ActionValue: q.Get("action"),
		Context: q.Get("context"),
	}
	if p.Actor != "" && p.Resource != "" && p.ActionValue != "" {
		sim := evaluatePolicy(p.Actor, p.Resource, p.ActionValue)
		p.Simulation = &sim
	}
	return p
}

func registerPolicy(mux *http.ServeMux, operator string) {
	mux.HandleFunc("GET /policy", func(w http.ResponseWriter, r *http.Request) {
		render(r.Context(), w, PolicyPage(PolicyView{
			Operator: operator,
			Props:    policyPropsFor(r.URL.Query()),
		}))
	})
}
