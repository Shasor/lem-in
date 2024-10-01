package visualizer

import (
	"bufio"
	"fmt"
	"io"
	"lem-in/pkg/core"
	"lem-in/pkg/input"
	"os"
	"strconv"
	"strings"
)

// ParseInput analyse l'entrée et retourne un GameState
func ParseInput(r io.Reader) (*GameState, error) {
	scanner := bufio.NewScanner(r)
	state := &GameState{
		Rooms:     []core.Room{},
		Movements: [][]AntMove{},
	}

	// movementStarted := false
	var nextIsStart, nextIsEnd bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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

		// Lire le nombre de fourmis
		if line != "0" {
			if len(core.Ants) == 0 {
				numAnts, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Println("ERROR: invalid data format")
					os.Exit(1)
				}
				for i := 1; i <= numAnts; i++ {
					input.CreateAnt(i)
				}
				state.NumAnts = numAnts
			}
		} else {
			fmt.Println("ERROR: invalid data format")
			os.Exit(0)
		}

		// Lire les lignes contenant les informations sur les salles
		if strings.Contains(line, " ") && !strings.Contains(line, ":") {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				fmt.Println("ERROR: invalid data format")
				os.Exit(0)
			}

			name := parts[0]
			x, err := strconv.Atoi(parts[1])
			core.ErrorsHandler(err)
			y, err := strconv.Atoi(parts[2])
			core.ErrorsHandler(err)
			if strings.HasPrefix(name, "L") {
				fmt.Println("ERROR: invalid room name")
				os.Exit(0)
			}

			if nextIsStart {
				input.CreateRoom(name, x, y, nextIsStart, false)
				nextIsStart = false
			} else if nextIsEnd {
				input.CreateRoom(name, x, y, false, nextIsEnd)
				nextIsEnd = false
			} else {
				input.CreateRoom(name, x, y, false, false)
			}
		}

		input.AddLinks(line)
		state.Rooms = core.Rooms

		if strings.Contains(line, "L") && strings.Contains(line, ":") {
			turnNumber, moves := parseMoveLine(line)
			if len(moves) > 0 {
				// Assurez-vous que le slice Movements a suffisamment d'espace
				for len(state.Movements) <= turnNumber-1 {
					state.Movements = append(state.Movements, nil)
				}
				state.Movements[turnNumber-1] = moves
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de l'entrée standard: %v", err)
	}
	return state, nil
}

// parseMoveLine analyse une ligne de mouvement
func parseMoveLine(line string) (int, []AntMove) {
	moves := []AntMove{}
	parts := strings.Split(line, " ")

	// Initialiser le numéro de tour
	turnNumber := 0
	startIndex := 0

	// Vérifier si le premier élément est un numéro de tour
	if len(parts) > 0 && strings.HasSuffix(parts[0], ":") {
		turnStr := strings.TrimSuffix(parts[0], ":")
		if num, err := strconv.Atoi(turnStr); err == nil {
			turnNumber = num
			startIndex = 1 // Commencer à partir du deuxième élément si un numéro de tour est présent
		}
	}

	// Parcourir les mouvements, en commençant par le bon index
	for _, part := range parts[startIndex:] {
		if strings.HasPrefix(part, "L") {
			moveParts := strings.Split(part, "-")
			if len(moveParts) == 2 {
				antID, _ := strconv.Atoi(strings.TrimPrefix(moveParts[0], "L"))
				moves = append(moves, AntMove{
					AntID: antID,
					Room:  moveParts[1],
				})
			}
		}
	}
	return turnNumber, moves
}
