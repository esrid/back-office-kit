package main

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/esrid/back-office-kit/ui"
)

type RecordsView struct {
	Operator    string
	Subtitle    string
	Items       []ui.MasterDetailItem
	Current     *ui.MasterDetailItem
	CurrentHref string
	Columns     []ui.ManagedColumn
	Density     ui.TableDensity
}

type ApprovalsView struct {
	Operator string
	Subtitle string
	Items    []ui.ApprovalRequest
}

func seedRecords() []ui.MasterDetailItem {
	rows := []struct{ id, title, desc, meta string }{
		{"DOS-4192", "Litige facturation Acme", "Contestation de la facture #4192", "En cours"},
		{"DOS-4193", "Demande de remboursement", "Commande annulée après expédition", "En attente"},
		{"DOS-4194", "Accès révoqué par erreur", "Le compte a perdu ses droits admin", "Résolu"},
		{"DOS-4195", "Double prélèvement", "Deux débits le même jour", "En cours"},
	}
	out := make([]ui.MasterDetailItem, len(rows))
	for i, r := range rows {
		out[i] = ui.MasterDetailItem{
			ID: r.id, Title: r.title, Description: r.desc, Meta: r.meta,
			Href: "/records?dossier=" + r.id,
		}
	}
	return out
}

// selectRecord resolves the query string to a record. Pure.
func selectRecord(items []ui.MasterDetailItem, q url.Values) (*ui.MasterDetailItem, string) {
	want := q.Get("dossier")
	if want == "" {
		return nil, ""
	}
	for i := range items {
		if items[i].ID == want {
			return &items[i], items[i].Href
		}
	}
	return nil, ""
}

// selectColumns applies the persisted visibility to the managed columns. Pure.
func selectColumns(q url.Values) ([]ui.ManagedColumn, ui.TableDensity) {
	defaults := []ui.ManagedColumn{
		{Key: "ref", Label: "Référence", Visible: true, Locked: true},
		{Key: "title", Label: "Objet", Visible: true},
		{Key: "owner", Label: "Responsable", Visible: true},
		{Key: "updated", Label: "Mise à jour", Visible: false},
		{Key: "amount", Label: "Montant", Visible: false},
	}
	if hidden, ok := q["hide"]; ok {
		for _, key := range hidden {
			for i := range defaults {
				if defaults[i].Key == key && !defaults[i].Locked {
					defaults[i].Visible = false
				}
			}
		}
	}
	density := ui.DensityComfortable
	if d := q.Get("density"); d != "" {
		density = ui.TableDensity(d)
	}
	return defaults, density
}

func seedApprovals() []ui.ApprovalRequest {
	return []ui.ApprovalRequest{
		{
			ID: "req_1", PlanID: "plan_42", PlanVersion: "7",
			Title:     "Suspendre 12 comptes en retard",
			Summary:   "Le backend a identifié 12 comptes selon la politique de facturation.",
			Requester: "Assistant opérations", RequestedAt: "il y a 8 minutes",
			ExpiresAt: "dans 52 minutes", Risk: ui.AgentRiskSensitive,
			ReviewHref: "/approvals/req_1",
		},
		{
			ID: "req_2", PlanID: "plan_43", PlanVersion: "2",
			Title:     "Recalculer les taxes du trimestre",
			Summary:   "Recalcul sur 1 284 factures, réversible.",
			Requester: "Nora Martin", RequestedAt: "hier",
			Risk: ui.AgentRiskWrite, ReviewHref: "/approvals/req_2",
		},
	}
}

func registerRecords(mux *http.ServeMux, operator string) {
	records := seedRecords()
	approvals := seedApprovals()

	mux.HandleFunc("GET /records", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		current, href := selectRecord(records, q)
		columns, density := selectColumns(q)
		render(r.Context(), w, RecordsPage(RecordsView{
			Operator:    operator,
			Subtitle:    strconv.Itoa(len(records)) + " dossiers ouverts",
			Items:       records,
			Current:     current,
			CurrentHref: href,
			Columns:     columns,
			Density:     density,
		}))
	})

	registerReview(mux, operator, approvals)
	registerAgent(mux, operator)
	registerPolicy(mux, operator)

	mux.HandleFunc("GET /approvals", func(w http.ResponseWriter, r *http.Request) {
		render(r.Context(), w, ApprovalsPage(ApprovalsView{
			Operator: operator,
			Subtitle: strconv.Itoa(len(approvals)) + " demandes en attente",
			Items:    approvals,
		}))
	})
}
