package main

import (
	"net/http"

	"github.com/esrid/back-office-kit/ui"
)

type ReviewView struct {
	Operator string
	Request  ui.ApprovalRequest
	Plan     ui.ActionPlanProps
	Approval ui.ApprovalCardProps
}

// findApproval resolves a request id. Pure.
func findApproval(items []ui.ApprovalRequest, id string) (ui.ApprovalRequest, bool) {
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	return ui.ApprovalRequest{}, false
}

func reviewFor(req ui.ApprovalRequest) ReviewView {
	return ReviewView{
		Request: req,
		Plan: ui.ActionPlanProps{
			PlanID: req.PlanID, Version: req.PlanVersion,
			Title: req.Title, Summary: req.Summary, Risk: req.Risk,
			Steps: []ui.AgentPlanStep{
				{Title: "Identifier les comptes concernés", Description: "Facture échue depuis plus de 30 jours.",
					Tool: "billing.search", Target: "comptes actifs", Risk: ui.AgentRiskRead},
				{Title: "Suspendre 12 comptes", Description: "Bloquer les connexions en conservant les données.",
					Tool: "accounts.suspend", Target: "12 comptes", Risk: ui.AgentRiskWrite},
				{Title: "Informer les propriétaires", Description: "Un e-mail par compte avec le lien de régularisation.",
					Tool: "mail.send", Target: "12 destinataires", Risk: ui.AgentRiskSensitive},
			},
		},
		Approval: ui.ApprovalCardProps{
			PlanID: req.PlanID, PlanVersion: req.PlanVersion, PlanDigest: "sha256:8f31c2a4",
			Title: "Approuver le plan complet", Risk: req.Risk,
			Summary:      "Cette approbation autorise les trois étapes exactement telles qu'elles sont présentées.",
			Action:       "/approvals/" + req.ID + "/approve",
			RejectHref:   "/approvals/" + req.ID + "/reject",
			ConfirmLabel: "Approuver et exécuter",
			RequireText:  "SUSPENDRE 12 COMPTES",
			ExpiresAt:    "dans 10 minutes",
			GlobalApproval: true,
		},
	}
}

func registerReview(mux *http.ServeMux, operator string, items []ui.ApprovalRequest) {
	mux.HandleFunc("GET /approvals/{id}", func(w http.ResponseWriter, r *http.Request) {
		req, ok := findApproval(items, r.PathValue("id"))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			render(r.Context(), w, ui.ErrorPage(ui.ErrorPageProps{
				Code: "404", Title: "Demande introuvable",
				Message:  "Cette demande d'approbation n'existe pas ou a déjà été traitée.",
				HomeHref: "/approvals",
			}))
			return
		}
		v := reviewFor(req)
		v.Operator = operator
		render(r.Context(), w, ReviewPage(v))
	})
}
