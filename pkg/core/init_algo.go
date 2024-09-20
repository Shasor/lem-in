package core

func AlgoInit() {
	var start, end int

	// Trouver les index des salles de départ et d'arrivée
	for i, room := range Rooms {
		if room.IsStart {
			start = i
		}
		if room.IsEnd {
			end = i
		}
	}

	visited := make([]bool, len(Rooms))
	dfs(start, end, []int{}, visited, &Paths)

	SortPaths()
}
