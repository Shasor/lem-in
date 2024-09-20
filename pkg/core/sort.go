package core

import (
	"fmt"
	"sort"
)

func SortPaths() {
	sort.Slice(Paths, func(i, j int) bool {
		return len(Paths[i]) < len(Paths[j])
	})

	var result [][]string
	for _, path := range Paths {
		if !isConflict(path, result) {
			result = append(result, path)
		}
	}

	for _, path := range result {
		fmt.Println(path)
	}
}

func isConflict(path []string, result [][]string) bool {
	for _, res := range result {
		for i := 1; i < len(res)-1; i++ {
			for j := 1; j < len(path)-1; j++ {
				if path[j] == res[i] {
					return true
				}
			}
		}
	}
	return false
}
