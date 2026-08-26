package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func postCollaboration(t *testing.T, mux http.Handler, path string, form url.Values) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func TestCollaborationScreenRendersIntegratedComponents(t *testing.T) {
	mux := http.NewServeMux()
	registerCollaboration(mux, "amelie@acme.co")
	req := httptest.NewRequest(http.MethodGet, "/collaboration", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("GET /collaboration = %d", rr.Code)
	}
	html := rr.Body.String()
	for _, want := range []string{"collaboration-assignee", "collaboration-comments", "collaboration-tasks", "collaboration-presence", "collaboration-activity"} {
		if !strings.Contains(html, want) {
			t.Errorf("integrated screen missing %q", want)
		}
	}
	if got := strings.Count(html, "<h1"); got != 1 {
		t.Errorf("screen has %d h1 elements, want exactly one", got)
	}
}

func TestCollaborationAssignmentIsVersioned(t *testing.T) {
	mux := http.NewServeMux()
	state := registerCollaboration(mux, "amelie@acme.co")

	rr := postCollaboration(t, mux, "/collaboration/assignee", url.Values{
		"record_version": {"18"},
		"assignee_id":    {"usr_nora"},
	})
	if rr.Code != http.StatusSeeOther || rr.Header().Get("Location") != "/collaboration?notice=assigned" {
		t.Fatalf("assignment response = %d %q", rr.Code, rr.Header().Get("Location"))
	}
	v := state.view("operator", nil)
	if v.AssigneeID != "usr_nora" || v.RecordVersion != "19" {
		t.Fatalf("assignment state = assignee %q version %q", v.AssigneeID, v.RecordVersion)
	}

	stale := postCollaboration(t, mux, "/collaboration/assignee", url.Values{
		"record_version": {"18"},
		"assignee_id":    {"usr_bruno"},
	})
	if stale.Header().Get("Location") != "/collaboration?conflict=assignment" {
		t.Fatalf("stale assignment should surface a conflict, got %q", stale.Header().Get("Location"))
	}
	if got := state.view("operator", nil).AssigneeID; got != "usr_nora" {
		t.Fatalf("stale assignment mutated state to %q", got)
	}
}

func TestCollaborationCommentAndTaskMutations(t *testing.T) {
	mux := http.NewServeMux()
	state := registerCollaboration(mux, "amelie@acme.co")

	comment := postCollaboration(t, mux, "/collaboration/comments", url.Values{
		"thread_version": {"7"},
		"body":           {`<script>alert("x")</script>`},
	})
	if comment.Code != http.StatusSeeOther {
		t.Fatalf("comment response = %d", comment.Code)
	}
	v := state.view("operator", nil)
	if v.ThreadVersion != "8" || len(v.Comments) != 3 {
		t.Fatalf("comment state = version %q, %d comments", v.ThreadVersion, len(v.Comments))
	}

	task := postCollaboration(t, mux, "/collaboration/tasks", url.Values{
		"task_version": {"11"},
		"intent":       {"complete:task_81"},
	})
	if task.Code != http.StatusSeeOther {
		t.Fatalf("task response = %d", task.Code)
	}
	v = state.view("operator", nil)
	if v.TaskVersion != "12" || !v.Tasks[0].Completed {
		t.Fatalf("task state = version %q completed %v", v.TaskVersion, v.Tasks[0].Completed)
	}
}

func TestCollaborationMutationsRejectGET(t *testing.T) {
	mux := http.NewServeMux()
	state := registerCollaboration(mux, "amelie@acme.co")
	req := httptest.NewRequest(http.MethodGet, "/collaboration/tasks?intent=complete:task_81&task_version=11", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET mutation route = %d, want 405", rr.Code)
	}
	if state.view("operator", nil).Tasks[0].Completed {
		t.Fatal("GET completed a task")
	}
}
