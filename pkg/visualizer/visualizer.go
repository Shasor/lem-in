package visualizer

import (
	"fmt"
	"io"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// App représente l'application principale du visualiseur
type App struct {
	width, height     int
	gameState         *GameState
	gridSize          int
	margin            int
	maxX, maxY        int
	ants              []Ant
	currentMove       int
	animationProgress float64
	animationSpeed    float64
	isPaused          bool
	antsMoved         int
	antsInMotion      int
	waitForNextMove   bool
	earthImage        *ebiten.Image
	grassImage        *ebiten.Image
	skyImage          *ebiten.Image
	isCelebrating     bool
	celebrationLights []CelebrationLight
}

// NewApp crée une nouvelle instance de l'application
func NewApp(width, height int, getReader func() io.Reader) (*App, error) {
	gameState, err := ParseInput(getReader())
	if err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de l'entrée: %v", err)
	}

	maxX, maxY := 0, 0
	for _, room := range gameState.Rooms {
		if room.X > maxX {
			maxX = room.X
		}
		if room.Y > maxY {
			maxY = room.Y
		}
	}

	app := &App{
		width:             width,
		height:            height,
		gameState:         gameState,
		gridSize:          60,
		margin:            200,
		maxX:              maxX,
		maxY:              maxY,
		currentMove:       -1,
		animationProgress: 0,
		animationSpeed:    0.002,
		antsMoved:         0,
		antsInMotion:      0,
		waitForNextMove:   false,
	}

	windowWidth := (maxX+1)*app.gridSize + 2*app.margin
	windowHeight := (maxY+1)*app.gridSize + 2*app.margin
	app.width = windowWidth
	app.height = windowHeight

	app.generateEarthImage()
	app.generateGrassImage()
	app.generateSkyImage()
	app.createAnts()

	return app, nil
}

