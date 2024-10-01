package input

import (
	"fmt"
	"lem-in/pkg/core"
	"os"
	"strconv"
	"strings"
)

func CreateAnts(line string) {
	if line != "0" {
		if len(core.Ants) == 0 {
			nbr_ants, err := strconv.Atoi(line)
			core.ErrorsHandler(err)
			for i := 1; i <= nbr_ants; i++ {
				CreateAnt(i)
			}
		}
	} else {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}
}

func CreateRooms(line string) {
	if strings.Contains(line, " ") {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			fmt.Println("ERROR: invalid data format")
			os.Exit(1)
		}
		name := parts[0]
		x, err := strconv.Atoi(parts[1])
		core.ErrorsHandler(err)
		y, err := strconv.Atoi(parts[2])
		core.ErrorsHandler(err)
		if strings.HasPrefix(name, "L") {
			fmt.Println("ERROR: invalid room name")
			os.Exit(1)
		}
		if nextIsStart {
			CreateRoom(name, x, y, nextIsStart, false)
			nextIsStart = false
		} else if nextIsEnd {
			CreateRoom(name, x, y, false, nextIsEnd)
			nextIsEnd = false
		} else {
			CreateRoom(name, x, y, false, false)
		}
	}
}

func AddLinks(line string) {
	if strings.Contains(line, "-") && !strings.Contains(line, ":") {
		parts := strings.Split(line, "-")
		if parts[0] == parts[1] {
			fmt.Println("ERROR: invalid link syntax")
			os.Exit(1)
		}

		// Check if room name exist
		var count0, count1 int
		for _, room := range core.Rooms {
			if parts[0] == room.Name {
				count0++
			}
			if parts[1] == room.Name {
				count1++
			}
		}
		if count0 == 0 || count1 == 0 {
			fmt.Println("give me a valid link motherfucker!")
			os.Exit(1)
		}

		for i := range core.Rooms {
			if core.Rooms[i].Name == parts[0] {
				core.Rooms[i].Links = append(core.Rooms[i].Links, parts[1])
			}
			if core.Rooms[i].Name == parts[1] {
				core.Rooms[i].Links = append(core.Rooms[i].Links, parts[0])
			}
		}
	}
}

func CreateRoom(name string, x, y int, isStart, isEnd bool) {
	r := core.Room{
		Name:    name,
		X:       x,
		Y:       y,
		IsStart: isStart,
		IsEnd:   isEnd,
	}
	core.Rooms = append(core.Rooms, r)
}

func CreateAnt(i int) {
	a := core.Ant{
		Index: i,
		Path:  []string{},
	}
	core.Ants = append(core.Ants, a)
}
