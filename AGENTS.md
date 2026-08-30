# Travailler sur ce dépôt

Bibliothèque de composants server-rendered : **templ** pour le rendu, **daisyUI 5**
pour les classes, **Unpoly 3** pour éviter les rechargements. Deux paquets :
`ui` (back-office) et `ui/marketing` (pages publiques).

Ce fichier existe parce que la plupart des bugs trouvés ici étaient **silencieux** :
ils compilaient, s'affichaient correctement, et ne faisaient rien. Les règles
ci-dessous viennent chacune d'un bug réel.

## Commandes

```sh
make generate   # templ generate + détache le commentaire de version des générés
make css        # tailwind + daisyui -> dist/app.css
make gallery    # écrit dist/gallery.html et llms.txt
make dev        # lance example/ sur :8080
make test
```

`make generate`, jamais `templ generate` seul : templ colle
`// templ: version: X` au `package ui` de chaque fichier généré, et Go prend
alors ces 63 commentaires pour la documentation du paquet — `go doc ./ui`
devient illisible.

## Invariants — ne pas les casser

**Aucun changement d'état ne passe par un GET.** `http.CrossOriginProtection`
laisse toujours passer GET, HEAD et OPTIONS, et un préchargement ou `[up-preload]`
peut déclencher un lien sans clic. Les actions destructives sont des `<button
type="submit">` avec `[formaction]`, jamais des `<a>`. Gardé par
`ui/nostategetlinks_test.go`.

**`[up-layer]` n'accepte que `new`, `swap`, `shatter`.** Le type d'overlay va
dans `[up-mode]`. Écrire `up-layer="new drawer"` compile, s'affiche, et le tiroir
ne s'ouvre jamais comme tiroir. Gardé par `TestNoComponentFoldsTheLayerMode`.

**Les classes de couleur sont écrites en toutes lettres.** `"badge-" + tone`
compile mais le scanner Tailwind ne le voit pas et le CSS manque à l'exécution.
Table littérale dans `ui/tone.go`.

**Un fragment qui porte de l'état a besoin de `[up-history="true"]`.** C'est le
défaut le plus fréquent du dépôt : cinq composants l'avaient, tous trouvés en
exécutant, aucun en relisant. Unpoly ne
change l'URL que lorsqu'une cible principale est rendue — « This is to prevent
location changes when rendering a minor fragment ». Un tableau ciblé par `#users`
n'en est pas une : sans cet attribut le tri s'applique sans apparaître dans l'URL
et disparaît au rechargement. Concerne `DataTable`, `Pagination`, `FilterBar`,
`DetailTabs`. Ne concerne pas `InlineEdit`, `RepeaterField`, `AsyncState` — ce
sont de vrais fragments mineurs.

**Chaque champ de formulaire est son propre `<fieldset>`.** C'est ce que
`[up-validate]` cible par défaut : `X-Up-Target: fieldset:has(#f-email)`.

**Rien ne lit l'horloge ni le hasard.** `Timestamp` prend `Now` en paramètre. Un
composant qui appelle `time.Now()` rend différemment à chaque appel et devient
intestable.

**L'argent est un entier en unités mineures.** Jamais un float. Le nombre de
décimales suit la devise — supposer deux transforme `1000 ¥` en `10,00 ¥`.

**Les helpers d'URL ne mutent jamais l'`url.Values` reçue.** `SortHref`,
`PageHref`, `ResetHref` clonent. Testé.

### Règles propres à `ui/marketing`

Un seul `<h1>` par page, rendu par `Hero` ; les autres sections ouvrent en `<h2>`.
Toute image porte `alt`, `width` et `height` — sans `alt` elle ne rend pas du
tout. Aucun JavaScript : la FAQ est un `<details>`.

## Pièges de templ

**`{` démarre une expression dans le contenu.** `@Money(MoneyProps{...})` écrit
au milieu d'un élément ne parse pas. Passer par une variable :

```templ
{{ montant := MoneyProps{Cents: it.TotalCents, Currency: p.Currency} }}
@Money(montant)
```

