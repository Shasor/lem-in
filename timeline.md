## Gestion Fichier (argument):

- ouvrir le fichier
- lire ligne par ligne

- ### Interpréter:

  - premier ligne = nbr de fourmis (if int && > 0 && < 1001)
  - ligne commençant par "##":

    - suivi de "start" = ligne suivante représente la première room (point de départ de toutes les fourmis)
    - suivi de "end" = ligne suivante représente la dernière room (point d'arrivé de toutes les fourmis)

  - ligne commençant par "#" = commentaire

  - #### une room est défini par:

    - une ligne avec deux " " (espace) = 3 éléments:

      - élément 1: alpha-numérique = nom de la room
      - 2 et 3e élément = numérique (0-9) à convertir en int = coordonnées x, y de la room

  - #### un tunnel/link est défini par:

    - une ligne avec un seul "-" = 2 éléments:

      - chaque élément doit être = string et = au nom d'une des rooms

- ### stocker le tout dans des structures

- ### retourner une erreur si:

  - la longueur arguments != 2
  - le format du fichier != .txt
  - le nom d'une room commance par "L"
  - les coordonnées ne sont pas int
  - il apparaît plus d'une ligne "##start" ou "##end"
  - aucune ligne avec "##start" ou "##end"
  - un tunnel/link a deux éléments identiques (ex: 2-2)
  - un tunnel/link apparaît deux fois (ex: 2-3, 3-2)

## Algo

### Règles:

- une room ne peut contenir qu'une seule fourmi à la fois, sauf ##start et ##end

- chaque fourmi doit partir de ##start est arriver à ##end en moins de tour possible

- une fourmi ne peut avancer que d'une seule room/tunnel par tour

### Logique:

- un tour est défini par tous les mouvements possibles à la suite du mouvement des fourmis plus avancées

- trouver tous les chemins possibles, définir ceux les plus courts

- définir les chemins sans conflits à partir du nombre de rooms liées à ##start (si < : trouver une autre stratégie)

## Output

- Afficher le fichier de départ tel quel
- puis, une ligne par tour, afficher tous les mouvements réalisés, selon cette syntaxe :

  - "Tour 1: LA-R LA-R" où A est l'index de la fourmi et R le nom de la room d'arrivé

    ex: "Tour 1: L1-3 L2-2", signifie que la fourmi 1 va dans la room nommée "3" et la fourmi 2 dans la room 2
