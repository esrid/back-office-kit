// Command gallery renders a static, self-contained preview of the complete
// back-office component kit. The stylesheet is embedded so the page needs no
// server and no CDN.
package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/esrid/back-office-kit/ui"
)

// Section is one entry of the gallery.
type Section struct {
	N          int
	Signature  string
	Search     string
	Group      string
	ID         string
	Title      string
	Purpose    string
	Uses       []string
	StageClass string
	Snippet    string
	Note       string
	Demo       templ.Component
}

type User struct {
	Name   string
	Email  string
	Role   string
	Status string
	Seats  int
}

func demoNav() []ui.NavItem {
	return []ui.NavItem{
		{Label: "Tableau de bord", Href: "/"},
		{Label: "Utilisateurs", Href: "/users"},
		{Label: "Facturation", Href: "/billing"},
		{Label: "Journal d'audit", Href: "/audit"},
		{Label: "Réglages", Href: "/settings"},
	}
}

func demoUsers() []User {
	return []User{
		{"Amélie Rousseau", "amelie@acme.co", "Administrateur", "Active", 12},
		{"Bruno Keller", "bruno@acme.co", "Membre", "Pending", 3},
		{"Chloé Nguyen", "chloe@acme.co", "Membre", "Active", 5},
		{"Dario Esposito", "dario@partner.io", "Lecture seule", "Failed", 1},
		{"Elin Sørensen", "elin@acme.co", "Administrateur", "Draft", 8},
	}
}

func demoFilters() []ui.Filter {
	return []ui.Filter{
		{Name: "q", Label: "Recherche", Placeholder: "Nom ou e-mail"},
		{Name: "status", Label: "Statut", Options: []ui.Option{
			{Value: "active", Label: "Actif"},
			{Value: "pending", Label: "En attente"},
			{Value: "failed", Label: "Échec"},
		}},
		{Name: "role", Label: "Rôle", Options: []ui.Option{
			{Value: "admin", Label: "Administrateur"},
			{Value: "member", Label: "Membre"},
		}},
	}
}

func settingsNav() []ui.NavItem {
	return []ui.NavItem{
		{Label: "Général", Href: "/settings"},
		{Label: "Notifications", Href: "/settings/notifications"},
		{Label: "Sécurité", Href: "/settings/security"},
		{Label: "Facturation", Href: "/settings/billing"},
		{Label: "Clés d'API", Href: "/settings/keys"},
	}
}

func demoQuery() url.Values {
	return url.Values{"sort": {"email"}, "dir": {"asc"}, "page": {"3"}, "q": {"acme"}}
}

func demoColumns() []ui.Column[User] {
	return []ui.Column[User]{
		ui.Text("Nom", "name", func(u User) string { return u.Name }),
		ui.Text("E-mail", "email", func(u User) string { return u.Email }),
		ui.Text("Rôle", "", func(u User) string { return u.Role }),
		{
			Header: "Statut",
			Sort:   "status",
			Cell:   func(u User) templ.Component { return ui.StatusBadge(u.Status) },
		},
		seats(),
	}
}

func seats() ui.Column[User] {
	c := ui.Text("Sièges", "seats", func(u User) int { return u.Seats })
	c.Class = "text-right tabular-nums"
	return c
}

