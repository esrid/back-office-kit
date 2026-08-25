// Command example is a runnable page built with the kit: filtering, sorting and
// pagination all driven by the query string, with no client state.
//
//	go run ./example    →  http://localhost:8080
package main

import (
	"cmp"
	"context"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/esrid/back-office-kit/ui"
)

const perPage = 8

type User struct {
	Name   string
	Email  string
	Role   string
	Status string
	Seats  int
	Joined time.Time
}

// UsersView is everything the page needs. Handlers build it; templates only read.
type UsersView struct {
	Operator string
	Subtitle string
	Rows     []User
	Page     int
	Total    int
	Query    url.Values
	Flashes  []ui.Flash
}

func nav() []ui.NavItem {
	return []ui.NavItem{
		{Label: "Utilisateurs", Href: "/"},
		{Label: "Facturation", Href: "/billing"},
		{Label: "Réglages", Href: "/settings"},
	}
}

func filters() []ui.Filter {
	return []ui.Filter{
		{Name: "q", Label: "Recherche", Placeholder: "Nom ou e-mail"},
		{Name: "status", Label: "Statut", Options: []ui.Option{
			{Value: "active", Label: "Actif"},
			{Value: "pending", Label: "En attente"},
			{Value: "failed", Label: "Échec"},
		}},
	}
}

func columns(v UsersView) []ui.Column[User] {
	seats := ui.Text("Sièges", "seats", func(u User) int { return u.Seats })
	seats.Class = "text-right tabular-nums"

	return []ui.Column[User]{
		ui.Text("Nom", "name", func(u User) string { return u.Name }),
		ui.Text("E-mail", "email", func(u User) string { return u.Email }),
		{Header: "Rôle", Cell: func(u User) templ.Component { return roleCell(u.Role) }},
		{Header: "Statut", Sort: "status", Cell: func(u User) templ.Component {
			return ui.StatusBadge(u.Status)
		}},
		{Header: "Inscrit", Sort: "joined", Cell: func(u User) templ.Component {
			return ui.Timestamp(ui.TimestampProps{At: u.Joined, Now: time.Now(), Location: time.UTC})
		}},
		seats,
	}
}

// selectUsers applies the query string to the dataset. Pure: same query, same
// page, every time.
func selectUsers(all []User, q url.Values) (rows []User, total int) {
	needle := strings.ToLower(strings.TrimSpace(q.Get("q")))
	status := q.Get("status")

	for _, u := range all {
		if needle != "" &&
			!strings.Contains(strings.ToLower(u.Name), needle) &&
			!strings.Contains(strings.ToLower(u.Email), needle) {
			continue
		}
		if status != "" && !strings.EqualFold(u.Status, status) {
			continue
		}
		rows = append(rows, u)
	}

	dir := 1
	if q.Get("dir") == "desc" {
		dir = -1
	}
	slices.SortStableFunc(rows, func(a, b User) int {
		switch q.Get("sort") {
		case "email":
			return dir * cmp.Compare(a.Email, b.Email)
		case "status":
			return dir * cmp.Compare(a.Status, b.Status)
		case "seats":
			return dir * cmp.Compare(a.Seats, b.Seats)
		case "joined":
			return dir * a.Joined.Compare(b.Joined)
		default:
			return dir * cmp.Compare(a.Name, b.Name)
		}
	})

	total = len(rows)
	page := max(atoiOr(q.Get("page"), 1), 1)
	start := min((page-1)*perPage, total)
	return rows[start:min(start+perPage, total)], total
}

func atoiOr(s string, fallback int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return fallback
}

func main() {
	all := seed()

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("dist"))))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		rows, total := selectUsers(all, q)

		view := UsersView{
			Operator: "amelie@acme.co",
			Subtitle: strconv.Itoa(total) + " comptes",
			Rows:     rows,
			Page:     max(atoiOr(q.Get("page"), 1), 1),
			Total:    total,
			Query:    q,
		}
		render(r.Context(), w, UsersPage(view))
	})

	// Cross-origin POSTs are rejected without a token or a cookie. Go 1.25+.
	protected := http.NewCrossOriginProtection().Handler(mux)

	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", protected))
}

func render(ctx context.Context, w http.ResponseWriter, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := c.Render(ctx, w); err != nil {
		log.Println("render:", err)
	}
}

func seed() []User {
	base := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	names := []struct {
		Name, Email, Role, Status string
		Seats, DaysAgo            int
	}{
		{"Amélie Rousseau", "amelie@acme.co", "Administrateur", "Active", 12, 530},
		{"Bruno Keller", "bruno@acme.co", "Membre", "Pending", 3, 12},
		{"Chloé Nguyen", "chloe@acme.co", "Membre", "Active", 5, 210},
		{"Dario Esposito", "dario@partner.io", "Lecture seule", "Failed", 1, 3},
		{"Elin Sørensen", "elin@acme.co", "Administrateur", "Draft", 8, 90},
		{"Farid Benali", "farid@acme.co", "Membre", "Active", 4, 45},
		{"Greta Lindqvist", "greta@acme.co", "Membre", "Pending", 2, 1},
		{"Hugo Marchand", "hugo@acme.co", "Lecture seule", "Active", 6, 400},
		{"Ines Ferreira", "ines@acme.co", "Membre", "Active", 3, 150},
		{"Jonas Weber", "jonas@partner.io", "Membre", "Failed", 1, 7},
		{"Klara Novak", "klara@acme.co", "Administrateur", "Active", 15, 620},
		{"Léo Fontaine", "leo@acme.co", "Membre", "Draft", 2, 30},
	}
	users := make([]User, len(names))
	for i, n := range names {
		users[i] = User{
			Name: n.Name, Email: n.Email, Role: n.Role, Status: n.Status, Seats: n.Seats,
			Joined: base.AddDate(0, 0, -n.DaysAgo),
		}
	}
	return users
}