// Update met à jour l'état de l'application
func (a *App) Update() error {
	if len(a.gameState.Movements) == 0 {
		return nil
	}

	// Ajout de débogage pour voir l'état actuel
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		a.animationSpeed *= 1.1 // Augmente la vitesse de 10%
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		a.animationSpeed *= 0.9 // Réduit la vitesse de 10%
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		a.isPaused = !a.isPaused
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		if a.isCelebrating {
			a.isCelebrating = false
			a.resetAnts()
		} else if a.allAntsAtEnd() {
			a.isCelebrating = true
			a.startCelebration()
		} else {
			a.waitForNextMove = false
		}
	}

	if !a.isPaused && !a.waitForNextMove {
		a.animationProgress += a.animationSpeed
		if a.animationProgress >= 1 {
			a.animationProgress = 0
			a.currentMove++
			if a.currentMove >= len(a.gameState.Movements) {
				a.currentMove = 0 // Recommencer l'animation test
			}

			if a.isCelebrating {
				a.updateCelebration()
			}

			// Mettre à jour la position des fourmis
			currentMoves := a.gameState.Movements[a.currentMove]
			for _, move := range currentMoves {
				for i := range a.ants {
					if a.ants[i].ID == move.AntID {
						targetRoom, exists := a.findRoom(move.Room)
						if exists {
							oldX, oldY := a.ants[i].X, a.ants[i].Y
							a.ants[i].Room = move.Room
							a.ants[i].X = float32(targetRoom.X)
							a.ants[i].Y = float32(targetRoom.Y)

							// Calculer la direction
							dx := a.ants[i].X - oldX
							dy := a.ants[i].Y - oldY
							a.ants[i].Direction = float32(math.Atan2(float64(dy), float64(dx)))
						}
					}
				}
			}

			// Réinitialiser les compteurs de mouvement
			a.antsMoved = 0
			a.antsInMotion = len(a.gameState.Movements[a.currentMove])

			// Mettre à jour la position des fourmis
			currentMoves = a.gameState.Movements[a.currentMove]
			for _, move := range currentMoves {
				for i := range a.ants {
					if a.ants[i].ID == move.AntID {

						// Mettre à jour la position des fourmis
						targetRoom, exists := a.findRoom(move.Room)
						if exists {
							a.ants[i].Room = move.Room
							a.ants[i].X = float32(targetRoom.X)
							a.ants[i].Y = float32(targetRoom.Y)
							a.antsMoved++
						}
					}
				}

				// Mettre à jour la position des fourmis en fonction de l'animation
				for i := range a.ants {
					dx := a.ants[i].TargetX - a.ants[i].X
					dy := a.ants[i].TargetY - a.ants[i].Y

					// Calculer la nouvelle position
					newX := a.ants[i].X + dx*float32(a.animationProgress)
					newY := a.ants[i].Y + dy*float32(a.animationProgress)

					// Calculer la distance parcourue
					distance := math.Sqrt(float64((newX-a.ants[i].X)*(newX-a.ants[i].X) + (newY-a.ants[i].Y)*(newY-a.ants[i].Y)))

					// Mettre à jour la position
					a.ants[i].X = newX
					a.ants[i].Y = newY

					// Calculer la direction si la fourmi se déplace
					if dx != 0 || dy != 0 {
						a.ants[i].Direction = float32(math.Atan2(float64(dy), float64(dx)))

						// Ajuster la vitesse des pattes en fonction de la distance parcourue
						a.ants[i].LegSpeed = distance * 10 // Ajustez ce facteur pour changer la vitesse de l'animation

						// Mettre à jour la phase des pattes
						a.ants[i].LegPhase += a.ants[i].LegSpeed
						if a.ants[i].LegPhase > math.Pi*2 {
							a.ants[i].LegPhase -= math.Pi * 2
						}
					} else {
						// Ralentir le mouvement des pattes si la fourmi ne se déplace pas
						a.ants[i].LegSpeed *= 0.9
						if a.ants[i].LegSpeed < 0.01 {
							a.ants[i].LegSpeed = 0
						}
						a.ants[i].LegPhase += a.ants[i].LegSpeed
					}
				}
			}

			// Vérifier si toutes les fourmis ont bougé
			if a.antsMoved >= a.antsInMotion {
				a.waitForNextMove = true // Met l'animation en pause jusqu'au prochain mouvement
			}
		}
	}
	return nil
}

// Draw dessine l'état actuel de l'application sur l'écran
func (a *App) Draw(screen *ebiten.Image) {
	screen.DrawImage(a.skyImage, nil)
	earthHeight := int(float64(a.height) * 0.80)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(0, float64(a.height)*0.05) // 10% du haut
	screen.DrawImage(a.grassImage, op)

	// Dessiner la terre (fourmilière, en bas)
	op = &ebiten.DrawImageOptions{}

	op.GeoM.Translate(0, float64(a.height)-float64(earthHeight))
	screen.DrawImage(a.earthImage, op)

	// Ajuster la zone de dessin pour la fourmilière
	a.margin = a.height - earthHeight

	a.drawConnections(screen)
	a.drawRooms(screen)
	a.drawAnts(screen)
	if a.isCelebrating {
		a.drawCelebrationLights(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Nombre de fourmis: %d", a.gameState.NumAnts))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouvements parsés: %d", len(a.gameState.Movements)), 0, 20)
	if len(a.gameState.Movements) > 0 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouvement actuel: %d/%d", a.currentMove+1, len(a.gameState.Movements)), 0, 40)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Progression: %.2f", a.animationProgress), 0, 60)
	}
}

// Layout définit la taille de l'écran
func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return a.width, a.height
}

// Run lance l'application
func (a *App) Run() error {
	ebiten.SetWindowSize(a.width, a.height)
	ebiten.SetWindowTitle("Lem-in Visualizer")
	ebiten.SetTPS(30)
	return ebiten.RunGame(a)
}