func sections() []Section {
	q := demoQuery()

	secs := []Section{
		{
			Group: "Tier 1 — Composants", ID: "shell", Title: "Shell",
			Purpose:    "Le cadre : sidebar rétractable, barre haute, zone de flashes, région <main>.",
			Uses:       []string{"drawer", "navbar", "menu"},
			StageClass: "gal-stage--app",
			Demo:       demoShell(),
			Snippet: `@ui.Shell("Utilisateurs", nav, "/users", user, flashes) {
    <h1 class="text-lg font-semibold">Utilisateurs</h1>
    ...
}`,
			Note: "<main> est une cible principale Unpoly par défaut : un lien [up-follow] ne remplace que le contenu, sans [up-target]. L'état actif du menu est calculé côté serveur — la page reste juste sans JS.",
		},
		{
			Group: "Tier 1 — Composants", ID: "flashes", Title: "Flashes",
			Purpose: "Messages de confirmation après action, mis à jour à chaque rendu.",
			Uses:    []string{"alert", "up-flashes"},
			Demo:    demoFlashes(),
			Snippet: `@ui.Flashes([]ui.Flash{
    {Tone: "success", Text: "Utilisateur mis à jour."},
})`,
			Note: "Unpoly met à jour [up-flashes] même quand l'élément n'est pas ciblé, et un conteneur vide n'efface pas les messages existants — donc on l'émet toujours. À la fermeture d'un overlay, les flashes remontent au layer parent.",
		},
		{
			Group: "Tier 1 — Composants", ID: "table", Title: "DataTable",
			Purpose: "Tableau typé et triable ; les en-têtes sont des liens, l'état de tri vit dans l'URL.",
			Uses:    []string{"table", "up-target"},
			Demo:    demoTable(demoColumns(), demoUsers(), q),
			Snippet: `cols := []ui.Column[User]{
    ui.Text("Nom", "name", func(u User) string { return u.Name }),
    {Header: "Statut", Sort: "status",
     Cell: func(u User) templ.Component { return ui.StatusBadge(u.Status) }},
}

@ui.DataTable("users", cols, users, r.URL.Query(), nil) {
    @ui.Pagination(page, total, 20, r.URL.Query(), "#users")
}`,
			Note: "ui.Text infère les deux paramètres de type depuis l'accesseur, donc pas de func(T) templ.Component à écrire pour une cellule texte. Les liens de tri ciblent #users, qui enveloppe le tableau ET la pagination : les deux sont échangés ensemble.",
		},
		{
			Group: "Tier 1 — Composants", ID: "pagination", Title: "Pagination",
			Purpose: "Sélecteur de page avec ellipses, calculé par une fonction pure et testée.",
			Uses:    []string{"join", "btn"},
			Demo:    ui.Pagination(3, 142, 20, q, "#users"),
			Snippet: `@ui.Pagination(current, totalItems, 20, r.URL.Query(), "#users")

// ui.Pages(5, 10) -> [1 0 4 5 6 0 10]   (0 = ellipse)`,
			Note: "PageHref et SortHref ne mutent jamais l'url.Values reçue et préservent les filtres en place ; page=1 est retiré de l'URL. Un test parcourt toutes les combinaisons jusqu'à 40 pages pour vérifier qu'aucune page n'est dupliquée ni sautée sans ellipse.",
		},
		{
			Group: "Tier 1 — Composants", ID: "empty", Title: "EmptyState",
			Purpose: "Ce qui remplit la place des lignes quand il n'y en a pas — avec une sortie.",
			Uses:    []string{"utilities"},
			Demo:    demoEmpty(),
			Snippet: `@ui.EmptyState("Aucun utilisateur",
    "Personne ne correspond à ce filtre.", inviteButton())`,
			Note: "Distinguez « rien ne correspond au filtre » de « rien n'existe encore » : ce ne sont pas les mêmes mots ni la même action.",
		},
		{
			Group: "Tier 1 — Composants", ID: "field", Title: "Field · SelectField",
			Purpose: "Champ étiqueté avec indice et erreur serveur, prêt pour la validation Unpoly.",
			Uses:    []string{"fieldset", "input", "select"},
			Demo:    demoFields(),
			Snippet: `@ui.Field(ui.FieldProps{
    Name: "email", Label: "E-mail", Type: "email",
    Value: form.Email, Error: form.Errors["email"],
    Required: true, Validate: true,
})`,
			Note: "Chaque champ est son propre <fieldset> : c'est ce que [up-validate] cible par défaut (X-Up-Target: fieldset:has(#f-email)). Sans ce conteneur, le serveur n'a aucun fragment stable à re-rendre. Le serveur valide sans committer et renvoie le formulaire ; seul le champ modifié bouge.",
		},
		{
			Group: "Tier 1 — Composants", ID: "status", Title: "StatusBadge",
			Purpose: "Un état a la même couleur sur tous les écrans, décidé à un seul endroit.",
			Uses:    []string{"badge", "status"},
			Demo:    demoStatus(),
			Snippet: `@ui.StatusBadge(u.Status)   // ToneFor() est pure, insensible à la casse`,
			Note:    "Les classes sont écrites en toutes lettres dans une table Go, jamais concaténées : \"badge-\" + tone compile mais le scanner Tailwind ne le voit pas, et le CSS manque à l'exécution. Un statut inconnu retombe sur neutral au lieu de disparaître.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "pageheader", Title: "PageHeader",
			Purpose: "Le haut de chaque écran : fil d'Ariane, titre, actions secondaires, action primaire.",
			Uses:    []string{"breadcrumbs", "btn-primary"},
			Demo:    demoPageHeaderComplete(),
			Snippet: `@ui.PageHeader(ui.PageProps{
    Crumbs:  []ui.Crumb{{Label: "Utilisateurs", Href: "/users"}, {Label: u.Name}},
    Title:   u.Name,
    Actions: secondaryActions(), PrimaryAction: editAction(),
})`,
			Note: "Un Crumb sans Href rend du texte marqué aria-current=\"page\" : le dernier maillon n'est pas un lien vers la page où l'on est déjà. PrimaryAction est rendue après les actions secondaires pour garder la décision principale au bord droit.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "filterbar", Title: "FilterBar",
			Purpose: "Recherche, listes et période, toutes dans une query string partageable.",
			Uses:    []string{"input", "select", "date", "up-autosubmit"},
			Demo:    demoFilterBarComplete(),
			Snippet: `@ui.FilterBar([]ui.Filter{
    {Name: "q", Label: "Recherche", Placeholder: "Nom ou e-mail"},
    {Name: "status", Label: "Statut", Options: statusOptions},
    {Kind: ui.FilterDateRange, Name: "from", EndName: "to", Label: "Du", EndLabel: "au"},
}, r.URL.Query(), "#results")`,
			Note: "Un <form> portant [up-target] n'a pas besoin de [up-submit] : Unpoly le prend en charge. [up-autosubmit] soumet au changement, [up-watch-delay=\"300\"] sur la recherche évite une requête par frappe. Sans JS il reste un formulaire GET, d'où le bouton dans <noscript>. Le bouton Réinitialiser retire les deux bornes de date et la pagination mais garde le tri : effacer une recherche ne doit pas jeter la colonne triée.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "indexpage", Title: "IndexPage",
			Purpose: "L'écran dont un back-office a vingt exemplaires : en-tête, filtres, tableau.",
			Uses:    []string{"PageHeader", "FilterBar", "DataTable"},
			Demo:    demoIndexPage(demoColumns(), demoUsers(), demoQuery()),
			Snippet: `@ui.IndexPage(ui.PageProps{
    Title: "Utilisateurs", Actions: inviteButton(),
}, ui.FilterBar(filters, q, "#users")) {
    @ui.DataTable("users", cols, users, q, nil) {
        @ui.Pagination(page, total, 20, q, "#users")
    }
}`,
			Note: "Filtres, tri et pagination ciblent tous #users, qui enveloppe le tableau et sa pagination : une seule zone échangée, l'en-tête et les filtres ne clignotent pas.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "detailpage", Title: "DetailPage · Panel · DetailList",
			Purpose: "Une fiche : colonne principale et panneau de métadonnées qui passe dessous en étroit.",
			Uses:    []string{"card", "list"},
			Demo:    demoDetailPage(),
			Snippet: `@ui.DetailPage(props, recordAside()) {
    @ui.Panel("Activité récente", nil) { ... }
}

@ui.DetailList([]ui.Option{
    {Label: "Identifiant", Value: u.ID},
    {Label: "Créé le", Value: u.CreatedAt.Format("2 January 2006")},
})`,
			Note: "L'aside est un paramètre templ.Component, pas un second bloc children : un composant templ n'a qu'une seule fente children. Passer nil enlève la colonne et le contenu reprend toute la largeur. DetailList sort en dl/dt/dd, ce qui conserve la relation clé-valeur pour les technologies d'assistance.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "formpage", Title: "FormPage",
			Purpose: "Création et édition : champs groupés en panneaux, actions collées au bas de l'écran.",
			Uses:    []string{"card", "fieldset", "up-submit"},
			Demo:    demoFormPage(),
			Snippet: `@ui.FormPage(props, "/users/42", saveButtons()) {
    @ui.Panel("Identité", nil) { @ui.Field(...) }
    @ui.Panel("Accès", nil)    { @ui.SelectField(...) }
}`,
			Note: "La barre d'actions est sticky bottom-0 : sur un formulaire long, Enregistrer reste atteignable sans faire défiler jusqu'en bas.",
		},
		{
			Group: "Tier 1.5 — Écrans", ID: "settingspage", Title: "SettingsPage · SettingsRow",
			Purpose: "Réglages : sous-navigation à gauche, explication à gauche de chaque contrôle.",
			Uses:    []string{"menu", "card", "toggle"},
			Demo:    demoSettingsPage(),
			Snippet: `@ui.SettingsPage(props, settingsNav(), "/settings/notifications") {
    @ui.Panel("Notifications", nil) {
        @ui.SettingsRow("Résumé quotidien", "Un e-mail chaque matin...") {
            <input type="checkbox" class="toggle toggle-sm" checked/>
        }
    }
}`,
			Note: "La description à côté du contrôle, pas dans une infobulle : lire la ligne doit suffire à savoir ce que l'interrupteur déclenche.",
		},
	}
	secs = append(secs, tier2Sections()...)
	secs = append(secs, tier3Sections()...)
	secs = append(secs, tier4Sections()...)
	secs = append(secs, tier5Sections()...)
	secs = append(secs, tier6Sections()...)
	secs = append(secs, tier7Sections()...)
	secs = append(secs, tier8Sections()...)
	// Commerce et Marketing ne sont pas des tiers : ce sont d'autres domaines,
	// donc ils ferment la liste au lieu de s'intercaler dans la progression.
	secs = append(secs, commerceSections()...)
	secs = append(secs, marketingSections()...)
	// Numbering is assigned here, not written into each literal: the sections
	// are a build order, so adding or removing one must not leave a gap or a
	// duplicate numeral behind.
	sigs := signatures("ui", "ui/marketing")
	for i := range secs {
		secs[i].N = i + 1
		secs[i].Signature = signatureFor(sigs, secs[i].Title)
		secs[i].Search = searchText(secs[i])
	}
	return secs
}

// searchText is the haystack the client filter matches against. Built here so
// the browser never has to scrape the DOM to know what a section is about.
func searchText(s Section) string {
	parts := append([]string{s.Title, s.Purpose, s.Group, s.Signature}, s.Uses...)
	return strings.ToLower(strings.Join(parts, " "))
}

func main() {
	if err := os.MkdirAll("dist", 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("dist/gallery.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := writePage(f, "dist/app.css", renderGallery); err != nil {
		log.Fatal(err)
	}
	log.Println("wrote dist/gallery.html")
}