**Pas de commentaire `//` dans un bloc d'attributs.** templ ne les parse pas :

```templ
<a href={ x }
   up-target={ t }
   // ceci casse la génération
   up-history="true">
```

L'explication va dans le commentaire de doc du composant, au-dessus du `templ`.

**Les imports sont par fichier.** Un `.templ` a son propre bloc `import`, un
`.go` aussi, même dans le même paquet. templ accepte du Go de premier niveau,
donc un composant et sa logique peuvent tenir dans un seul `.templ`.

**`gofmt` ne voit pas les `.templ`.** C'est `templ fmt .` qui les formate.

**Les `*_templ.go` sont versionnés.** Sans eux, `go get` livre un paquet qui
compile et n'expose aucun composant — un échec silencieux. Les commiter après
chaque `make generate`.

## Vérifier son travail

Le compilateur ne dit presque rien sur ce qui compte ici. Ces quatre contrôles
ont trouvé chacun au moins un bug réel.

**1. Le rendu, pas la relecture.** Construire, servir, regarder.

```sh
make gallery && (cd dist && python3 -m http.server 8731)
```

**2. Les classes daisyUI existent-elles dans le CSS construit ?**

```sh
grep -c '\.badge-success' dist/app.css     # 0 = la classe manquera à l'exécution
```

**3. Les attributs Unpoly existent-ils ?** Lire la doc officielle avant d'en
écrire un nouveau. `up-layer`, `up-mode`, `up-position`, `up-poll`,
`up-interval`, `up-source`, `up-autosubmit`, `up-watch-delay`, `up-flashes`,
`up-validate`, `up-confirm`, `up-dismiss`, `up-history` sont vérifiés. Un
attribut inventé ne produit aucune erreur.

**4. Le contraste se mesure, il ne s'estime pas.** Une sonde in-page lit la
`color` avec un thème de retard juste après un changement de `data-theme`,
alors que `background-color` suit tout de suite — d'où des ratios de 1,01 sur
une page parfaitement lisible. Mesurer sur un **élément neuf** créé après le
changement, ou recharger la page dans le thème visé. Sonde canvas dans la page,
avec un contrôle de vraisemblance noir sur blanc qui doit donner 21. Trois de mes
mesures ont été fausses avant d'être justes : alpha ignoré, `oklch()` mal parsé,
puis style périmé lu juste après un changement de thème. Toujours relire jusqu'à
stabilisation, et trancher au pixel en cas de doute.

## Limites connues

- Le thème sombre de daisyUI donnait **4,13:1** sur `btn-primary`. Corrigé dans
  `assets/app.css` en assombrissant `--color-primary` à `oklch(50% …)` : 5,85:1,
  clair inchangé à 6,75. C'est l'exemple de ce que « changer le style » veut
  dire ici — deux sélecteurs, aucun composant touché.
- `example/` couvre six écrans : utilisateurs, dossiers, approbations, revue,
  assistant, politique. Le Tier 8 (sécurité) et la plupart du Tier 3 n'ont
  toujours aucun écran.
- La couverture de tests porte sur la logique pure (URL, bornages, argent,
  dates). Le balisage est vérifié par quelques assertions de rendu et à l'œil.

## Décisions closes — ne pas les rouvrir

**Les libellés sont en français, en dur, et c'est un choix.** Pas une dette à
combler. Le rattrapage resterait bon marché le jour où il servirait — une
recherche par `context` ne change aucune signature de composant — donc rien ne
justifie de le faire d'avance. `formatAddress` commute déjà sur `CountryISO`,
là où le français aurait produit des étiquettes illisibles par un transporteur ;
c'est le seul endroit où ça comptait.

## Ce qu'il ne faut pas construire ici

Une boutique publique — panier, tunnel de paiement, galerie produit — est un
troisième produit : mobile-first, état client, PCI. Un éditeur riche, un Kanban
avec glisser-déposer, une bibliothèque de graphiques : chacun est un projet
JavaScript déguisé en composant.
