package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestTier9MutationsAreVersionedPOSTs(t *testing.T) {
	const mutation = "/MUTATION"
	cases := map[string]struct {
		component templ.Component
		version   string
	}{
		"assignment": {
			component: AssigneePicker(AssigneePickerProps{ID: "assignee", Action: mutation, Version: "12", Options: []Option{{Value: "u1", Label: "Nora"}}}),
			version:   `name="record_version" value="12"`,
		},
		"comment": {
			component: CommentThread(CommentThreadProps{ID: "comments", Action: mutation, Version: "8"}),
			version:   `name="thread_version" value="8"`,
		},
		"watch": {
			component: WatchButton(WatchButtonProps{ID: "watch", Action: mutation, RecordID: "r1", Version: "4"}),
			version:   `name="record_version" value="4"`,
		},
	}

	for name, tc := range cases {
		html := renderTier4(t, tc.component)
		if !strings.Contains(html, `method="post"`) || !strings.Contains(html, tc.version) {
			t.Errorf("%s must be a versioned POST: %s", name, html)
		}
		for _, href := range anchorHrefs(html) {
			if href == mutation {
				t.Errorf("%s mutation leaked into a GET link: %s", name, html)
			}
		}
	}
}

func TestTier9FormControlsHaveFieldsets(t *testing.T) {
	assignment := renderTier4(t, AssigneePicker(AssigneePickerProps{ID: "assignee", Action: "/assign"}))
	if !strings.Contains(assignment, `<fieldset`) || !strings.Contains(assignment, `name="assignee_id"`) {
		t.Fatalf("assignment control needs its validation fieldset: %s", assignment)
	}

	comments := renderTier4(t, CommentThread(CommentThreadProps{ID: "comments", Action: "/comments"}))
	if !strings.Contains(comments, `<fieldset`) || !strings.Contains(comments, `name="body"`) {
		t.Fatalf("comment body needs its validation fieldset: %s", comments)
	}
}

func TestCommentThreadEscapesOperatorTextAndTargetsItself(t *testing.T) {
	html := renderTier4(t, CommentThread(CommentThreadProps{
		ID: "comments", Action: "/comments", Comments: []Comment{{ID: "c1", Author: "Nora", Body: `<script>alert("x")</script>`}},
	}))
	if strings.Contains(html, `<script>`) || !strings.Contains(html, `&lt;script&gt;`) {
		t.Fatalf("comment bodies must remain escaped plain text: %s", html)
	}
	if !strings.Contains(html, `up-target="#comments"`) {
		t.Fatalf("a successful comment must replace its stable thread fragment: %s", html)
	}
}

func TestWatchButtonSubmitsExplicitIntent(t *testing.T) {
	watching := renderTier4(t, WatchButton(WatchButtonProps{ID: "watch", Action: "/watch", Watching: true, Count: 3}))
	if !strings.Contains(watching, `value="unwatch"`) || !strings.Contains(watching, `aria-pressed="true"`) {
		t.Fatalf("watching state needs an explicit inverse intent: %s", watching)
	}

	notWatching := renderTier4(t, WatchButton(WatchButtonProps{ID: "watch", Action: "/watch"}))
	if !strings.Contains(notWatching, `value="watch"`) || !strings.Contains(notWatching, `aria-pressed="false"`) {
		t.Fatalf("unwatched state needs an explicit watch intent: %s", notWatching)
	}
}

func TestActivityDetailNavigationPreservesHistory(t *testing.T) {
	html := renderTier4(t, ActivityFeed(ActivityFeedProps{
		ID: "activity", DetailTarget: "#detail", Items: []ActivityItem{{ID: "a1", Actor: "Nora", Action: "a commenté", Href: "/activity/a1"}},
	}))
	for _, want := range []string{`href="/activity/a1"`, `up-target="#detail"`, `up-history="true"`} {
		if !strings.Contains(html, want) {
			t.Errorf("targeted activity navigation missing %q: %s", want, html)
		}
	}
}

func TestDeadlineStateUsesExplicitNow(t *testing.T) {
	now := time.Date(2026, 8, 25, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		p    DeadlineStatusProps
		want string
	}{
		{"scheduled", DeadlineStatusProps{At: now.Add(48 * time.Hour), Now: now}, "Planifié"},
		{"due soon", DeadlineStatusProps{At: now.Add(2 * time.Hour), Now: now, WarningBefore: 4 * time.Hour}, "Bientôt"},
		{"overdue", DeadlineStatusProps{At: now.Add(-time.Minute), Now: now}, "En retard"},
		{"completed", DeadlineStatusProps{At: now.Add(-time.Hour), Now: now, Completed: true}, "Terminé"},
	}
	for _, tc := range cases {
		html := renderTier4(t, DeadlineStatus(tc.p))
		if !strings.Contains(html, tc.want) {
			t.Errorf("%s: expected %q: %s", tc.name, tc.want, html)
		}
	}

	if html := renderTier4(t, DeadlineStatus(DeadlineStatusProps{Now: now})); html != "" {
		t.Fatalf("a missing deadline should render no misleading status, got %q", html)
	}
}

func TestRecordPresenceIsAdvisoryAndPollable(t *testing.T) {
	html := renderTier4(t, RecordPresence(RecordPresenceProps{
		ID: "presence", Busy: true, PollSource: "/records/r1/presence",
		People: []PresencePerson{{ID: "u1", Name: "Nora", Initials: "NM", Mode: PresenceEditing}},
	}))
	for _, want := range []string{`up-poll`, `up-interval="5000"`, `up-source="/records/r1/presence"`, `data-user-id="u1"`, "modifie cette fiche"} {
		if !strings.Contains(html, want) {
			t.Errorf("presence surface missing %q: %s", want, html)
		}
	}
}

func TestTaskListUsesExactVersionedPOSTIntents(t *testing.T) {
	html := renderTier4(t, TaskList(TaskListProps{
		ID: "tasks", Action: "/tasks", Version: "9", Items: []TaskItem{
			{ID: "t1", Title: "Appeler le client"},
			{ID: "t2", Title: "Envoyer le devis", Completed: true},
			{ID: "t3", Title: "Privée", ReadOnly: true},
		},
	}))
	for _, want := range []string{`method="post"`, `name="task_version" value="9"`, `value="complete:t1"`, `value="reopen:t2"`, `aria-pressed="true"`} {
		if !strings.Contains(html, want) {
			t.Errorf("task list missing %q: %s", want, html)
		}
	}
	if !strings.Contains(html, `disabled`) || !strings.Contains(html, "Lecture seule") {
		t.Fatalf("read-only tasks must explain why they cannot change: %s", html)
	}
	for _, href := range anchorHrefs(html) {
		if href == "/tasks" {
			t.Fatalf("task mutation leaked into a GET link: %s", html)
		}
	}
}

func TestEmptyTaskListHasNoMutationControls(t *testing.T) {
	html := renderTier4(t, TaskList(TaskListProps{ID: "tasks", Action: "/tasks"}))
	if !strings.Contains(html, "Aucune tâche") || strings.Contains(html, `<form`) {
		t.Fatalf("empty task list should explain itself without a meaningless form: %s", html)
	}
}
