package main

func tier9Sections() []Section {
	const g = "Tier 9 — Collaboration humaine"
	return []Section{
		{
			Group: g, ID: "assignee-picker", Title: "AssigneePicker",
			Purpose: "Attribuer une fiche à un opérateur sans écraser une affectation concurrente.",
			Uses:    []string{"fieldset", "record version", "POST"}, Demo: demoAssigneePicker(),
			Snippet: `@ui.AssigneePicker(ui.AssigneePickerProps{
    ID: "assignee", Action: "/records/42/assignee",
    Version: record.Version, Current: record.AssigneeID,
    Options: operators,
})`,
			Note: "L’affectation transporte la version de la fiche et reste un fragment mineur : le backend refuse une version obsolète, puis renvoie le même contrôle avec la valeur autoritaire.",
		},
		{
			Group: g, ID: "comment-thread", Title: "CommentThread",
			Purpose: "Conserver le contexte humain sur une fiche et publier une note sans rechargement.",
			Uses:    []string{"plain text", "thread version", "up-poll"}, Demo: demoCommentThread(),
			Snippet: `@ui.CommentThread(ui.CommentThreadProps{
    ID: "discussion", Action: "/records/42/comments",
    Version: thread.Version, Comments: thread.Items,
})`,
			Note: "Les commentaires restent du texte brut. Le serveur résout les mentions, crée les notifications et rejette une version obsolète ; aucun HTML fourni par un opérateur n’est réinjecté.",
		},
		{
			Group: g, ID: "watch-button", Title: "WatchButton",
			Purpose: "Suivre une fiche et rendre visible le nombre d’opérateurs abonnés.",
			Uses:    []string{"aria-pressed", "record_id", "POST"}, Demo: demoWatchButton(),
			Snippet: `@ui.WatchButton(ui.WatchButtonProps{
    Action: "/records/42/watch", RecordID: record.ID,
    Version: record.Version, Watching: record.Watching,
    Count: record.WatcherCount,
})`,
			Note: "Suivre et ne plus suivre sont deux intentions POST explicites. Le compteur vient du backend et ne dérive jamais d’un état optimiste caché dans le navigateur.",
		},
		{
			Group: g, ID: "activity-feed", Title: "ActivityFeed",
			Purpose: "Réunir commentaires, affectations et transitions métier dans un fil lisible.",
			Uses:    []string{"timeline", "up-poll", "up-history"}, Demo: demoActivityFeed(),
			Snippet: `@ui.ActivityFeed(ui.ActivityFeedProps{
    ID: "activity", Items: record.Activity,
    DetailTarget: "#record-detail",
})`,
			Note: "Ce fil aide à reprendre le travail ; ce n’est pas une preuve immuable. Un lien de détail ciblé porte up-history=true puisqu’il ouvre une navigation partageable. AuditTimeline reste la surface d’audit.",
		},
		{
			Group: g, ID: "deadline-status", Title: "DeadlineStatus",
			Purpose: "Afficher une échéance ou un SLA dans le fuseau de l’opérateur avec un état mesurable.",
			Uses:    []string{"time.Time", "explicit Now", "timezone"}, Demo: demoDeadlineStatus(),
			Snippet: `@ui.DeadlineStatus(ui.DeadlineStatusProps{
    Label: "SLA de réponse", At: ticket.DueAt,
    Now: requestNow, Location: operator.Location,
    WarningBefore: 4 * time.Hour,
})`,
			Note: "Now est injecté par le handler : le composant ne lit jamais l’horloge. Le seuil d’avertissement est métier et explicite ; la date absolue affiche toujours le fuseau.",
		},
		{
			Group: g, ID: "task-list", Title: "TaskList",
			Purpose: "Suivre des relances assignées, datées et terminables sans état client caché.",
			Uses:    []string{"task version", "exact intent", "DeadlineStatus"}, Demo: demoTaskList(),
			Snippet: `@ui.TaskList(ui.TaskListProps{
    ID: "tasks", Action: "/records/42/tasks",
    Version: tasks.Version, Items: tasks.Items,
    Now: requestNow, Location: operator.Location,
})`,
			Note: "Chaque bouton envoie complete:<task_id> ou reopen:<task_id> avec la version de la liste. Le backend refuse une transition obsolète et décide si la tâche reste modifiable.",
		},
		{
			Group: g, ID: "record-presence", Title: "RecordPresence",
			Purpose: "Signaler qui consulte ou modifie la même fiche sans promettre un verrou inexistant.",
			Uses:    []string{"advisory presence", "aria-live", "up-poll"}, Demo: demoRecordPresence(),
			Snippet: `@ui.RecordPresence(ui.RecordPresenceProps{
    ID: "presence", People: presence.People,
    Busy: true, PollSource: "/records/42/presence",
})`,
			Note: "La présence est une indication périssable, jamais une garantie d’exclusivité. Chaque sauvegarde doit encore transporter une version et le backend doit détecter les conflits.",
		},
		{
			Group: g, ID: "collaboration-states", Title: "États de collaboration",
			Purpose: "Vérifier visuellement conflit, désactivation, vide et actualisation avant l’intégration.",
			Uses:    []string{"error", "disabled", "empty", "busy"}, Demo: demoCollaborationStates(),
			Snippet: `// Les erreurs et états viennent du serveur.
@ui.AssigneePicker(ui.AssigneePickerProps{Error: form.Error, ...})
@ui.RecordPresence(ui.RecordPresenceProps{Busy: loading, ...})`,
			Note: "Une galerie limitée au happy path laisse les états silencieux arriver pour la première fois en production. Cette matrice garde visibles les contrats de validation, d’indisponibilité et de polling.",
		},
	}
}
