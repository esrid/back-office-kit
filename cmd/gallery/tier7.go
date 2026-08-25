package main

func tier7Sections() []Section {
	const g = "Tier 7 — Exploitation"
	return []Section{
		{
			Group: g, ID: "timestamp", Title: "Timestamp",
			Purpose: "Un instant, dit une seule fois, avec son fuseau assumé.",
			Uses:    []string{"time", "datetime"},
			Demo:    demoTimestamp(),
			Snippet: `@ui.Timestamp(ui.TimestampProps{
    At: event.CreatedAt, Now: req.Time, Location: operator.Zone,
})`,
			Note: "Now est un paramètre, jamais time.Now() : un composant qui lit l'horloge rend différemment à chaque appel et devient intestable. Le fuseau est toujours écrit, jamais sous-entendu — un opérateur à Paris qui lit un journal en UTC prend des décisions fausses avec assurance. La phrase relative et la date complète coexistent : l'une est affichée, l'autre est dans le title, et datetime reste lisible par une machine.",
		},
		{
			Group: g, ID: "sparkline", Title: "Sparkline",
			Purpose: "La tendance à côté du nombre, parce qu'un nombre seul ment.",
			Uses:    []string{"svg", "role=img"},
			Demo:    demoSparkline(),
			Snippet: `@ui.Sparkline(ui.SparklineProps{
    Values: revenue, Label: "Revenu mensuel", Suffix: " k€",
})`,
			Note: "SVG rendu côté serveur : aucune bibliothèque graphique, aucun état client, rien à hydrater. Une seule série donc pas de légende, et l'encre reste neutre par défaut — le chiffre voisin porte déjà bon ou mauvais, et les couleurs de statut sont réservées au statut. Le <title> tient lieu d'infobulle sans JavaScript et l'aria-label énonce minimum, maximum et valeur courante. Une série plate se pose sur la ligne médiane au lieu de diviser par un intervalle nul.",
		},
		{
			Group: g, ID: "undo", Title: "UndoSnackbar",
			Purpose: "Le filet des actions de masse, qui manquait au kit.",
			Uses:    []string{"toast", "alert"},
			Demo:    demoUndoSnackbar(),
			Snippet: `@ui.UndoSnackbar(ui.UndoProps{
    Message: "12 comptes suspendus", Action: "/undo",
    ReceiptID: receipt.ID, ExpiresIn: "30 secondes",
})`,
			Note: "BulkActionBar peut suspendre deux cents comptes d'un clic ; rien ne les rattrapait. Ce n'est pas une confirmation : l'action a déjà eu lieu. L'annulation vise un ReceiptID précis, jamais « la dernière action » — deux onglets ouverts et « la dernière » ne veut plus rien dire. C'est un formulaire POST, pas un lien.",
		},
		{
			Group: g, ID: "impersonation", Title: "ImpersonationBanner",
			Purpose: "Agir au nom de quelqu'un sans jamais l'oublier.",
			Uses:    []string{"warning", "sticky"},
			Demo:    demoImpersonation(),
			Snippet: `@ui.ImpersonationBanner(ui.ImpersonationProps{
    Subject: target.Email, Actor: session.Email,
    Since: "8 minutes", Action: "/impersonation/stop",
})`,
			Note: "Le danger n'est pas d'entrer en usurpation, c'est d'oublier qu'on y est. D'où le bandeau collant pleine largeur sur la surface warning plutôt qu'une pastille discrète, et la sortie contenue dans le bandeau lui-même. L'opérateur réel est nommé : ce qui est écrit à l'écran doit correspondre à ce qui sera écrit au journal d'audit.",
		},
		{
			Group: g, ID: "errorpage", Title: "ErrorPage",
			Purpose:    "Le seul écran certain d'être vu, et celui qu'on ne dessine jamais.",
			Uses:       []string{"badge", "hors Shell"},
			StageClass: "gal-stage--app",
			Demo:       demoErrorPage(),
			Snippet: `@ui.ErrorPage(ui.ErrorPageProps{
    Code: "500", Title: "Le rapport n'a pas pu être généré",
    RetryHref: r.URL.Path, RequestID: middleware.RequestID(ctx),
})`,
			Note: "Rendu hors du Shell : ce qui a échoué peut être le Shell. Le RequestID donne à l'opérateur quelque chose de précis à citer au support, au lieu de « ça marche pas ». Un 404 reste neutre, un 500 est une erreur, un 403 un avertissement — trois situations différentes ne méritent pas la même couleur.",
		},
		{
			Group: g, ID: "importmapper", Title: "ImportMapper",
			Purpose: "L'étape où un humain décide que « Email address » veut dire e-mail.",
			Uses:    []string{"table", "select", "up-validate"},
			Demo:    demoImportMapper(),
			Snippet: `@ui.ImportMapper(ui.ImportMapperProps{
    ID: "import-users", Action: "/imports/42/run",
    Fields: schema, Columns: detected, RowCount: 1284,
})`,
			Note: "JobProgress couvre l'exécution de l'import ; c'est ici que les données se cassent vraiment. On associe sur pièces : les premières valeurs de chaque colonne sont affichées à côté du choix. Deux garde-fous testés — un champ obligatoire sans colonne désactive le lancement en disant pourquoi, et deux colonnes visant le même champ sont signalées puisque la dernière écraserait les autres en silence.",
		},
	}
}
