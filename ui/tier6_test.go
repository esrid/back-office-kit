package ui

import (
	"strings"
	"testing"
)

func TestCommandPaletteKeepsWritesOutOfLinks(t *testing.T) {
	const destructive = "/accounts/42/suspend"
	html := renderTier4(t, CommandPalette(CommandPaletteProps{
		ID: "commands", SearchAction: "/commands",
		Groups: []PaletteGroup{{Label: "Comptes", Commands: []PaletteCommand{
			{ID: "view", Label: "Voir", Href: "/accounts/42"},
			{ID: "suspend", Label: "Suspendre", Href: destructive, Method: "post", Confirm: "Confirmer ?"},
		}}},
	}))

	for _, href := range anchorHrefs(html) {
		if href == destructive {
			t.Fatalf("state-changing command must not be a link: %s", html)
		}
	}
	for _, want := range []string{`form="commands-actions"`, `formaction="/accounts/42/suspend"`, `up-method="post"`, `<a href="/accounts/42"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("command palette missing %q: %s", want, html)
		}
	}
}

func TestQueryBuilderSerializesNestedGroupsInOneForm(t *testing.T) {
	html := renderTier4(t, QueryBuilder(QueryBuilderProps{
		ID: "query", Action: "/accounts", Fields: []Option{{Value: "status", Label: "Statut"}}, Operators: []Option{{Value: "eq", Label: "est"}},
		Root: QueryGroup{ID: "root", Join: QueryAll, AllowAddGroup: true, Groups: []QueryGroup{{
			ID: "late", Join: QueryAny, AllowRemove: true,
			Rules: []QueryRule{{ID: "r1", Field: "status", Operator: "eq", Value: "late"}},
		}}},
	}))

	for _, want := range []string{
		`name="groups[root][join]"`,
		`name="groups[late][rules][r1][field]"`,
		`value="remove_group:late"`,
		`value="add_group:root"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("query builder missing %q: %s", want, html)
		}
	}
	if count := strings.Count(html, "<form"); count != 1 {
		t.Fatalf("query builder must remain one valid form, got %d: %s", count, html)
	}
}

func TestColumnManagerSubmitsLockedVisibleColumns(t *testing.T) {
	html := renderTier4(t, ColumnManager(ColumnManagerProps{
		ID: "columns", Action: "/views/columns", Density: DensityCompact,
		Columns: []ManagedColumn{{Key: "name", Label: "Nom", Visible: true, Locked: true}, {Key: "email", Label: "E-mail", Visible: true}},
	}))

	for _, want := range []string{
		`type="hidden" name="columns" value="name"`,
		`name="column_order" value="name"`,
		`name="density" value="compact" checked`,
		`value="move_down:name"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("column manager missing %q: %s", want, html)
		}
	}
}

func TestMasterDetailNavigationTargetsDetailRegion(t *testing.T) {
	nav := MasterDetailNav("Comptes", []MasterDetailItem{{ID: "42", Title: "Acme", Href: "/accounts/42"}}, "42", "#accounts-detail")
	html := renderTier4(t, MasterDetail(MasterDetailProps{
		ID: "accounts", Title: "Comptes", List: nav, DetailLabel: "Détail du compte",
	}))
	for _, want := range []string{`id="accounts-list"`, `id="accounts-detail"`, `up-target="#accounts-detail"`, `aria-current="page"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("master/detail missing %q: %s", want, html)
		}
	}
}
