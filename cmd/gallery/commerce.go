package main

func commerceSections() []Section {
	const g = "Commerce — back-office"
	return []Section{
		{
			Group: g, ID: "money", Title: "Money",
			Purpose: "Le pendant de Timestamp pour l'argent : une seule façon d'écrire un montant.",
			Uses:    []string{"tabular-nums", "int64"},
			Demo:    demoMoney(),
			Snippet: `@ui.Money(ui.MoneyProps{Cents: order.TotalCents, Currency: "EUR"})

ui.FormatMoney(-4500, "EUR")   // "−45,00 €"`,
			Note: "Cents est un entier en unités mineures, jamais un float : 0,1 + 0,2 ne fait pas 0,3 en binaire, et de l'argent qui dérive d'un centime par opération est le plus vieux bug du commerce. Le nombre de décimales suit la devise — supposer deux transforme 1000 ¥ en 10,00 ¥. Le signe négatif est un vrai moins U+2212, pas un trait d'union : il s'aligne avec les chiffres. L'espace avant le symbole est insécable, « € » ne part jamais seul à la ligne.",
		},
		{
			Group: g, ID: "lineitems", Title: "LineItemsTable",
			Purpose: "Une commande ou une facture, avec ses ajustements et son total.",
			Uses:    []string{"table", "tfoot", "Money"},
			Demo:    demoLineItems(),
			Snippet: `@ui.LineItemsTable(ui.LineItemsProps{
    Currency: "EUR", Items: order.Lines,
    Adjustments: order.Adjustments, TotalCents: order.TotalCents,
})`,
			Note: "Le total affiché est celui du serveur, jamais une addition faite ici : remises par ligne, prorata et règles d'arrondi vivent dans le code de facturation. Le composant additionne quand même les lignes — mais seulement pour vérifier.",
		},
		{
			Group: g, ID: "lineitems-mismatch", Title: "LineItemsTable · écart détecté",
			Purpose: "Ce qui arrive quand le total contredit ses propres lignes.",
			Uses:    []string{"alert"},
			Demo:    demoLineItemsMismatch(),
			Snippet: `// mêmes props, mais TotalCents ne correspond pas aux lignes`,
			Note:    "Afficher en silence un total qui contredit son propre détail, c'est facturer le mauvais montant sans que personne ne le voie. L'écart est nommé et chiffré, et le total du serveur reste affiché — le composant n'en invente jamais un autre. Un test couvre les deux cas.",
		},
		{
			Group: g, ID: "address", Title: "AddressBlock",
			Purpose: "Une adresse dans l'ordre attendu par son pays de destination.",
			Uses:    []string{"address"},
			Demo:    demoAddress(),
			Snippet: `@ui.AddressBlock("Livraison", ui.Address{
    Name: o.Name, Line1: o.Line1, PostalCode: o.Zip, City: o.City,
    Country: o.Country, CountryISO: o.CountryISO,
})`,
			Note: "L'ordre des lignes change selon le pays : code postal avant la ville en Europe continentale, « ville, région code » aux États-Unis, code postal sur sa propre ligne après la ville au Royaume-Uni. Écrire l'adresse à la main dans un template donne des étiquettes que le transporteur ne sait pas lire. Les parties vides sont supprimées plutôt que laissées en lignes blanches.",
		},
		{
			Group: g, ID: "refund", Title: "RefundForm",
			Purpose: "Rembourser tout ou partie, borné par ce qui reste réellement remboursable.",
			Uses:    []string{"Money", "fieldset"},
			Demo:    demoRefund(),
			Snippet: `@ui.RefundForm(ui.RefundFormProps{
    OrderID: o.ID, Currency: "EUR",
    TotalCents: o.TotalCents, AlreadyRefundedCents: o.RefundedCents,
    Action: "/orders/" + o.ID + "/refund", Reasons: motifs,
})`,
			Note: "AlreadyRefundedCents vient du serveur : le navigateur ne décide jamais de ce qui reste, et un second onglet ne doit pas pouvoir rembourser deux fois la même somme. Le champ propose le reste, pas le total, et si le serveur annonce plus remboursé que facturé le formulaire affiche zéro au lieu d'un montant négatif. Une commande intégralement remboursée n'offre aucun champ — testé.",
		},
	}
}

