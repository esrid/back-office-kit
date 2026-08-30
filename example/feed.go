package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/esrid/back-office-kit/ui"
)

// feedPage is small on purpose: a cursor list is read in one screen, and the
// button is the only way forward.
const feedPage = 12

// Event is one line of an append-only log -- the shape LoadMore exists for:
// counting it would cost a full scan, and a new row shifts every offset.
type Event struct {
	ID    string
	Actor string
	Verb  string
	At    time.Time
}

type FeedView struct {
	Operator string
	Events   []Event
	NextHref string // empty on the last page
	Shown    int
}

// afterCursor returns the page starting after the event with the given id, and
// the cursor of the page after it. An unknown cursor starts from the top rather
// than erroring: a stale link should show something. Pure.
func afterCursor(all []Event, cursor string, size int) (page []Event, next string) {
	start := 0
	if cursor != "" {
		for i, e := range all {
			if e.ID == cursor {
				start = i + 1
				break
			}
		}
	}
	end := min(start+size, len(all))
	page = all[start:end]
	if end < len(all) && len(page) > 0 {
		next = page[len(page)-1].ID
	}
	return page, next
}

func registerFeed(mux *http.ServeMux, operator string) {
	events := seedEvents()

	mux.HandleFunc("GET /feed", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, next := afterCursor(events, q.Get("cursor"), feedPage)

		view := FeedView{Operator: operator, Events: page, Shown: len(page)}
		if next != "" {
			view.NextHref = "/feed" + ui.CursorHref(url.Values{}, next)
		}
		render(r.Context(), w, FeedPage(view))
	})
}

func seedEvents() []Event {
	actors := []string{"amelie@acme.co", "karim@acme.co", "sofia@acme.co", "lucas@acme.co"}
	verbs := []string{
		"a suspendu un compte", "a validé une facture", "a exporté 1 240 lignes",
		"a changé un rôle", "a révoqué une clé d'API", "a rejoué un webhook",
	}
	base := time.Date(2026, 3, 14, 9, 0, 0, 0, time.UTC)

	out := make([]Event, 48)
	for i := range out {
		out[i] = Event{
			ID:    fmt.Sprintf("e_%03d", i),
			Actor: actors[i%len(actors)],
			Verb:  verbs[i%len(verbs)],
			At:    base.Add(-time.Duration(i) * 17 * time.Minute),
		}
	}
	return out
}
