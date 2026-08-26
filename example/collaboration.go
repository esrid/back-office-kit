package main

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/esrid/back-office-kit/ui"
)

type CollaborationView struct {
	Operator      string
	RecordVersion string
	ThreadVersion string
	TaskVersion   string
	AssigneeID    string
	AssigneeError string
	CommentError  string
	TaskError     string
	Watching      bool
	WatcherCount  int
	Comments      []ui.Comment
	Tasks         []ui.TaskItem
	Activity      []ui.ActivityItem
	People        []ui.PresencePerson
	Now           time.Time
	DueAt         time.Time
	Location      *time.Location
	Flashes       []ui.Flash
}

type collaborationState struct {
	mu            sync.Mutex
	recordVersion int
	threadVersion int
	taskVersion   int
	assigneeID    string
	watching      bool
	watcherCount  int
	comments      []ui.Comment
	tasks         []ui.TaskItem
	now           time.Time
	location      *time.Location
}

func newCollaborationState() *collaborationState {
	location := time.FixedZone("AST", -4*60*60)
	now := time.Date(2026, 8, 25, 15, 20, 0, 0, location)
	return &collaborationState{
		recordVersion: 18,
		threadVersion: 7,
		taskVersion:   11,
		assigneeID:    "usr_amelie",
		watching:      true,
		watcherCount:  4,
		now:           now,
		location:      location,
		comments: []ui.Comment{
			{ID: "comment_71", Author: "Nora Martin", Initials: "NM", Body: "Le budget est confirmé. Il manque encore la date de déploiement.", At: "hier à 16:42", Datetime: "2026-08-24T16:42:00-04:00"},
			{ID: "comment_74", Author: "Amélie Rousseau", Initials: "AR", Body: "J’appelle le client cet après-midi et je complète la fiche.", At: "il y a 28 minutes", Datetime: "2026-08-25T14:52:00-04:00"},
		},
		tasks: []ui.TaskItem{
			{ID: "task_81", Title: "Confirmer la date de déploiement", Assignee: "Amélie Rousseau", DueAt: now.Add(2 * time.Hour)},
			{ID: "task_82", Title: "Envoyer la proposition révisée", Assignee: "Nora Martin", DueAt: now.Add(26 * time.Hour)},
			{ID: "task_79", Title: "Valider le budget", Assignee: "Bruno Keller", DueAt: now.Add(-24 * time.Hour), Completed: true},
		},
	}
}

func collaborationAssignees() []ui.Option {
	return []ui.Option{
		{Value: "usr_amelie", Label: "Amélie Rousseau"},
		{Value: "usr_nora", Label: "Nora Martin"},
		{Value: "usr_bruno", Label: "Bruno Keller"},
	}
}

func collaborationActivities() []ui.ActivityItem {
	return []ui.ActivityItem{
		{ID: "activity_91", Kind: ui.ActivityComment, Actor: "Amélie Rousseau", Action: "a ajouté un commentaire", Detail: "La date de déploiement sera confirmée vendredi.", At: "il y a 28 minutes", Datetime: "2026-08-25T14:52:00-04:00"},
		{ID: "activity_88", Kind: ui.ActivityAssignment, Actor: "Nora Martin", Action: "a confié la fiche à Amélie Rousseau", At: "hier à 09:14", Datetime: "2026-08-24T09:14:00-04:00"},
		{ID: "activity_82", Kind: ui.ActivityStatus, Actor: "Bruno Keller", Action: "a passé l’opportunité à Proposition", Detail: "Valeur estimée : 24 000 €.", At: "le 22 août", Datetime: "2026-08-22T11:06:00-04:00"},
	}
}

func collaborationPresence() []ui.PresencePerson {
	return []ui.PresencePerson{
		{ID: "usr_nora", Name: "Nora Martin", Initials: "NM", Mode: ui.PresenceEditing},
		{ID: "usr_bruno", Name: "Bruno Keller", Initials: "BK", Mode: ui.PresenceViewing},
	}
}

func collaborationFlashes(q url.Values) []ui.Flash {
	switch q.Get("notice") {
	case "assigned":
		return []ui.Flash{{Tone: "success", Text: "Responsable mis à jour."}}
	case "watched":
		return []ui.Flash{{Tone: "success", Text: "Abonnement mis à jour."}}
	case "commented":
		return []ui.Flash{{Tone: "success", Text: "Commentaire publié."}}
	case "task":
		return []ui.Flash{{Tone: "success", Text: "Tâche mise à jour."}}
	}
	if q.Get("conflict") != "" {
		return []ui.Flash{{Tone: "error", Text: "La fiche a changé entre-temps. Vérifiez la version actuelle avant de réessayer."}}
	}
	return nil
}

func (s *collaborationState) view(operator string, q url.Values) CollaborationView {
	s.mu.Lock()
	defer s.mu.Unlock()

	v := CollaborationView{
		Operator:      operator,
		RecordVersion: strconv.Itoa(s.recordVersion), ThreadVersion: strconv.Itoa(s.threadVersion), TaskVersion: strconv.Itoa(s.taskVersion),
		AssigneeID: s.assigneeID, Watching: s.watching, WatcherCount: s.watcherCount,
		Comments: append([]ui.Comment(nil), s.comments...), Tasks: append([]ui.TaskItem(nil), s.tasks...),
		Activity: collaborationActivities(), People: collaborationPresence(),
		Now: s.now, DueAt: s.now.Add(2 * time.Hour), Location: s.location,
		Flashes: collaborationFlashes(q),
	}
	switch q.Get("conflict") {
	case "assignment":
		v.AssigneeError = "Cette affectation partait d’une ancienne version. La valeur actuelle a été rechargée."
	case "comment":
		v.CommentError = "La discussion a changé. Relisez les nouveaux commentaires avant de publier."
	case "task":
		v.TaskError = "La liste a changé. Vérifiez l’état actuel des tâches avant de réessayer."
	}
	switch q.Get("error") {
	case "assignment":
		v.AssigneeError = "Choisissez un responsable connu."
	case "comment":
		v.CommentError = "Le commentaire ne peut pas être vide."
	case "task":
		v.TaskError = "Cette tâche ou cette transition n’existe plus."
	}
	return v
}

