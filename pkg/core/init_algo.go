package core

import (
	"fmt"
)

func AlgoInit() {
	Paths = FindPaths()
	// Obtenez le nombre de fourmis et la salle de départ
	antCount := len(Ants)
	startRoom := GetStartRoom()

	// Optimisez la sélection des chemins
	bestCombination := OptimizePathSelection(Paths, antCount, startRoom)

	if len(bestCombination) > 0 {
		SimulateAntMovement(bestCombination)
	} else {
		fmt.Println("Aucun chemin optimal trouvé.")
	}
}
