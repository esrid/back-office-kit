package main

func tier2Sections() []Section {
	return []Section{
		{
			Group: "Tier 2 — Après", ID: "statrow", Title: "StatRow",
			Purpose: "Une ligne de KPI lisible qui devient verticale sur petit écran.", Uses: []string{"stats", "responsive"},
			Demo: demoStatRow(), Snippet: `@ui.StatRow([]ui.Stat{{Label: "MRR", Value: "48 920 €", Change: "+8,4 %", Tone: "success"}})`,
			Note: "Les couleurs de tendance restent sémantiques et les valeurs utilisent une taille mesurée pour ne pas casser les tuiles avec de vrais nombres.",
		},
		{
			Group: "Tier 2 — Après", ID: "slideover", Title: "SlideOver",
			Purpose: "Éditer dans un drawer droit puis rafraîchir la liste parente à l'acceptation.", Uses: []string{"up-layer", "drawer", "up-accept-location"},
			Demo: demoSlideOver(), Snippet: `@ui.SlideOverTrigger(ui.SlideOverProps{Href: "/users/42/edit", Label: "Modifier", AcceptLocation: "/users/*", ReloadTarget: "#users"})

@ui.SlideOver("Modifier l'utilisateur") { ... }`,
			Note: "Le trigger utilise la syntaxe Unpoly actuelle new drawer et up-position=right ; la page du drawer reste une réponse serveur ordinaire.",
		},
		{
			Group: "Tier 2 — Après", ID: "detailtabs", Title: "DetailTabs",
			Purpose: "Navigation entre sections d'une fiche, avec état actif porté par l'URL.", Uses: []string{"tabs", "up-follow"},
			Demo: demoDetailTabs(), Snippet: `@ui.DetailTabs(tabs, r.URL.Path, "#record-detail")`,
			Note: "Ce sont de vrais liens : historique navigateur, ouverture dans un nouvel onglet et fonctionnement sans JavaScript restent disponibles.",
		},
		{
			Group: "Tier 2 — Après", ID: "bulkactionbar", Title: "BulkActionBar",
			Purpose: "Actions groupées visibles uniquement lorsqu'au moins une ligne est sélectionnée.", Uses: []string{"sticky", "aria-live"},
			Demo: demoBulkActionBar(), Snippet: `@ui.BulkActionBar(len(selected), "utilisateur sélectionné", "utilisateurs sélectionnés") { ... }`,
			Note: "La barre est annoncée comme statut dynamique et reste collée au bas de la zone pendant le défilement.",
		},
		{
			Group: "Tier 2 — Après", ID: "authcard", Title: "AuthCard",
			Purpose: "Login concentré dans une card, soutenu par un hero produit sur grand écran.", Uses: []string{"card", "hero"}, StageClass: "gal-stage--app",
			Demo: demoTier2AuthCard(), Snippet: `@ui.AuthCard(ui.AuthCardProps{Product: "Acme Admin", Title: "Connexion", HeroTitle: "Pilotez vos opérations."}) { ... }`,
			Note: "Le hero disparaît sur mobile pour réserver l'espace à la tâche ; le nom du produit reste visible dans la card.",
		},
		{
			Group: "Tier 2 — Après", ID: "filefield", Title: "FileField",
			Purpose: "Sélection de fichier, aide, erreur et progression d'envoi dans un seul fieldset.", Uses: []string{"file-input", "progress"},
			Demo: demoFileField(), Snippet: `@ui.FileField(ui.FileFieldProps{Name: "import", Label: "Fichier CSV", Accept: ".csv", Uploading: true, Progress: 64})`,
			Note: "Progress est borné entre 0 et 100. Le transport peut être natif, Unpoly ou direct-to-storage ; le composant ne couple pas l'UI à ce choix.",
		},
		{
			Group: "Tier 2 — Après", ID: "audittimeline", Title: "AuditTimeline",
			Purpose: "Historique chronologique avec auteur, horodatage et nature de l'événement.", Uses: []string{"timeline", "time"},
			Demo: demoAuditTimeline(), Snippet: `@ui.AuditTimeline(events)`,
			Note: "La timeline reste compacte et sur un seul axe : un journal d'audit se parcourt mieux qu'une alternance gauche/droite.",
		},
		{
			Group: "Tier 2 — Après", ID: "themetoggle", Title: "ThemeToggle",
			Purpose: "Basculer clair/sombre avec le theme-controller CSS de daisyUI.", Uses: []string{"theme-controller", "toggle"},
			Demo: demoThemeToggle(), Snippet: `@ui.ThemeToggle(isDark)`,
			Note: "Le contrôleur change le thème sans JavaScript. La persistance éventuelle reste une responsabilité de l'application.",
		},
	}
}
