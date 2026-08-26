// Package marketing fournit les composants des pages publiques : hero,
// grille de fonctionnalités, tarifs, témoignages, FAQ, pied de page.
//
// Il est séparé du paquet ui pour une raison : un back-office et une landing
// page n'ont pas les mêmes contraintes. Le premier est dense, authentifié,
// piloté par la donnée, sans référencement ni expression de marque. La seconde
// vit du référencement, des images et de la conversion.
//
// Trois règles s'appliquent ici et pas dans ui :
//
//   - Un seul <h1> par page, rendu par Hero. Les autres sections ouvrent en
//     <h2>, sinon la hiérarchie des titres ment aux moteurs et aux lecteurs
//     d'écran.
//   - Toute image porte alt, width et height. Sans dimensions, la page saute au
//     chargement et le Cumulative Layout Shift s'effondre.
//   - Aucun JavaScript. La FAQ utilise <details>, pas un accordéon scripté.
package marketing
