package visualizer

import (
	"fmt"
	"image/color"
	"lem-in/pkg/core"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/rand"
)

// findRoom trouve une salle par son nom
func (a *App) findRoom(name string) (core.Room, bool) {
	for _, room := range a.gameState.Rooms {
		if room.Name == name {
			return room, true
		}
	}
	return core.Room{}, false
}

// getStartRoomName retourne le nom de la salle de départ
func (a *App) getStartRoomName() string {
	for _, room := range a.gameState.Rooms {
		if room.IsStart {
			return room.Name
		}
	}
	return ""
}

// allAntsAtEnd vérifie si toutes les fourmis sont arrivées à la fin
func (a *App) allAntsAtEnd() bool {
	endRoom := ""
	for _, room := range a.gameState.Rooms {
		if room.IsEnd {
			endRoom = room.Name
			break
		}
	}
	if endRoom == "" {
		return false
	}

	for _, ant := range a.ants {
		if ant.Room != endRoom {
			return false
		}
	}
	return true
}

// createAnts initialise les fourmis
func (a *App) createAnts() {
	startRoomName := a.getStartRoomName()
	startRoom, exists := a.findRoom(startRoomName)
	if !exists {
		fmt.Printf("Erreur : salle de départ '%s' non trouvée\n", startRoomName)
		return
	}

	a.ants = make([]Ant, a.gameState.NumAnts)
	for i := range a.ants {
		a.ants[i] = Ant{
			ID:    i + 1,
			X:     float32(startRoom.X),
			Y:     float32(startRoom.Y),
			Room:  startRoom.Name,
			Color: color.RGBA{uint8(rand.Intn(200) + 55), uint8(rand.Intn(200) + 55), uint8(rand.Intn(200) + 55), 255},
		}
	}
}

// generateEarthImage crée l'image de la terre
func (a *App) generateEarthImage() {
	earthHeight := int(float64(a.height) * 0.75) // 3/4 de la hauteur totale
	a.earthImage = ebiten.NewImage(a.width, earthHeight)
	darkEarth := color.RGBA{180, 140, 100, 255}
	mediumEarth := color.RGBA{200, 160, 120, 255}
	lightEarth := color.RGBA{220, 180, 140, 255}

	for y := 0; y < earthHeight; y += 10 {
		for x := 0; x < a.width; x += 10 {
			var c color.Color
			r := rand.Float32()
			if r < 0.5 {
				c = darkEarth
			} else if r < 0.83 {
				c = mediumEarth
			} else {
				c = lightEarth
			}

			rVar := float64(rand.Intn(21) - 10)
			gVar := float64(rand.Intn(21) - 10)
			bVar := float64(rand.Intn(21) - 10)
			rc, gc, bc, ac := c.RGBA()
			rc = uint32(clamp(int(rc>>8)+int(rVar), 0, 255))
			gc = uint32(clamp(int(gc>>8)+int(gVar), 0, 255))
			bc = uint32(clamp(int(bc>>8)+int(bVar), 0, 255))
			variedColor := color.RGBA{uint8(rc), uint8(gc), uint8(bc), uint8(ac >> 8)}

			ebitenutil.DrawRect(a.earthImage, float64(x), float64(y), 10, 10, variedColor)
		}
	}
}

// generateGrassImage crée l'image de l'herbe
func (a *App) generateGrassImage() {
	grassHeight := int(float64(a.height) * 0.15) // 15% de la hauteur totale pour l'herbe
	a.grassImage = ebiten.NewImage(a.width, grassHeight)
	grassColor := color.RGBA{34, 139, 34, 180} // Vert foncé pour l'herbe avec transparence
	a.grassImage.Fill(grassColor)

	// Ajouter des détails d'herbe
	for i := 0; i < 2000; i++ { // Augmenté le nombre de brins d'herbe
		x := rand.Float64() * float64(a.width)
		y := rand.Float64() * float64(grassHeight)
		length := 5 + rand.Float64()*10
		grassShade := color.RGBA{0, uint8(100 + rand.Intn(100)), 0, 255}

		vector.StrokeLine(a.grassImage, float32(x), float32(y), float32(x), float32(y-length), 1, grassShade, false)
	}
}

// generateSkyImage crée l'image du ciel
func (a *App) generateSkyImage() {
	skyHeight := int(float64(a.height) * 0.10) // 10% restant pour le ciel
	a.skyImage = ebiten.NewImage(a.width, skyHeight)
	skyColor := color.RGBA{135, 206, 235, 255} // Bleu ciel
	a.skyImage.Fill(skyColor)
}

// startCelebration initialise la célébration
func (a *App) startCelebration() {
	a.celebrationLights = make([]CelebrationLight, 0)
}

// updateCelebration met à jour l'état de la célébration
func (a *App) updateCelebration() {
	endRoom := a.getEndRoom()

	// Déplacer les fourmis aléatoirement
	for i := range a.ants {
		a.ants[i].X += float32(rand.Float32()*4 - 2)
		a.ants[i].Y += float32(rand.Float32()*4 - 2)
	}

	// Ajouter de nouvelles lumières
	if rand.Float32() < 0.3 {
		light := CelebrationLight{
			X: float32(endRoom.X) + float32(rand.Float32()*40-20),
			Y: float32(endRoom.Y) + float32(rand.Float32()*40-20),
			Color: color.RGBA{
				uint8(rand.Intn(256)),
				uint8(rand.Intn(256)),
				uint8(rand.Intn(256)),
				255,
			},
			Lifetime: 30 + rand.Intn(30),
		}
		a.celebrationLights = append(a.celebrationLights, light)
	}

	// Mettre à jour les lumières existantes
	for i := 0; i < len(a.celebrationLights); i++ {
		a.celebrationLights[i].Lifetime--
		if a.celebrationLights[i].Lifetime <= 0 {
			// Supprimer la lumière
			a.celebrationLights = append(a.celebrationLights[:i], a.celebrationLights[i+1:]...)
			i--
		}
	}
}

// resetAnts réinitialise la position des fourmis
func (a *App) resetAnts() {
	startRoom := a.getStartRoom()
	for i := range a.ants {
		a.ants[i].X = float32(startRoom.X)
		a.ants[i].Y = float32(startRoom.Y)
		a.ants[i].Room = startRoom.Name
	}
	a.currentMove = -1
	a.animationProgress = 0
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func (a *App) getEndRoom() core.Room {
	for _, room := range a.gameState.Rooms {
		if room.IsEnd {
			return room
		}
	}
	return core.Room{}
}

func (a *App) getStartRoom() core.Room {
	for _, room := range a.gameState.Rooms {
		if room.IsStart {
			return room
		}
	}
	return core.Room{}
}

func calcSegmentPosition(startX, startY, length, angle float32) (endX, endY float32) {
	endX = startX + length*float32(math.Cos(float64(angle)))
	endY = startY + length*float32(math.Sin(float64(angle)))
	return
}

func (a *App) findRoomByName(name string) core.Room {
	for _, room := range a.gameState.Rooms {
		if room.Name == name {
			return room
		}
	}
	return core.Room{}
}
