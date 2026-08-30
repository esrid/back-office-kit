package main

func tier2Sections() []Section {
	const g = "Tier 2 — Navigation & contexte"
	return []Section{
		{
			Group: g, ID: "statrow", Title: "StatRow",
			Purpose: "Une ligne de KPI lisible qui devient verticale sur petit écran.", Uses: []string{"stats", "responsive"},
			Demo: demoStatRow(), Snippet: `@ui.StatRow([]ui.Stat{{Label: "MRR", Value: "48 920 €", Change: "+8,4 %", Tone: "success"}})`,
			Note: "Les couleurs de tendance restent sémantiques et les valeurs utilisent une taille mesurée pour ne pas casser les tuiles avec de vrais nombres.",
		},
		{
			Group: g, ID: "slideover", Title: "SlideOver",
			Purpose: "Éditer dans un drawer droit puis rafraîchir la liste parente à l'acceptation.", Uses: []string{"up-layer", "drawer", "up-accept-location"},
			Demo: demoSlideOver(), Snippet: `@ui.SlideOverTrigger(ui.SlideOverProps{Href: "/users/42/edit", Label: "Modifier", AcceptLocation: "/users/*", ReloadTarget: "#users"})

@ui.SlideOver("Modifier l'utilisateur") { ... }`,
			Note: "Le trigger utilise la syntaxe Unpoly actuelle new drawer et up-position=right ; la page du drawer reste une réponse serveur ordinaire.",
		},
		{
			Group: g, ID: "detailtabs", Title: "DetailTabs",
			Purpose: "Navigation entre sections d'une fiche, avec état actif porté par l'URL.", Uses: []string{"tabs", "up-follow"},
			Demo: demoDetailTabs(), Snippet: `@ui.DetailTabs(tabs, r.URL.Path, "#record-detail")`,
			Note: "Ce sont de vrais liens : historique navigateur, ouverture dans un nouvel onglet et fonctionnement sans JavaScript restent disponibles.",
		},
		{
			Group: g, ID: "bulkactionbar", Title: "BulkActionBar",
			Purpose: "Actions groupées visibles uniquement lorsqu'au moins une ligne est sélectionnée.", Uses: []string{"sticky", "aria-live"},
			Demo: demoBulkActionBar(), Snippet: `@ui.BulkActionBar(len(selected), "utilisateur sélectionné", "utilisateurs sélectionnés") { ... }`,
			Note: "La barre est annoncée comme statut dynamique et reste collée au bas de la zone pendant le défilement.",
		},
		{
			Group: g, ID: "authcard", Title: "AuthCard",
			Purpose: "Login concentré dans une card, soutenu par un hero produit sur grand écran.", Uses: []string{"card", "hero"}, StageClass: "gal-stage--app",
			Demo: demoTier2AuthCard(), Snippet: `@ui.AuthCard(ui.AuthCardProps{Product: "Acme Admin", Title: "Connexion", HeroTitle: "Pilotez vos opérations."}) { ... }`,
			Note: "Le hero disparaît sur mobile pour réserver l'espace à la tâche ; le nom du produit reste visible dans la card.",
		},
		{
			Group: g, ID: "filefield", Title: "FileField",
			Purpose: "Sélection de fichier, aide, erreur et progression d'envoi dans un seul fieldset.", Uses: []string{"file-input", "progress"},
			Demo: demoFileField(), Snippet: `@ui.FileField(ui.FileFieldProps{Name: "import", Label: "Fichier CSV", Accept: ".csv", Uploading: true, Progress: 64})`,
			Note: "Progress est borné entre 0 et 100. Le transport peut être natif, Unpoly ou direct-to-storage ; le composant ne couple pas l'UI à ce choix.",
		},
		{
			Group: g, ID: "audittimeline", Title: "AuditTimeline",
			Purpose: "Historique chronologique avec auteur, horodatage et nature de l'événement.", Uses: []string{"timeline", "time"},
			Demo: demoAuditTimeline(), Snippet: `@ui.AuditTimeline(events)`,
			Note: "La timeline reste compacte et sur un seul axe : un journal d'audit se parcourt mieux qu'une alternance gauche/droite.",
		},
		{
			Group: g, ID: "themetoggle", Title: "ThemeToggle",
			Purpose: "Basculer clair/sombre avec le theme-controller CSS de daisyUI.", Uses: []string{"theme-controller", "radio"},
			Demo: demoThemeToggle(), Snippet: `@ui.ThemeToggle(isDark)`,
			Note: "Deux radios, pas une case : daisyUI compile --prefersdark en @media (prefers-color-scheme: dark){ :root:not([data-theme]) }, donc une case unique n'exprime que « sombre » et « défaut » — et sur une machine en sombre système, « défaut » vaut sombre, ce qui rendait la position Clair inopérante. Chaque radio pose son :root:has(...:checked), plus spécifique que le média, donc les deux choix gagnent. Sans JavaScript : le choix survit aux échanges de fragments Unpoly, pas à un rechargement complet. La persistance reste à l'application.",
		},
		{
			Group: g, ID: "orgswitcher", Title: "OrgSwitcher",
			Purpose: "Changer de client sans quitter l'écran, et sans qu'un GET puisse le faire.", Uses: []string{"details", "formaction", "multi-tenant"},
			Demo: demoOrgSwitcher(), Snippet: `@ui.OrgSwitcher(ui.OrgSwitcherProps{
    ID: "org", Action: "/orgs/switch",
    Orgs: orgs, CurrentID: session.OrgID,
})`,
			Note: "Chaque entrée est un bouton de soumission : changer de tenant change tout l'écran, et un changement d'état ne passe pas par un lien. Volontairement sans up-submit — la barre haute et la barre latérale vivent hors de <main>, donc un échange de fragment laisserait l'ancienne organisation affichée dans l'en-tête. Une seule organisation n'est pas un choix : le nom s'affiche sans menu.",
		},
	}
}
