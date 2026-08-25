# back-office-kit

Composants de back-office en [templ](https://templ.guide) + [daisyUI 5](https://daisyui.com) + [Unpoly 3](https://unpoly.com).

57 composants, des primitives (tableau triable, pagination, champs) jusqu'aux
écrans complets (liste, fiche, formulaire, réglages), aux interfaces agentiques
(plan d'action, carte d'approbation, journal d'exécution) et à la gouvernance.

## Documentation

La galerie rend chaque composant pour de vrai, avec son code d'appel. Elle se
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

## Démarrer

```sh
npm install
templ generate                 # les *_templ.go ne sont pas versionnés
npx tailwindcss -i assets/app.css -o dist/app.css --minify
go run ./cmd/gallery           # écrit dist/gallery.html
go test ./ui/
```

## Utiliser les composants

```go
cols := []ui.Column[User]{
    ui.Text("Nom", "name", func(u User) string { return u.Name }),
    {Header: "Statut", Sort: "status",
     Cell: func(u User) templ.Component { return ui.StatusBadge(u.Status) }},
}
```

```templ
@ui.IndexPage(ui.PageProps{Title: "Utilisateurs", PrimaryAction: inviteButton()},
    ui.FilterBar(filters, r.URL.Query(), "#users")) {
    @ui.DataTable("users", cols, users, r.URL.Query(), nil) {
        @ui.Pagination(page, total, 20, r.URL.Query(), "#users")
    }
}
```

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

## Une règle Tailwind à ne pas oublier

Les classes de couleur sont écrites **en toutes lettres** dans `ui/tone.go`,
jamais concaténées. `"badge-" + tone` compile, mais le scanner Tailwind ne le
voit pas et le CSS manque à l'exécution.

## Licence

MIT
