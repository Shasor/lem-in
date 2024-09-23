package core

import (
	"fmt"
	"math"
	"sort"
)

func FindPaths() [][]string {
	var startRoom, endRoom string

	// Trouver les salles de départ et d'arrivée
	for _, room := range Rooms {
		if room.IsStart {
			startRoom = room.Name
		}
		if room.IsEnd {
			endRoom = room.Name
		}
	}

	if startRoom == "" || endRoom == "" {
		ErrorsHandler(fmt.Errorf("start or end room not found"))
	}

	// Initialiser avec le chemin de départ
	Paths = append(Paths, []string{startRoom})

	// Boucle principale pour étendre les chemins
	for {
		newPaths := [][]string{}
		pathsToRemove := []int{}

		for i, path := range Paths {
			lastRoom := path[len(path)-1]

			if lastRoom == endRoom {
				continue // Ce chemin est complet
			}

			currentRoom := getRoomByName(lastRoom)
			if currentRoom == nil {
				continue
			}

			for _, link := range currentRoom.Links {
				if contains(path, link) && link != endRoom {
					continue // Éviter les boucles
				}

				newPath := make([]string, len(path))
				copy(newPath, path)
				newPath = append(newPath, link)
				newPaths = append(newPaths, newPath)
			}

			pathsToRemove = append(pathsToRemove, i)
		}

		// Supprimer les chemins qui ont été étendus
		for i := len(pathsToRemove) - 1; i >= 0; i-- {
			Paths = removeIndex(Paths, pathsToRemove[i])
		}

		// Ajouter les nouveaux chemins
		Paths = append(Paths, newPaths...)

		// Vérifier si tous les chemins se terminent par la salle de fin
		allComplete := true
		for _, path := range Paths {
			if path[len(path)-1] != endRoom {
				allComplete = false
				break
			}
		}

		if allComplete {
			break
		}
	}

	return Paths
}

func FilterUniquePaths(paths [][]string) [][]string {
	var uniquePaths [][]string

	for i, path := range paths {
		if len(path) < 3 {
			// Ignorer les chemins trop courts
			continue
		}

		hasUniquePath := false
		for j, otherPath := range paths {
			if i == j || len(otherPath) < 3 {
				continue
			}

			// Si le deuxième élément est le même, ne pas comparer ces chemins
			if path[1] == otherPath[1] {
				continue
			}

			// Comparer les éléments du milieu (en excluant le premier et le dernier)
			if !hasSameMiddle(path[1:len(path)-1], otherPath[1:len(otherPath)-1]) {
				hasUniquePath = true
				break
			}
		}

		if hasUniquePath {
			uniquePaths = append(uniquePaths, path)
		}
	}

	return uniquePaths
}

func FindOptimalPaths() [][][]string {
	Paths := FindPaths()
	nbrChemOpt := findOptimalPathCount()

	var optimalPaths [][][]string

	for i := 1; i <= len(Paths); i++ {
		uniquePaths := filterPathsByUniqueCount(Paths, i)
		if len(uniquePaths) > 0 {
			optimalPaths = append(optimalPaths, uniquePaths)
		}

		// Si nous avons atteint exactement nbrChemOpt chemins sans éléments communs, nous arrêtons
		if len(uniquePaths) == nbrChemOpt {
			break
		}
	}

	return optimalPaths
}

func FindOptimalPaths2() [][][]string {
	Paths := FindPaths()
	startRoom := getStartRoom()
	if startRoom == nil {
		fmt.Println("Salle de départ non trouvée.")
		return nil
	}

	// Grouper les chemins par leur première salle après le départ
	pathsByFirstRoom := make(map[string][][]string)
	for _, path := range Paths {
		if len(path) > 1 {
			firstRoom := path[1]
			pathsByFirstRoom[firstRoom] = append(pathsByFirstRoom[firstRoom], path)
		}
	}

	optimalPaths := make([][]string, 0, len(startRoom.Links))

	// Pour chaque salle liée à la salle de départ
	for _, link := range startRoom.Links {
		paths := pathsByFirstRoom[link]
		if len(paths) == 0 {
			continue
		}

		// Trouver le chemin avec le moins de conflits
		bestPath := findLeastConflictingPath(paths, optimalPaths)
		optimalPaths = append(optimalPaths, bestPath)
	}

	// Trier les chemins optimaux par ordre de taille (du plus court au plus long)
	sort.Slice(optimalPaths, func(i, j int) bool {
		return len(optimalPaths[i]) < len(optimalPaths[j])
	})

	antCount := len(Ants)
	bestCombination := OptimizePathSelection(optimalPaths, antCount, startRoom)

	return [][][]string{bestCombination}
}

