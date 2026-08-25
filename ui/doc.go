// Package ui fournit les composants d'un back-office server-rendered :
// templ pour le rendu, daisyUI pour les classes, Unpoly pour éviter les
// rechargements.
//
// Tout état de liste vit dans la query string. Les helpers SortHref, PageHref
// et ResetHref la construisent ; votre handler l'applique aux données. Le
// paquet ne stocke rien.
//
// Aucun changement d'état ne passe par un GET : les actions destructives sont
// des boutons submit, jamais des liens, pour que http.NewCrossOriginProtection
// les couvre.
//
// La galerie rendue de tous les composants se construit avec
// « go run ./cmd/gallery ».
package ui
