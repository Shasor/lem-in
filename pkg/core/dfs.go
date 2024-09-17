package core

func dfs(current, end int, path []int, visited []bool, paths *[][]string) {
	visited[current] = true
	path = append(path, current)

	if current == end {
		// Convertir les indices en noms de salles pour le chemin final
		namePath := make([]string, len(path))
		for i, idx := range path {
			namePath[i] = Rooms[idx].Name
		}
		*paths = append(*paths, namePath)
	} else {
		for _, linkName := range Rooms[current].Links {
			nextIndex := findRoomIndex(Rooms, linkName)
			if nextIndex != -1 && !visited[nextIndex] {
				dfs(nextIndex, end, path, visited, paths)
			}
		}
	}

	visited[current] = false
	// path = path[:len(path)-1]
}

func findRoomIndex(Rooms []Room, name string) int {
	for i, room := range Rooms {
		if room.Name == name {
			return i
		}
	}
	return -1
}