func marketingSections() []Section {
	const g = "Marketing — pages publiques"
	const full = "gal-stage--app"
	return []Section{
		{
			Group: g, ID: "hero", Title: "Hero", StageClass: full,
			Purpose: "Le seul <h1> de la page.",
			Uses:    []string{"h1", "img"},
			Demo:    demoHero(),
			Snippet: `@marketing.Hero(marketing.HeroProps{
    Title: "Pilotez vos opérations sans perdre le fil",
    Actions: []marketing.Action{{Label: "Essayer", Href: "/signup", Primary: true}},
    Visual: marketing.Image{Src: "/hero.png", Alt: "Aperçu", Width: 480, Height: 300},
})`,
			Note: "Un seul h1 par page, et c'est Hero qui le porte ; toutes les autres sections ouvrent en h2. Une hiérarchie de titres fausse ment aux moteurs de recherche comme aux lecteurs d'écran, et c'est invisible à l'œil. Un test le vérifie sur les cinq sections.",
		},
		{
			Group: g, ID: "features", Title: "FeatureGrid", StageClass: full,
			Purpose: "Trois à six arguments, pas vingt.",
			Uses:    []string{"grid", "h2"},
			Demo:    demoFeatures(),
			Snippet: `@marketing.FeatureGrid("Ce que vous y gagnez", sousTitre, features)`,
			Note:    "text-pretty sur les descriptions et text-balance sur les titres : les lignes orphelines d'un mot sont la marque des pages générées à la va-vite.",
		},
		{
			Group: g, ID: "pricing", Title: "PricingTable", StageClass: full,
			Purpose: "Les paliers, dans un ordre qui veut dire quelque chose.",
			Uses:    []string{"card", "ol"},
			Demo:    demoPricing(),
			Snippet: `@marketing.PricingTable("Tarifs", sousTitre, []marketing.Plan{
    {Name: "Équipe", Price: "89 €", Period: "par mois", Highlighted: true, Badge: "Populaire"},
})`,
			Note: "Une <ol> et non une <div> : l'ordre des paliers porte du sens, et un lecteur d'écran annonce « 2 sur 3 ». Le prix est une chaîne, pas un ui.Money — un prix d'affichage n'est pas un montant calculé, et « Sur devis » doit rester possible.",
		},
		{
			Group: g, ID: "testimonials", Title: "TestimonialGrid", StageClass: full,
			Purpose: "Des citations attribuées, pas des phrases flottantes.",
			Uses:    []string{"blockquote", "cite"},
			Demo:    demoTestimonials(),
			Snippet: `@marketing.TestimonialGrid("Ils l'utilisent tous les jours", quotes)`,
			Note:    "blockquote, figcaption et cite : l'attribution est une donnée, pas une décoration. Sans ce balisage, un lecteur d'écran ne sait pas où finit la citation ni qui l'a dite.",
		},
		{
			Group: g, ID: "faq", Title: "FAQ", StageClass: full,
			Purpose: "Un accordéon sans une ligne de JavaScript.",
			Uses:    []string{"details", "collapse"},
			Demo:    demoFAQ(),
			Snippet: `@marketing.FAQ("Questions fréquentes", []marketing.QA{
    {Question: "Puis-je changer de palier ?", Answer: "Oui, au prorata."},
})`,
			Note: "details/summary : la recherche du navigateur atteint les réponses fermées, ce qu'un accordéon scripté empêche. Le name partagé rend l'accordéon exclusif là où le navigateur le gère, et laisse simplement plusieurs réponses ouvertes ailleurs — dégradation sans casse.",
		},
		{
			Group: g, ID: "cta", Title: "CTASection", StageClass: full,
			Purpose: "La dernière chose avant le pied de page.",
			Uses:    []string{"btn"},
			Demo:    demoCTA(),
			Snippet: `@marketing.CTASection(titre, sousTitre, actions, "Aucune carte bancaire demandée.")`,
			Note:    "La note sous les boutons lève l'objection au moment où elle se pose, pas dans une FAQ trois écrans plus haut.",
		},
		{
			Group: g, ID: "logos", Title: "LogoCloud", StageClass: full,
			Purpose: "Qui l'utilise déjà.",
			Uses:    []string{"img", "alt"},
			Demo:    demoLogos(),
			Snippet: `@marketing.LogoCloud("Ils nous font confiance", logos)`,
			Note:    "Chaque logo garde son texte alternatif : un mur d'images non décrites ne dit rien à un lecteur d'écran, et les noms de clients sont précisément ce qu'on veut voir indexé.",
		},
		{
			Group: g, ID: "marketing-footer", Title: "MarketingFooter", StageClass: full,
			Purpose: "Colonnes de liens et mentions légales.",
			Uses:    []string{"footer", "nav"},
			Demo:    demoFooter(),
			Snippet: `@marketing.MarketingFooter(produit, accroche, colonnes, "© 2026 Acme SAS")`,
			Note:    "Chaque colonne est un <nav> avec son propre aria-label : un lecteur d'écran peut sauter de groupe en groupe au lieu de traverser trente liens.",
		},
	}
}
