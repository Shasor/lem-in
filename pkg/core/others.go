package core

import "math"

func getRoomByName(name string) *Room {
	for i := range Rooms {
		if Rooms[i].Name == name {
			return &Rooms[i]
		}
	}
	return nil
}

func getStartRoom() *Room {
	for i := range Rooms {
		if Rooms[i].IsStart {
			return &Rooms[i]
		}
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func removeIndex(s [][]string, index int) [][]string {
	return append(s[:index], s[index+1:]...)
}

func countCommonRooms(path1, path2 []string) int {
	common := 0
	set := make(map[string]bool)

	// Ignorer la première et la dernière salle
	for _, room := range path1[1 : len(path1)-1] {
		set[room] = true
	}

	for _, room := range path2[1 : len(path2)-1] {
		if set[room] {
			common++
		}
	}

	return common
}

func countConflicts(path []string, existingPaths [][]string) int {
	conflicts := 0
	for _, existingPath := range existingPaths {
		conflicts += countCommonRooms(path, existingPath)
	}
	return conflicts
}

func hasSameMiddle(path1, path2 []string) bool {
	set := make(map[string]bool)

	// Ajouter tous les éléments de path1 à l'ensemble
	for _, elem := range path1 {
		set[elem] = true
	}

	// Vérifier si un élément de path2 est dans l'ensemble
	for _, elem := range path2 {
		if set[elem] {
			return true
		}
	}

	return false
}

func hasConflicts(paths [][]string) bool {
	seen := make(map[string]bool)
	for _, path := range paths {
		for _, room := range path[1 : len(path)-1] { // Ignorer start et end
			if seen[room] {
				return true
			}
			seen[room] = true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateCombinations(paths [][]string, k int) [][][]string {
	var result [][][]string
	var current [][]string
	generateCombinationsHelper(paths, k, 0, current, &result)
	return result
}

func selectBestCombo(combos [][][]string, antCount int) [][]string {
	bestScore := math.MaxFloat64
	var bestCombo [][]string

	for _, combo := range combos {
		score := evaluateCombo(combo, antCount)
		if score < bestScore {
			bestScore = score
			bestCombo = combo
		}
	}

	return bestCombo
}

func evaluateCombo(combo [][]string, antCount int) float64 {
	lengths := make([]float64, len(combo))
	for i, path := range combo {
		lengths[i] = float64(len(path) - 1)
	}

	mean := 0.0
	for _, l := range lengths {
		mean += l
	}
	mean /= float64(len(lengths))

	variance := 0.0
	for _, l := range lengths {
		variance += math.Pow(l-mean, 2)
	}
	variance /= float64(len(lengths))

	// Un score plus bas est meilleur
	return variance + mean/float64(antCount)
}
