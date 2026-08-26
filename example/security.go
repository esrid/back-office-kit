package main

import (
	"net/http"
	"net/url"

	"github.com/esrid/back-office-kit/ui"
)

type SecurityView struct {
	Operator string
	Severity string
	Events   []ui.SecurityEvent
	Sessions []ui.SecuritySession
}

func severityOptions() []ui.Option {
	return []ui.Option{
		{Value: "", Label: "Toutes gravités"},
		{Value: "critical", Label: "Critique"},
		{Value: "warning", Label: "Avertissement"},
		{Value: "info", Label: "Information"},
	}
}

func allSecurityEvents() []ui.SecurityEvent {
	return []ui.SecurityEvent{
		{ID: "e1", Title: "Connexion depuis un appareil inconnu", Severity: ui.SecurityEventWarning,
			Summary: "Première connexion depuis ce navigateur.", Actor: "amelie@acme.co",
			IP: "88.120.4.17", OccurredAt: "il y a 12 minutes", RequestID: "req_a91f", Href: "/security?event=a91f"},
		{ID: "e2", Title: "Clé d'API révoquée", Severity: ui.SecurityEventInfo,
			Summary: "La clé « intégration compta » a été supprimée.", Actor: "nora@acme.co",
			OccurredAt: "il y a 2 heures", RequestID: "req_b30c", Href: "/security?event=b30c"},
		{ID: "e3", Title: "Cinq échecs d'authentification consécutifs", Severity: ui.SecurityEventCritical,
			Summary: "Le compte a été temporairement verrouillé.", Actor: "dario@partner.io",
			IP: "203.0.113.44", OccurredAt: "hier", RequestID: "req_c72d", Href: "/security?event=c72d"},
		{ID: "e4", Title: "Rôle élevé en Administrateur", Severity: ui.SecurityEventWarning,
			Summary: "Élévation approuvée par nora@acme.co.", Actor: "amelie@acme.co",
			OccurredAt: "24 août", RequestID: "req_d18e", Href: "/security?event=d18e"},
	}
}

// filterEvents keeps the severity filter in the query string. Pure.
func filterEvents(all []ui.SecurityEvent, severity string) []ui.SecurityEvent {
	if severity == "" {
		return all
	}
	var out []ui.SecurityEvent
	for _, e := range all {
		if string(e.Severity) == severity {
			out = append(out, e)
		}
	}
	return out
}

func demoSessions() []ui.SecuritySession {
	return []ui.SecuritySession{
		{ID: "s1", Device: "MacBook Pro", Browser: "Chrome 141", Location: "Lyon, FR",
			IP: "88.120.4.17", LastSeen: "à l'instant", Current: true},
		{ID: "s2", Device: "iPhone 17", Browser: "Safari", Location: "Lyon, FR",
			IP: "88.120.4.17", LastSeen: "il y a 3 heures", Revocable: true},
		{ID: "s3", Device: "Windows", Browser: "Firefox 140", Location: "Francfort, DE",
			IP: "203.0.113.44", LastSeen: "il y a 6 jours", Revocable: true},
	}
}

func registerSecurity(mux *http.ServeMux, operator string) {
	events := allSecurityEvents()
	mux.HandleFunc("GET /security", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		render(r.Context(), w, SecurityPage(SecurityView{
			Operator: operator,
			Severity: q.Get("severity"),
			Events:   filterEvents(events, q.Get("severity")),
			Sessions: demoSessions(),
		}))
	})
	_ = url.Values{}
}
