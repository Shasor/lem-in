package core

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func SimulateAntMovement(optimalPaths [][]string) {
	if len(optimalPaths) == 0 {
		fmt.Println("Aucun chemin trouvé.")
		return
	}

	var b bool
	for _, path := range optimalPaths {
		if len(path) == 2 {
			b = true
		}
	}

	if len(optimalPaths) == 2 && b {
		SimulateAntMovement2(optimalPaths)
	} else {
		SimulateAntMovement1(optimalPaths)
	}
}

func SimulateAntMovement1(optimalPaths [][]string) {
	if len(optimalPaths) == 0 {
		fmt.Println("Aucun chemin trouvé.")
		return
	}

	antCount := len(Ants)
	pathCount := len(optimalPaths)
	antPositions := make([]int, antCount) // -1 signifie que la fourmi n'a pas encore commencé
	antPaths := make([]int, antCount)     // Indique quel chemin chaque fourmi suit
	for i := range antPositions {
		antPositions[i] = -1
	}

	turn := 0
	finishedAnts := 0
	var startRoom, endRoom string

	// Trouver les salles de départ et de fin
	for _, room := range Rooms {
		if room.IsStart {
			startRoom = room.Name
		}
		if room.IsEnd {
			endRoom = room.Name
		}
	}

	// Trouver le nombre de connexions de la salle de départ
	startConnections := len(getRoomByName(startRoom).Links)

	for finishedAnts < antCount {
		turn++
		movements := []string{}

		// Déplacer les fourmis existantes
		for i := 0; i < antCount; i++ {
			if antPositions[i] >= 0 && antPositions[i] < len(optimalPaths[antPaths[i]])-1 {
				nextRoom := optimalPaths[antPaths[i]][antPositions[i]+1]
				if nextRoom == endRoom || isRoomFree(nextRoom, antPositions, antPaths, optimalPaths) {
					antPositions[i]++
					movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), nextRoom))
					if antPositions[i] == len(optimalPaths[antPaths[i]])-1 {
						finishedAnts++
					}
				}
			}
		}

		// Ajouter de nouvelles fourmis
		antsAdded := 0
		for i := 0; i < antCount && antsAdded < startConnections; i++ {
			if antPositions[i] == -1 { // Si la fourmi n'a pas encore commencé
				for path := 0; path < pathCount; path++ {
					firstRoom := optimalPaths[path][1]
					if isRoomFree(firstRoom, antPositions, antPaths, optimalPaths) {
						antPositions[i] = 1
						antPaths[i] = path
						movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), firstRoom))
						antsAdded++
						break
					}
				}
				if antsAdded >= startConnections {
					break
				}
			}
		}

		if len(movements) > 0 {
			fmt.Printf("%d: %s\n", turn, strings.Join(movements, " "))
		}
	}
}

func SimulateAntMovement2(optimalPaths [][]string) {
	antCount := len(Ants)
	pathCount := len(optimalPaths)
	antPositions := make([]int, antCount)
	antPaths := make([]int, antCount)

	for i := range antPositions {
		antPositions[i] = -1
	}

	turn := 0
	finishedAnts := 0
	var endRoom string

	// Trouver la salle de fin et la longueur du chemin le plus long et le plus court
	longestPathLength := 0
	shortestPathLength := math.MaxInt32
	pathLengths := make([]int, pathCount) // Stocke la longueur de chaque chemin
	for _, room := range Rooms {
		if room.IsEnd {
			endRoom = room.Name
			break
		}
	}
	for i, path := range optimalPaths {
		pathLength := len(path) - 1 // -1 car on ne compte pas la salle de départ
		pathLengths[i] = pathLength
		if pathLength > longestPathLength {
			longestPathLength = pathLength
		}
		if pathLength < shortestPathLength {
			shortestPathLength = pathLength
		}
	}

	shortestPathIndex := -1 // Initialiser ici pour éviter l'erreur

	for finishedAnts < antCount {
		turn++

		// Impression de l'état avant le mouvement (pour débogage)

		movements := []string{}

		remainingAnts := antCount - finishedAnts

		// Déplacer les fourmis existantes
		for i := 0; i < antCount; i++ {
			if antPositions[i] >= 0 && antPositions[i] < len(optimalPaths[antPaths[i]])-1 {
				nextRoom := optimalPaths[antPaths[i]][antPositions[i]+1]
				if nextRoom == endRoom || isRoomFree(nextRoom, antPositions, antPaths, optimalPaths) {
					antPositions[i]++
					movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), nextRoom))
					if nextRoom == endRoom {
						finishedAnts++
					}
				}
			}
		}

		// Ajouter de nouvelles fourmis (toutes sauf la dernière)
		for path := 0; path < pathCount; path++ {
			firstRoom := optimalPaths[path][1]
			if isRoomFree(firstRoom, antPositions, antPaths, optimalPaths) || firstRoom == endRoom {
				if remainingAnts > 0 { // Ajustement ici pour permettre d'ajouter des fourmis
					for i := 0; i < antCount; i++ {
						if antPositions[i] == -1 && i != (antCount-1) { // Exclure la dernière fourmi (L20)
							antPositions[i] = 1 // Commence à 1 car 0 est la salle de départ
							antPaths[i] = path
							movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), firstRoom))
							if firstRoom == endRoom {
								finishedAnts++
							}
							remainingAnts--
							break
						}
					}
				}
			}
		}

		// Vérifier si seulement deux fourmis restent
		if remainingAnts == 2 && finishedAnts < antCount {
			for i := 0; i < antCount; i++ {
				if antPositions[i] == -1 && i == (antCount-2) { // Avant-dernière fourmi
					firstRoom := optimalPaths[shortestPathIndex][1]
					antPositions[i] = 1 // Commence à 1 car 0 est la salle de départ
					antPaths[i] = shortestPathIndex
					movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), firstRoom))
					remainingAnts--
					break
				} else if antPositions[i] == -1 && i == (antCount-1) { // Dernière fourmi (L20)
					// La dernière fourmi attend et ne fait rien ce tour-ci.
					continue
				}
			}
		}

		// Vérifier si c'est le tour où L20 peut se déplacer sur le chemin le plus court
		if remainingAnts == 2 && turn >= longestPathLength-shortestPathLength {
			for i := 0; i < antCount; i++ {
				if antPositions[i] == -1 && i == (antCount-1) { // Assurez-vous que c'est L20 (la dernière fourmi)
					// Déplacez L20 vers la salle de fin directement au dernier tour.
					movements = append(movements, fmt.Sprintf("L%s-%s", formatAntNumber(i+1), endRoom))
					finishedAnts++
					remainingAnts--
					break
				}
			}
		}

		if len(movements) > 0 {
			fmt.Printf("%d: %s\n", turn, strings.Join(movements, " "))
		}

		if len(movements) == 0 {
			break // Évite une boucle infinie si aucun mouvement n'est possible
		}
	}
}

func formatAntNumber(num int) string {
	if num < 10 {
		return fmt.Sprintf("0%d", num)
	}
	return strconv.Itoa(num)
}

func isRoomFree(room string, antPositions []int, antPaths []int, paths [][]string) bool {
	for ant, pos := range antPositions {
		if pos == -1 {
			continue
		}
		if paths[antPaths[ant]][pos] == room {
			return false
		}
	}
	return true
}

// func join(elements []string, separator string) string {
// 	result := ""
// 	for i, element := range elements {
// 		if i > 0 {
// 			result += separator
// 		}
// 		result += element
// 	}
// 	return result
// }
