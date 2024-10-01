package input

import (
	"bufio"
	"fmt"
	"lem-in/pkg/core"
	"os"
	"strings"
)

var nextIsStart, nextIsEnd bool

func InputHandler(file string) {
	content, err := os.Open(file)
	core.ErrorsHandler(err)
	defer content.Close()

	scanner := bufio.NewScanner(content)
	var history string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") && line != "##start" && line != "##end" {
			history += line + "\n"
			continue
		}

		if line == "##start" {
			nextIsStart = true
			history += line + "\n"
			continue
		}
		if line == "##end" {
			nextIsEnd = true
			history += line + "\n"
			continue
		}

		CreateAnts(line)

		CreateRooms(line)

		// Add all their links to each room
		AddLinks(line)
		history += line + "\n"
	}
	fmt.Println(history)
	core.AlgoInit()
}
