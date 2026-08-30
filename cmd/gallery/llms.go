package main

import (
	"bufio"
	"io"
	"os"
)

// llmsHeader is everything an agent must know before writing a single call:
// the import paths, the build commands, and the invariants that fail silently.
// Kept in sync by hand with AGENTS.md; the component list below is generated.
const llmsHeader = `# back-office-kit

> Composants server-rendered pour back-office Go : templ pour le rendu, daisyUI 5
> pour les classes, Unpoly 3 pour éviter les rechargements. Deux paquets :
> ` + "`ui`" + ` (back-office) et ` + "`ui/marketing`" + ` (pages publiques). Le kit ne stocke
> aucun état : tri, filtres et pagination vivent dans la query string.

## Installation

    go get github.com/esrid/back-office-kit

Les ` + "`*_templ.go`" + ` sont versionnés : importer et rendre ne demande ni templ ni npm.
Tailwind doit voir les sources du kit, sinon le CSS manque à l'exécution — le
chemin se calcule, il ne s'écrit pas :

    KIT=$(go list -m -f '{{.Dir}}' github.com/esrid/back-office-kit)
    { echo '@source "./**/*.templ";'; echo "@source \"$KIT/ui\";"; } > sources.generated.css

Puis dans votre entrée Tailwind : ` + "`@import \"tailwindcss\"; @plugin \"daisyui\"; @import \"./sources.generated.css\";`" + `

## Invariants — les casser compile et ne se voit qu'à l'exécution

- Aucun changement d'état ne passe par un GET. http.CrossOriginProtection laisse
  toujours passer GET/HEAD/OPTIONS, et [up-preload] déclenche un lien sans clic.
  Les actions destructives sont des <button type="submit"> avec [formaction].
- [up-layer] n'accepte que new, swap, shatter. Le type d'overlay va dans [up-mode].
- Les classes de couleur s'écrivent en toutes lettres. "badge-" + tone compile,
  mais le scanner Tailwind ne le voit pas.
- Un fragment porteur d'état a besoin de [up-history="true"] : Unpoly ne change
  l'URL que pour une cible principale. Concerne DataTable, Pagination, FilterBar,
  DetailTabs. Exception documentée : LoadMore, où aucune URL ne peut décrire des
  pages ajoutées.
- Chaque champ de formulaire est son propre <fieldset> : c'est ce que
  [up-validate] cible (X-Up-Target: fieldset:has(#f-email)).
- L'argent est un entier en unités mineures, jamais un float.
- Rien ne lit l'horloge ni le hasard : Timestamp prend Now en paramètre.
- Dans ui/marketing : un seul <h1> (rendu par Hero), toute image porte alt/width/
  height, aucun JavaScript.

## Composants

Chaque entrée donne la signature exacte, extraite du source à la génération.
` + "`go doc github.com/esrid/back-office-kit/ui <Nom>`" + ` donne les champs de props.
`

// writeLLMs renders llms.txt: the same sections as the gallery, in the format
// an agent reads instead of a 470 KB page of rendered HTML.
func writeLLMs(path string, secs []Section) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if _, err := io.WriteString(w, llmsHeader); err != nil {
		return err
	}

	group := ""
	for _, s := range secs {
		if s.Group != group {
			group = s.Group
			if _, err := io.WriteString(w, "\n### "+group+"\n\n"); err != nil {
				return err
			}
		}
		line := "- **" + s.Title + "** — " + s.Purpose + "\n"
		if s.Signature != "" {
			line += "  `" + s.Signature + "`\n"
		}
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return w.Flush()
}
