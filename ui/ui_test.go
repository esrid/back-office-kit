package ui

import (
	"net/url"
	"reflect"
	"testing"
)

func TestPages(t *testing.T) {
	cases := []struct {
		current, total int
		want           []int
	}{
		{1, 0, nil},
		{1, 1, nil},
		{1, 3, []int{1, 2, 3}},
		{1, 10, []int{1, 2, 0, 10}},
		{5, 10, []int{1, 0, 4, 5, 6, 0, 10}},
		{10, 10, []int{1, 0, 9, 10}},
		{2, 10, []int{1, 2, 3, 0, 10}},
		{3, 10, []int{1, 2, 3, 4, 0, 10}},
		{0, 5, []int{1, 2, 0, 5}},  // clamped low
		{99, 5, []int{1, 0, 4, 5}}, // clamped high
	}
	for _, c := range cases {
		if got := Pages(c.current, c.total); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Pages(%d,%d) = %v, want %v", c.current, c.total, got, c.want)
		}
	}
}

func TestPagesNeverRepeatsOrGapsAdjacent(t *testing.T) {
	for total := 2; total <= 40; total++ {
		for cur := 1; cur <= total; cur++ {
			got := Pages(cur, total)
			seen := map[int]bool{}
			prev := -1
			for i, p := range got {
				if p == gap {
					if i == 0 || i == len(got)-1 || got[i-1] == gap {
						t.Fatalf("Pages(%d,%d)=%v: misplaced gap at %d", cur, total, got, i)
					}
					prev = -1
					continue
				}
				if seen[p] {
					t.Fatalf("Pages(%d,%d)=%v: duplicate %d", cur, total, got, p)
				}
				if prev != -1 && p != prev+1 {
					t.Fatalf("Pages(%d,%d)=%v: gap not marked between %d and %d", cur, total, got, prev, p)
				}
				seen[p] = true
				prev = p
			}
			if !seen[1] || !seen[total] {
				t.Fatalf("Pages(%d,%d)=%v: missing first or last", cur, total, got)
			}
			if !seen[cur] {
				t.Fatalf("Pages(%d,%d)=%v: missing current", cur, total, got)
			}
		}
	}
}

func TestTotalPages(t *testing.T) {
	cases := []struct{ n, per, want int }{
		{0, 20, 0}, {1, 20, 1}, {20, 20, 1}, {21, 20, 2}, {40, 20, 2}, {41, 20, 3}, {10, 0, 0},
	}
	for _, c := range cases {
		if got := TotalPages(c.n, c.per); got != c.want {
			t.Errorf("TotalPages(%d,%d) = %d, want %d", c.n, c.per, got, c.want)
		}
	}
}

func TestPageHrefPreservesFiltersAndDropsPageOne(t *testing.T) {
	q := url.Values{"q": {"bob"}, "status": {"active"}, "page": {"7"}}
	if got, want := PageHref(q, 3), "?page=3&q=bob&status=active"; got != want {
		t.Errorf("PageHref = %q, want %q", got, want)
	}
	if got, want := PageHref(q, 1), "?q=bob&status=active"; got != want {
		t.Errorf("PageHref page 1 = %q, want %q", got, want)
	}
	if q.Get("page") != "7" {
		t.Error("PageHref mutated the caller's url.Values")
	}
}

func TestSortHrefToggles(t *testing.T) {
	base := url.Values{"q": {"bob"}, "page": {"4"}}

	first := SortHref(base, "email")
	if want := "?dir=asc&q=bob&sort=email"; first != want {
		t.Errorf("first sort = %q, want %q", first, want)
	}
	if base.Get("page") != "4" {
		t.Error("SortHref mutated the caller's url.Values")
	}

	asc := url.Values{"sort": {"email"}, "dir": {"asc"}}
	if got, want := SortHref(asc, "email"), "?dir=desc&sort=email"; got != want {
		t.Errorf("toggle = %q, want %q", got, want)
	}

	desc := url.Values{"sort": {"email"}, "dir": {"desc"}}
	if got, want := SortHref(desc, "email"), "?dir=asc&sort=email"; got != want {
		t.Errorf("toggle back = %q, want %q", got, want)
	}

	// switching column always starts ascending
	if got, want := SortHref(desc, "name"), "?dir=asc&sort=name"; got != want {
		t.Errorf("switch column = %q, want %q", got, want)
	}
}

func TestSortDir(t *testing.T) {
	q := url.Values{"sort": {"email"}, "dir": {"desc"}}
	if got := SortDir(q, "email"); got != "desc" {
		t.Errorf("SortDir active = %q", got)
	}
	if got := SortDir(q, "name"); got != "" {
		t.Errorf("SortDir inactive = %q", got)
	}
	if got := SortDir(q, ""); got != "" {
		t.Errorf("SortDir unsortable = %q", got)
	}
	if got := SortDir(url.Values{"sort": {"email"}}, "email"); got != "asc" {
		t.Errorf("SortDir default = %q", got)
	}
}

func TestToneFor(t *testing.T) {
	cases := map[string]string{
		"active": "success", "  Active ": "success", "PAID": "success",
		"pending": "warning", "failed": "error", "draft": "neutral",
		"wat": "neutral", "": "neutral",
	}
	for in, want := range cases {
		if got := ToneFor(in); got != want {
			t.Errorf("ToneFor(%q) = %q, want %q", in, got, want)
		}
	}
}
