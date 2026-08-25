package main

import (
	"net/url"
	"testing"
)

func firstNames(rows []User) []string {
	out := make([]string, len(rows))
	for i, u := range rows {
		out[i] = u.Name
	}
	return out
}

func TestSelectUsers(t *testing.T) {
	all := seed()

	t.Run("le total compte les lignes filtrées, pas la page", func(t *testing.T) {
		rows, total := selectUsers(all, url.Values{})
		if total != len(all) {
			t.Errorf("total = %d, want %d", total, len(all))
		}
		if len(rows) != perPage {
			t.Errorf("page = %d lignes, want %d", len(rows), perPage)
		}
	})

	t.Run("le filtre réduit aussi le total", func(t *testing.T) {
		_, total := selectUsers(all, url.Values{"status": {"failed"}})
		if total != 2 {
			t.Errorf("total filtré = %d, want 2", total)
		}
	})

	t.Run("la recherche regarde le nom et l'e-mail", func(t *testing.T) {
		_, byName := selectUsers(all, url.Values{"q": {"nguyen"}})
		_, byMail := selectUsers(all, url.Values{"q": {"partner.io"}})
		if byName != 1 || byMail != 2 {
			t.Errorf("nom=%d mail=%d, want 1 et 2", byName, byMail)
		}
	})

	t.Run("le tri descendant inverse l'ascendant", func(t *testing.T) {
		asc, _ := selectUsers(all, url.Values{"sort": {"seats"}, "dir": {"asc"}})
		desc, _ := selectUsers(all, url.Values{"sort": {"seats"}, "dir": {"desc"}})
		if asc[0].Seats > asc[1].Seats {
			t.Errorf("asc mal trié: %v", firstNames(asc))
		}
		if desc[0].Seats < desc[1].Seats {
			t.Errorf("desc mal trié: %v", firstNames(desc))
		}
		if asc[0].Name == desc[0].Name {
			t.Errorf("les deux sens donnent la même tête: %s", asc[0].Name)
		}
	})

	t.Run("une page au-delà de la fin ne panique pas", func(t *testing.T) {
		rows, total := selectUsers(all, url.Values{"page": {"999"}})
		if len(rows) != 0 || total != len(all) {
			t.Errorf("page hors bornes: %d lignes, total %d", len(rows), total)
		}
	})

	t.Run("un filtre sans résultat rend une page vide, pas une erreur", func(t *testing.T) {
		rows, total := selectUsers(all, url.Values{"q": {"zzzz"}})
		if len(rows) != 0 || total != 0 {
			t.Errorf("got %d lignes / total %d", len(rows), total)
		}
	})
}
