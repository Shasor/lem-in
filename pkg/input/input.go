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

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fmt.Println(line)

		if line == "" || strings.HasPrefix(line, "#") && line != "##start" && line != "##end" {
			continue
		}

		if line == "##start" {
			nextIsStart = true
			continue
		}
		if line == "##end" {
			nextIsEnd = true
			continue
		}

		CreateAnts(line)

		CreateRooms(line)

		// Add all their links to each room
		AddLinks(line)
	}
	fmt.Printf("\n")
	core.AlgoInit()
}