func collaborationRedirect(w http.ResponseWriter, r *http.Request, key, value string) {
	q := url.Values{}
	if key != "" {
		q.Set(key, value)
	}
	destination := "/collaboration"
	if encoded := q.Encode(); encoded != "" {
		destination += "?" + encoded
	}
	http.Redirect(w, r, destination, http.StatusSeeOther)
}

func validAssignee(id string) bool {
	if id == "" {
		return true
	}
	for _, option := range collaborationAssignees() {
		if option.Value == id {
			return true
		}
	}
	return false
}

func registerCollaboration(mux *http.ServeMux, operator string) *collaborationState {
	state := newCollaborationState()

	mux.HandleFunc("GET /collaboration", func(w http.ResponseWriter, r *http.Request) {
		render(r.Context(), w, CollaborationPage(state.view(operator, r.URL.Query())))
	})
	mux.HandleFunc("GET /collaboration/presence", func(w http.ResponseWriter, r *http.Request) {
		render(r.Context(), w, ui.RecordPresence(ui.RecordPresenceProps{ID: "collaboration-presence", People: collaborationPresence()}))
	})
	mux.HandleFunc("GET /collaboration/comments", func(w http.ResponseWriter, r *http.Request) {
		v := state.view(operator, nil)
		render(r.Context(), w, collaborationCommentThread(v))
	})

	mux.HandleFunc("POST /collaboration/assignee", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulaire invalide", http.StatusBadRequest)
			return
		}
		state.mu.Lock()
		if r.Form.Get("record_version") != strconv.Itoa(state.recordVersion) {
			state.mu.Unlock()
			collaborationRedirect(w, r, "conflict", "assignment")
			return
		}
		assignee := r.Form.Get("assignee_id")
		if !validAssignee(assignee) {
			state.mu.Unlock()
			collaborationRedirect(w, r, "error", "assignment")
			return
		}
		state.assigneeID = assignee
		state.recordVersion++
		state.mu.Unlock()
		collaborationRedirect(w, r, "notice", "assigned")
	})

	mux.HandleFunc("POST /collaboration/watch", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulaire invalide", http.StatusBadRequest)
			return
		}
		state.mu.Lock()
		if r.Form.Get("record_version") != strconv.Itoa(state.recordVersion) || r.Form.Get("record_id") != "opportunity_42" {
			state.mu.Unlock()
			collaborationRedirect(w, r, "conflict", "watch")
			return
		}
		switch r.Form.Get("intent") {
		case "watch":
			if !state.watching {
				state.watching = true
				state.watcherCount++
			}
		case "unwatch":
			if state.watching {
				state.watching = false
				state.watcherCount = max(0, state.watcherCount-1)
			}
		default:
			state.mu.Unlock()
			collaborationRedirect(w, r, "error", "watch")
			return
		}
		state.recordVersion++
		state.mu.Unlock()
		collaborationRedirect(w, r, "notice", "watched")
	})

	mux.HandleFunc("POST /collaboration/comments", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulaire invalide", http.StatusBadRequest)
			return
		}
		state.mu.Lock()
		if r.Form.Get("thread_version") != strconv.Itoa(state.threadVersion) {
			state.mu.Unlock()
			collaborationRedirect(w, r, "conflict", "comment")
			return
		}
		body := strings.TrimSpace(r.Form.Get("body"))
		if body == "" {
			state.mu.Unlock()
			collaborationRedirect(w, r, "error", "comment")
			return
		}
		state.comments = append(state.comments, ui.Comment{
			ID: "comment_" + strconv.Itoa(75+len(state.comments)), Author: "Amélie Rousseau", Initials: "AR",
			Body: body, At: "à l’instant", Datetime: "2026-08-25T15:20:00-04:00",
		})
		state.threadVersion++
		state.mu.Unlock()
		collaborationRedirect(w, r, "notice", "commented")
	})

	mux.HandleFunc("POST /collaboration/tasks", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulaire invalide", http.StatusBadRequest)
			return
		}
		state.mu.Lock()
		if r.Form.Get("task_version") != strconv.Itoa(state.taskVersion) {
			state.mu.Unlock()
			collaborationRedirect(w, r, "conflict", "task")
			return
		}
		intent, id, ok := strings.Cut(r.Form.Get("intent"), ":")
		changed := false
		for i := range state.tasks {
			if state.tasks[i].ID != id || state.tasks[i].ReadOnly {
				continue
			}
			switch intent {
			case "complete":
				state.tasks[i].Completed = true
				changed = true
			case "reopen":
				state.tasks[i].Completed = false
				changed = true
			}
			break
		}
		if !ok || !changed {
			state.mu.Unlock()
			collaborationRedirect(w, r, "error", "task")
			return
		}
		state.taskVersion++
		state.mu.Unlock()
		collaborationRedirect(w, r, "notice", "task")
	})

	return state
}
