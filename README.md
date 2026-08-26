# back-office-kit

Composants de back-office en [templ](https://templ.guide) + [daisyUI 5](https://daisyui.com) + [Unpoly 3](https://unpoly.com).

57 composants, des primitives (tableau triable, pagination, champs) jusqu'aux
écrans complets (liste, fiche, formulaire, réglages), aux interfaces agentiques
(plan d'action, carte d'approbation, journal d'exécution) et à la gouvernance.

## Documentation

Deux entrées, complémentaires.

**`go doc` pendant que vous codez.** C'est la référence d'API : signatures,
champs de chaque struct de props, et le pourquoi de chacun.

```sh
go doc github.com/esrid/back-office-kit/ui              # les 84 composants
go doc github.com/esrid/back-office-kit/ui FieldProps   # tous les champs, commentés
go doc github.com/esrid/back-office-kit/ui DataTable
```

**La galerie pour choisir.** Chaque composant rendu pour de vrai, sa signature
exacte extraite du source à la génération, et son code d'appel. Elle se
construit depuis le dépôt et s'ouvre en local :

```sh
npm install && templ generate
npx tailwindcss -i assets/app.css -o dist/app.css --minify
go run ./cmd/gallery && open dist/gallery.html
```

`dist/gallery.html` est autonome (CSS inline, aucune dépendance réseau hors
Google Fonts) et versionné : il s'ouvre aussi directement depuis un clone, sans
rien construire.

GitHub Pages n'est pas activé — il demande un dépôt public ou un plan payant.
Si le dépôt passe public, la source Pages devra être `/` ou `/docs` (`/dist`
n'est pas un chemin accepté), donc soit renommer le dossier de sortie, soit
publier via GitHub Actions.

## Principes

- **daisyUI donne les primitives.** Ce dépôt n'ajoute que ce qui manque à un
  back-office : tri et pagination portés par l'URL, validation serveur par
  champ, états vides, statuts cohérents, overlays d'édition.
- **Tout état vit dans l'URL.** Filtres, tri, pagination : rechargeable,
  partageable, sans état caché côté client.
- **Ça marche sans JavaScript.** Unpoly accélère, il n'est jamais requis.
- **La logique pure est testée.** Calcul des pages, construction des query
  strings, bornages : fonctions pures, jamais de mutation de l'appelant.

## L'utiliser dans votre projet

```sh
go get github.com/esrid/back-office-kit
```

Les fichiers `*_templ.go` sont versionnés, donc le paquet est utilisable sans
installer templ. Vous n'avez besoin de templ que pour vos propres templates.

Tailwind doit voir les classes du kit, sinon le CSS manquera à l'exécution.
Ajoutez la source dans votre `app.css` :

```css
@import "tailwindcss";
@plugin "daisyui";
@source "../path/vers/back-office-kit/ui";
```

Puis chargez Unpoly dans votre `<head>` — le kit fonctionne sans lui, Unpoly ne
fait qu'éviter les rechargements.

## Une page complète

`example/` est une page réelle : filtre, tri et pagination servis par un vrai
serveur Go, sans état client.

```sh
git clone https://github.com/esrid/back-office-kit && cd back-office-kit
npm install
npx tailwindcss -i assets/app.css -o dist/app.css --minify
go run ./example        # http://localhost:8080
```

L'ossature d'un écran de liste tient en un bloc :

```templ
@ui.Shell("Acme Admin", nav(), "/", operator, flashes) {
    @ui.IndexPage(ui.PageProps{Title: "Utilisateurs"}, ui.FilterBar(filters(), q, "#users")) {
        @ui.DataTable("users", columns(), rows, q, nil) {
            @ui.Pagination(page, total, 20, q, "#users")
        }
    }
}
```

Le handler lit `r.URL.Query()`, filtre, trie, pagine, et rend. Le kit ne garde
aucun état : `ui.SortHref`, `ui.PageHref` et `ui.ResetHref` construisent les
URL, votre code les applique aux données.

## Construire

```sh
make generate   # templ generate + détache le commentaire de version des générés
make css        # tailwind + daisyui -> dist/app.css
make gallery    # écrit dist/gallery.html
make dev        # lance example/ sur :8080
make test
```

`make generate` fait plus que `templ generate` : templ colle
`// templ: version: X` au `package ui` de chaque fichier généré, et Go prend
alors ces 63 commentaires pour la documentation du paquet — `go doc ./ui`
devient illisible. La cible insère une ligne vide qui les détache.

## Côté serveur

### Protection CSRF

Aucun composant n'émet de jeton, et aucun n'en a besoin :

```go
p := http.NewCrossOriginProtection()   // Go 1.25+
srv := &http.Server{Handler: p.Handler(mux)}
```

Elle rejette les requêtes cross-origin non-sûres via `Sec-Fetch-Site`, sans
jeton ni cookie. Corollaire assumé dans tout le kit : **aucun changement d'état
ne passe par un GET**, puisque `CrossOriginProtection` laisse toujours passer
GET, HEAD et OPTIONS. Les actions destructives sont des boutons `submit`
(`[formaction]`), jamais des liens. `ui/nostategetlinks_test.go` garde cette
propriété.

Sans JavaScript, `RowActionMenu` et `FormActions` envoient un POST vers l'URL
de l'action ; le verbe réel ne voyage que dans `[up-method]` pour Unpoly.

### Validation de formulaire

Chaque champ est son propre `<fieldset>`, ce que `[up-validate]` cible par
défaut (`X-Up-Target: fieldset:has(#f-email)`). Le serveur valide sans
committer et re-rend le formulaire ; seul le champ modifié bouge.

### Messages flash

Émettez toujours la zone `[up-flashes]` : Unpoly la met à jour même quand elle
n'est pas ciblée, et un conteneur vide n'efface pas les messages existants.

## Deux pièges qui ne se voient qu'à l'exécution

**L'historique du navigateur.** Unpoly ne change l'URL que lorsqu'une cible
principale est rendue — *« This is to prevent location changes when rendering a
minor fragment »*. Un tableau ciblé par `#users` n'en est pas une, donc tri et
pagination s'appliquaient sans apparaître dans l'URL et disparaissaient au
rechargement. `DataTable`, `Pagination`, `FilterBar` et `DetailTabs` portent
donc `[up-history="true"]`. Les fragments réellement mineurs — `InlineEdit`,
`RepeaterField`, `AsyncState` — gardent le défaut : ils ne doivent pas polluer
l'historique.

## Une règle Tailwind à ne pas oublier

Les classes de couleur sont écrites **en toutes lettres** dans `ui/tone.go`,
jamais concaténées. `"badge-" + tone` compile, mais le scanner Tailwind ne le
voit pas et le CSS manque à l'exécution.

## Deux paquets, pas un

`ui` est le back-office : dense, authentifié, piloté par la donnée. `ui/marketing`
est la page publique : référencement, images, conversion. Ils ne partagent ni les
contraintes ni les règles, d'où la séparation.

```go
import "github.com/esrid/back-office-kit/ui"
import "github.com/esrid/back-office-kit/ui/marketing"
```

Trois règles s'appliquent dans `marketing` et pas dans `ui` :

- **Un seul `<h1>` par page**, rendu par `Hero`. Les autres sections ouvrent en
  `<h2>`. Une hiérarchie de titres fausse ment aux moteurs et aux lecteurs
  d'écran sans que rien ne se voie à l'œil.
- **Toute image porte `alt`, `width` et `height`.** Sans `alt` elle ne rend pas
  du tout, plutôt que de partir cassée ; sans dimensions la page saute au
  chargement.
- **Aucun JavaScript.** La FAQ est un `<details>`, ce qui laisse la recherche du
  navigateur atteindre les réponses fermées.

## Licence

MIT