func findLeastConflictingPath(paths [][]string, existingPaths [][]string) []string {
	if len(paths) == 0 {
		return nil
	}

	var bestPath []string
	minConflicts := -1

	for _, path := range paths {
		conflicts := countConflicts(path, existingPaths)
		if minConflicts == -1 || conflicts < minConflicts {
			minConflicts = conflicts
			bestPath = path
		}
	}

	return bestPath
}

func findOptimalPathCount() int {
	var startLinks, endLinks int

	for _, room := range Rooms {
		if room.IsStart {
			startLinks = len(room.Links)
		}
		if room.IsEnd {
			endLinks = len(room.Links)
		}
	}

	return min(startLinks, endLinks)
}

func filterPathsByUniqueCount(paths [][]string, minUniqueCount int) [][]string {
	var uniquePaths [][]string

	for i, path := range paths {
		if len(path) < 3 {
			continue
		}

		uniqueCount := 0
		for j, otherPath := range paths {
			if i == j || len(otherPath) < 3 {
				continue
			}

			if path[1] == otherPath[1] {
				continue
			}

			if !hasSameMiddle(path[1:len(path)-1], otherPath[1:len(otherPath)-1]) {
				uniqueCount++
			}
		}

		if uniqueCount >= minUniqueCount {
			uniquePaths = append(uniquePaths, path)
		}
	}

	return uniquePaths
}

func OptimizePathSelection(Paths [][]string, antCount int, startRoom *Room) [][]string {

	// Trier tous les chemins par longueur
	sort.Slice(Paths, func(i, j int) bool {
		return len(Paths[i]) < len(Paths[j])
	})

	shortestPath := Paths[0]
	compatiblePaths := [][]string{shortestPath}

	// Rechercher tous les chemins compatibles avec le plus court
	for _, path := range Paths[1:] {
		if !hasConflicts(append(compatiblePaths, path)) {
			compatiblePaths = append(compatiblePaths, path)
		}
	}

	maxPaths := len(startRoom.Links)
	var globalBestCombinations [][][]string
	globalMinTurns := math.MaxInt32

	for numPaths := maxPaths; numPaths >= 1; numPaths-- {
		combinations := generateCombinations(Paths, numPaths)
		var bestCombinations [][][]string
		minTurns := math.MaxInt32

		for _, combo := range combinations {
			if !hasConflicts(combo) {
				turns := calculateTurns(combo, antCount)

				if turns < minTurns {
					minTurns = turns
					bestCombinations = [][][]string{combo}
				} else if turns == minTurns {
					bestCombinations = append(bestCombinations, combo)
				}
			}
		}

		if len(bestCombinations) > 0 {

			// Limiter à deux combinaisons au maximum
			maxCombos := 2
			if len(bestCombinations) < maxCombos {
				maxCombos = len(bestCombinations)
			}

			// for i := 0; i < maxCombos; i++ {
			// 	fmt.Printf("Combinaison %d:\n", i+1)
			// 	for j, path := range bestCombinations[i] {
			// 		fmt.Printf("  Chemin %d: %v\n", j+1, path)
			// 	}
			// 	fmt.Printf("Nombre de tours : %d\n\n", minTurns)
			// }

			if minTurns < globalMinTurns {
				globalMinTurns = minTurns
				globalBestCombinations = bestCombinations[:maxCombos]
			} else if minTurns == globalMinTurns && numPaths < len(globalBestCombinations[0]) {
				globalBestCombinations = bestCombinations[:maxCombos]
			}
		}
	}

	// Retourner la première combinaison parmi les meilleures
	return selectBestCombo(globalBestCombinations, antCount)
}

func generateCombinationsHelper(paths [][]string, k, start int, current [][]string, result *[][][]string) {
	if len(current) == k {
		combination := make([][]string, len(current))
		copy(combination, current)
		*result = append(*result, combination)
		return
	}

	for i := start; i < len(paths); i++ {
		current = append(current, paths[i])
		generateCombinationsHelper(paths, k, i+1, current, result)
		current = current[:len(current)-1]
	}
}

func calculateTurns(paths [][]string, antCount int) int {
	if len(paths) == 0 {
		return 0
	}

	pathLengths := make([]int, len(paths))
	for i, path := range paths {
		pathLengths[i] = len(path) - 1
	}

	// Trier les chemins du plus court au plus long
	sort.Ints(pathLengths)

	totalTurns := 0
	for i, length := range pathLengths {
		remainingAnts := antCount - i
		if remainingAnts <= 0 {
			break
		}
		turnsForPath := length + (remainingAnts-1)/len(paths)
		if turnsForPath > totalTurns {
			totalTurns = turnsForPath
		}
	}

	return totalTurns
}
