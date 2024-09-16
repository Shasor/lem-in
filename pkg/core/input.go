package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var Rooms []Room
var Ants []Ant

func InputHandler(file string) ([]Room, error) {
	content, err := os.Open(file)
	ErrorsHandler(err)
	defer content.Close()

	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") && line != "##start" && line != "##end" {
			continue
		}

		// Add Ants
		if len(Ants) == 0 {
			nbr_ants, err := strconv.Atoi(line)
			ErrorsHandler(err)
			for i := 1; i <= nbr_ants; i++ {
				CreateAnt(i)
			}
		}

		if line == "##start" {
			nextIsStart := true
			continue
		}
		if line == "##end" {
			nextIsEnd := true
			continue
		}

		if strings.Contains(line, " ") {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				fmt.Println("ERROR: invalid data format")
				os.Exit(1)
			}
			name := parts[0]
			x, err := strconv.Atoi(parts[1])
			ErrorsHandler(err)
			y, err := strconv.Atoi(parts[2])
			ErrorsHandler(err)
			if strings.HasPrefix(name, "L") {
				fmt.Println("ERROR: invalid room name")
				os.Exit(1)
			}
			CreateRoom(name, x, y)
		}
	}
	fmt.Println(Rooms)
	return Rooms, nil
}

func CreateRoom(name string, x, y int) {
	r := Room{
		Name: name,
		X:    x,
		Y:    y,
	}
	Rooms = append(Rooms, r)
}

func CreateAnt(i int) {
	a := Ant{
		Index: i,
		Path:  []string{},
	}
	Ants = append(Ants, a)
}
