package main

func tier8Sections() []Section {
	const g = "Tier 8 — Sécurité & accès"
	return []Section{
		{
			Group: g, ID: "permission-matrix", Title: "PermissionMatrix",
			Purpose: "Administrer les capacités d’un rôle sans confondre affichage et autorisation.",
			Uses:    []string{"table", "role version", "inherited grants"}, Demo: demoPermissionMatrix(),
			Snippet: `@ui.PermissionMatrix(ui.PermissionMatrixProps{
    RoleID: role.ID, RoleVersion: role.Version,
    Columns: capabilities, Rows: resources,
})`,
			Note: "Le backend fournit chaque cellule, son origine et son caractère modifiable. Les permissions héritées cochées sont également soumises par champ caché, mais le backend recalcule toujours les droits effectifs et refuse une version obsolète.",
		},
		{
			Group: g, ID: "access-grants", Title: "AccessGrantList",
			Purpose: "Voir et révoquer les accès directs par identifiant de grant.",
			Uses:    []string{"grant_id", "scope", "expiry"}, Demo: demoAccessGrants(),
			Snippet: `@ui.AccessGrantList(ui.AccessGrantListProps{
    Action: "/access/grants", Version: grants.Version,
    Items: grants.Items,
})`,
			Note: "La liste ne déduit jamais un accès effectif. Révoquer envoie grant_version et intent=revoke:<grant_id> ; un accès produit par une politique reste visible mais non révocable depuis cette surface.",
		},
		{
			Group: g, ID: "session-manager", Title: "SessionManager",
			Purpose: "Identifier les appareils connectés et révoquer une session précise.",
			Uses:    []string{"session_id", "current session", "POST"}, Demo: demoSessionManager(),
			Snippet: `@ui.SessionManager(ui.SessionManagerProps{
    SubjectID: user.ID, Version: sessions.Version,
    Sessions: sessions.Items, Action: "/security/sessions",
})`,
			Note: "La session courante est nommée et ne peut pas être révoquée par erreur. Une révocation globale signifie toujours « les autres sessions » et reste une mutation POST décidée par le backend.",
		},
		{
			Group: g, ID: "api-key-manager", Title: "APIKeyManager",
			Purpose: "Gérer des clés masquées, leurs portées, expiration et dernière utilisation.",
			Uses:    []string{"masked prefix", "scopes", "key version"}, Demo: demoAPIKeyManager(),
			Snippet: `@ui.APIKeyManager(ui.APIKeyManagerProps{
    Keys: keys.Items, Version: keys.Version,
    Action: "/security/api-keys",
})`,
			Note: "Ce composant ne reçoit jamais le secret complet. Le préfixe sert à reconnaître la clé ; l’identifiant opaque sert à la révoquer après une nouvelle vérification backend.",
		},
		{
			Group: g, ID: "secret-reveal", Title: "SecretReveal",
			Purpose: "Remettre un secret une seule fois puis enregistrer son acquittement.",
			Uses:    []string{"one-time secret", "readonly", "acknowledgement"}, Demo: demoSecretReveal(),
			Snippet: `@ui.SecretReveal(ui.SecretRevealProps{
    Secret: created.Plaintext, SecretID: created.ID,
    AcknowledgeAction: "/security/api-keys/ack",
})`,
			Note: "Le secret apparaît comme valeur visible sélectionnable, jamais dans data-* ou une URL. Après l’acquittement POST, le serveur détruit la copie affichable et ne rend plus ce composant.",
		},
		{
			Group: g, ID: "step-up-auth", Title: "StepUpAuthCard",
			Purpose: "Demander un facteur frais avant de reprendre une action sensible.",
			Uses:    []string{"challenge_id", "expiry", "fresh factor"}, Demo: demoStepUpAuth(),
			Snippet: `@ui.StepUpAuthCard(ui.StepUpAuthProps{
    ChallengeID: challenge.ID, ChallengeVersion: challenge.Version,
    Method: challenge.Method, Action: "/auth/step-up",
})`,
			Note: "Réussir la vérification ne donne aucun droit nouveau : le backend reprend ensuite l’action en attente, revérifie permission, cible et contexte, puis consomme le challenge.",
		},
		{
			Group: g, ID: "policy-simulator", Title: "PolicySimulator",
			Purpose: "Expliquer une décision acteur × ressource × action sans modifier les accès.",
			Uses:    []string{"GET", "decision_id", "policy facts"}, Demo: demoPolicySimulator(),
			Snippet: `@ui.PolicySimulator(ui.PolicySimulatorProps{
    Actors: actors, Resources: resources, Actions: actions,
    Simulation: decision,
})`,
			Note: "La requête est un GET vers le moteur de politique. Le résultat porte une référence auditable et rappelle explicitement qu’aucun accès n’a été accordé.",
		},
		{
			Group: g, ID: "security-event-feed", Title: "SecurityEventFeed",
			Purpose: "Surveiller les événements sensibles avec acteur, cible et référence de requête.",
			Uses:    []string{"audit projection", "severity", "up-poll"}, Demo: demoSecurityEventFeed(),
			Snippet: `@ui.SecurityEventFeed(ui.SecurityEventFeedProps{
    Events: events, Filters: filters, PollSource: "/security/events",
})`,
			Note: "Le flux reçoit une projection déjà filtrée par le backend. Il n’expose ni secret ni charge utile brute ; le détail lié applique encore les mêmes politiques d’accès.",
		},
	}
}
