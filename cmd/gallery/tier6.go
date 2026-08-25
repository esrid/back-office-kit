package main

import "github.com/esrid/back-office-kit/ui"

func demoPaletteGroups() []ui.PaletteGroup {
	return []ui.PaletteGroup{
		{Label: "Navigation", Commands: []ui.PaletteCommand{
			{ID: "users", Label: "Ouvrir les utilisateurs", Description: "142 comptes actifs", Href: "/users", Shortcut: "G U"},
			{ID: "billing", Label: "Ouvrir la facturation", Description: "12 factures nécessitent une intervention", Href: "/billing", Shortcut: "G F"},
		}},
		{Label: "Actions", Commands: []ui.PaletteCommand{
			{ID: "new-user", Label: "Inviter un utilisateur", Description: "Ouvrir le formulaire d’invitation", Href: "/users/new", Shortcut: "N U"},
			{ID: "suspend", Label: "Suspendre Acme Caraïbes", Description: "Cette commande modifie le compte", Href: "/accounts/acme-caribbean/suspend", Method: "post", Confirm: "Suspendre ce compte ?", Tone: "error"},
			{ID: "export", Label: "Exporter toutes les données", Description: "Votre rôle ne permet pas cette action", Disabled: true},
		}},
	}
}

func demoQueryFields() []ui.Option {
	return []ui.Option{{Value: "status", Label: "Statut"}, {Value: "amount", Label: "Montant dû"}, {Value: "country", Label: "Pays"}, {Value: "overdue_days", Label: "Jours de retard"}}
}

func demoQueryOperators() []ui.Option {
	return []ui.Option{{Value: "eq", Label: "est"}, {Value: "neq", Label: "n’est pas"}, {Value: "gt", Label: "est supérieur à"}, {Value: "contains", Label: "contient"}}
}

func demoQueryRoot() ui.QueryGroup {
	return ui.QueryGroup{
		ID: "root", Join: ui.QueryAll, AllowAddRule: true, AllowAddGroup: true,
		Rules: []ui.QueryRule{{
			ID: "status", Field: "status", Operator: "eq", Value: "overdue",
			ValueOptions: []ui.Option{{Value: "active", Label: "Actif"}, {Value: "overdue", Label: "En retard"}, {Value: "suspended", Label: "Suspendu"}},
		}},
		Groups: []ui.QueryGroup{{
			ID: "priority", Join: ui.QueryAny, AllowAddRule: true, AllowRemove: true,
			Rules: []ui.QueryRule{
				{ID: "amount", Field: "amount", Operator: "gt", Value: "10000", ValueType: "number"},
				{ID: "country", Field: "country", Operator: "eq", Value: "MQ", ValueOptions: []ui.Option{{Value: "FR", Label: "France"}, {Value: "GP", Label: "Guadeloupe"}, {Value: "MQ", Label: "Martinique"}}},
			},
		}},
	}
}

func demoManagedColumns() []ui.ManagedColumn {
	return []ui.ManagedColumn{
		{Key: "name", Label: "Nom", Visible: true, Locked: true},
		{Key: "company", Label: "Entreprise", Visible: true},
		{Key: "email", Label: "E-mail", Visible: true},
		{Key: "status", Label: "Statut", Visible: true},
		{Key: "created_at", Label: "Date de création", Visible: false},
		{Key: "last_seen", Label: "Dernière activité", Visible: false},
	}
}

func demoMasterDetailItems() []ui.MasterDetailItem {
	return []ui.MasterDetailItem{
		{ID: "acme-caribbean", Title: "Acme Caraïbes", Description: "12 factures · Fort-de-France", Meta: "12", Href: "/accounts/acme-caribbean"},
		{ID: "kalina", Title: "Kalina Conseil", Description: "3 factures · Le Lamentin", Meta: "3", Href: "/accounts/kalina"},
		{ID: "madinina", Title: "Madinina Services", Description: "1 facture · Schoelcher", Meta: "1", Href: "/accounts/madinina"},
		{ID: "bwa-galba", Title: "Bwa Galba", Description: "7 factures · Sainte-Marie", Meta: "7", Href: "/accounts/bwa-galba"},
	}
}

func tier6Sections() []Section {
	return []Section{
		{
			Group: "Tier 6 — Utilisateurs avancés", ID: "command-palette", Title: "CommandPalette",
			Purpose: "Accéder rapidement aux pages et actions sans transformer les mutations en liens.",
			Uses:    []string{"search", "modal", "safe commands"}, Demo: demoCommandPalette(),
			Snippet: `@ui.CommandPalette(props) · @ui.CommandPaletteTrigger(trigger)`,
			Note:    "La recherche est un formulaire serveur ordinaire. Les navigations restent des liens ; toute commande avec Method devient un bouton relié à un formulaire POST.",
		},
		{
			Group: "Tier 6 — Utilisateurs avancés", ID: "query-builder", Title: "QueryBuilder",
			Purpose: "Composer des règles AND/OR imbriquées dans une URL interprétable par le backend.",
			Uses:    []string{"nested fieldsets", "GET", "query groups"}, Demo: demoQueryBuilder(),
			Snippet: `@ui.QueryBuilder(ui.QueryBuilderProps{Fields: fields, Operators: operators, Root: query})`,
			Note:    "Le composant n’évalue jamais les règles. Le backend valide champs et opérateurs autorisés, reconstruit le brouillon après add/remove et applique la requête finale.",
		},
		{
			Group: "Tier 6 — Utilisateurs avancés", ID: "column-manager", Title: "ColumnManager",
			Purpose: "Configurer visibilité, ordre et densité d’un tableau sans dépendre du drag-and-drop.",
			Uses:    []string{"checkboxes", "order", "density"}, Demo: demoColumnManager(),
			Snippet: `@ui.ColumnManager(ui.ColumnManagerProps{Columns: columns, Density: density, Action: "/views/columns"})`,
			Note:    "Monter/descendre reste utilisable au clavier et sans JavaScript. Une colonne obligatoire est soumise par champ caché même si sa checkbox est désactivée.",
		},
		{
			Group: "Tier 6 — Utilisateurs avancés", ID: "master-detail", Title: "MasterDetail · MasterDetailNav",
			Purpose: "Traiter une liste et son record actif dans un espace dense qui se replie sur mobile.",
			Uses:    []string{"responsive grid", "URL selection", "fragment target"}, Demo: demoMasterDetail(),
			Snippet: `@ui.MasterDetail(props) { @accountDetail(selected) }`,
			Note:    "La sélection reste une URL normale. Unpoly peut ne remplacer que la région detail ; sur écran étroit, la liste et la fiche s’empilent naturellement.",
		},
	}
}
