package visualizer

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func (leg *Leg) Draw(screen *ebiten.Image) {
	// Calculer la fin de la cuisse
	thighEndX := leg.StartX + leg.LongCuisse*float32(math.Cos(float64(leg.ThighCuisseAngle)))
	thighEndY := leg.StartY + leg.LongCuisse*float32(math.Sin(float64(leg.ThighCuisseAngle)))

	// Calculer la fin du tibia (à partir de la fin de la cuisse)
	shinEndX := thighEndX + leg.LongTibia*float32(math.Cos(float64(leg.ThighCuisseAngle+leg.ShinTibiaAngle)))
	shinEndY := thighEndY + leg.LongTibia*float32(math.Sin(float64(leg.ThighCuisseAngle+leg.ShinTibiaAngle)))

	// Calculer la fin du pied (à partir de la fin du tibia)
	footEndX := shinEndX + leg.LongPied*float32(math.Cos(float64(leg.ThighCuisseAngle+leg.ShinTibiaAngle+leg.FootAngle)))
	footEndY := shinEndY + leg.LongPied*float32(math.Sin(float64(leg.ThighCuisseAngle+leg.ShinTibiaAngle+leg.FootAngle)))

	// Dessiner la cuisse
	vector.StrokeLine(screen, leg.StartX, leg.StartY, thighEndX, thighEndY, leg.StrokeWidth, leg.InnerColor, false)
	vector.StrokeLine(screen, leg.StartX, leg.StartY, thighEndX, thighEndY, leg.StrokeWidth/2, leg.OutlineColor, false)

	// Dessiner le tibia
	vector.StrokeLine(screen, thighEndX, thighEndY, shinEndX, shinEndY, leg.StrokeWidth, leg.InnerColor, false)
	vector.StrokeLine(screen, thighEndX, thighEndY, shinEndX, shinEndY, leg.StrokeWidth/2, leg.OutlineColor, false)

	// Dessiner le pied
	vector.StrokeLine(screen, shinEndX, shinEndY, footEndX, footEndY, leg.StrokeWidth, leg.InnerColor, false)
	vector.StrokeLine(screen, shinEndX, shinEndY, footEndX, footEndY, leg.StrokeWidth/2, leg.OutlineColor, false)
}

// drawRooms dessine les salles sur l'écran
func (a *App) drawRooms(screen *ebiten.Image) {
	for _, room := range a.gameState.Rooms {
		x := float32(room.X*a.gridSize + a.margin)
		y := float32(a.height - (room.Y*a.gridSize + a.margin))

		// Choisir la couleur en fonction du type de salle
		var roomColor color.RGBA
		if room.IsStart {
			roomColor = color.RGBA{34, 139, 34, 255} // Vert forêt // Vert pour la salle de départ
		} else if room.IsEnd {
			roomColor = color.RGBA{34, 139, 34, 255} // Rouge pour la salle d'arrivée
		} else {
			roomColor = color.RGBA{70, 70, 0, 255} // olive foncé pour les salles normales
		}

		// Dessiner le cercle de la salle
		vector.DrawFilledCircle(screen, x, y, float32(a.gridSize)/3, roomColor, false)

		// Dessiner le contour du cercle
		vector.StrokeCircle(screen, x, y, float32(a.gridSize)/3, 2, color.RGBA{128, 128, 0, 255}, false)

		// Dessiner le nom de la salle
		ebitenutil.DebugPrintAt(screen, room.Name, int(x)-len(room.Name)*3, int(y)-6)
	}
}

// drawConnections dessine les connexions entre les salles
func (a *App) drawConnections(screen *ebiten.Image) {
	drawnConnections := make(map[string]bool)
	for _, room := range a.gameState.Rooms {
		x1 := float32(room.X*a.gridSize + a.margin)
		y1 := float32(a.height - (room.Y*a.gridSize + a.margin))
		for _, linkName := range room.Links {
			linkedRoom, exists := a.findRoom(linkName)
			if exists {
				x2 := float32(linkedRoom.X*a.gridSize + a.margin)
				y2 := float32(a.height - (linkedRoom.Y*a.gridSize + a.margin))

				// Créer une clé unique pour chaque connexion
				connectionKey := fmt.Sprintf("%s-%s", room.Name, linkedRoom.Name)
				reverseKey := fmt.Sprintf("%s-%s", linkedRoom.Name, room.Name)

				// Vérifier si la connexion n'a pas déjà été dessinée
				if !drawnConnections[connectionKey] && !drawnConnections[reverseKey] {

					// Dessiner la ligne de connexion
					vector.StrokeLine(screen, x1, y1, x2, y2, 25, color.RGBA{70, 70, 0, 255}, false)
					drawnConnections[connectionKey] = true
				}
			}
		}
	}
}

// drawAnts dessine les fourmis sur l'écran
func (a *App) drawAnts(screen *ebiten.Image) {
	if len(a.ants) == 0 {
		ebitenutil.DebugPrint(screen, "Aucune fourmi à afficher")
		return
	}

	for _, ant := range a.ants {
		x, y := ant.X, ant.Y
		if len(a.gameState.Movements) > 0 {
			nextMoveIndex := (a.currentMove + 1) % len(a.gameState.Movements)
			nextMoves := a.gameState.Movements[nextMoveIndex]
			currentRoom := a.findRoomByName(ant.Room)
			targetRoom := currentRoom
			for _, move := range nextMoves {
				if move.AntID == ant.ID {
					targetRoom = a.findRoomByName(move.Room)
					break
				}
			}
			x = float32(currentRoom.X) + float32(targetRoom.X-currentRoom.X)*float32(a.animationProgress)
			y = float32(currentRoom.Y) + float32(targetRoom.Y-currentRoom.Y)*float32(a.animationProgress)
		}

		// Couleur de la fourmi
		antColor := ant.Color

		// Calcul des positions de base
		screenX := x*float32(a.gridSize) + float32(a.margin)
		screenY := float32(a.height) - (y*float32(a.gridSize) + float32(a.margin))

		// Paramètres pour la fourmi
		bodyRadius := float32(a.gridSize) / 9    //taille des cercles
		middleRadius := float32(a.gridSize) / 15 // Taille réduite du cercle central
		bodySpacing := bodyRadius * 1.2          // Espace entre les cercles du corps

		// Dessiner les quatre cercles du corps de la fourmi
		for i := 0; i < 4; i++ {
			radius := bodyRadius
			if i >= 1 && i <= 2 {
				radius = middleRadius
			}
			vector.DrawFilledCircle(screen, screenX-float32(i)*bodySpacing, screenY, radius, antColor, false)
			vector.StrokeCircle(screen, screenX-float32(i)*bodySpacing, screenY, radius, 1, color.RGBA{0, 0, 0, 255}, false)
		}

		////////////// PATTES

		// Définition des pattes avec des angles et positions spécifiques
		legParams := []Leg{
			{StartX: screenX - bodySpacing - 3, StartY: screenY, LongCuisse: 8, LongTibia: 6, LongPied: 4, ThighCuisseAngle: float32(math.Pi / 4), ShinTibiaAngle: float32(math.Pi / 3), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 2, OutlineColor: color.RGBA{R: 128, G: 128, B: 128, A: 255}, InnerColor: antColor}, // Patte avant gauche
			{StartX: screenX - bodySpacing - 3, StartY: screenY, LongCuisse: 6, LongTibia: 4, LongPied: 2, ThighCuisseAngle: float32(7 * math.Pi / 12), ShinTibiaAngle: float32(7 * math.Pi / 4), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 1, OutlineColor: color.RGBA{R: 169, G: 169, B: 169, A: 255}, InnerColor: antColor},
			{StartX: screenX - bodySpacing - 8, StartY: screenY, LongCuisse: 8, LongTibia: 6, LongPied: 4, ThighCuisseAngle: float32(math.Pi / 4), ShinTibiaAngle: float32(math.Pi / 3), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 2, OutlineColor: color.RGBA{R: 128, G: 128, B: 128, A: 255}, InnerColor: antColor}, // Patte avant gauche
			{StartX: screenX - bodySpacing - 8, StartY: screenY, LongCuisse: 6, LongTibia: 4, LongPied: 2, ThighCuisseAngle: float32(7 * math.Pi / 12), ShinTibiaAngle: float32(7 * math.Pi / 4), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 1, OutlineColor: color.RGBA{R: 169, G: 169, B: 169, A: 255}, InnerColor: antColor},
			{StartX: screenX - bodySpacing - 14, StartY: screenY, LongCuisse: 8, LongTibia: 6, LongPied: 4, ThighCuisseAngle: float32(math.Pi / 4), ShinTibiaAngle: float32(math.Pi / 3), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 2, OutlineColor: color.RGBA{R: 128, G: 128, B: 128, A: 255}, InnerColor: antColor}, // Patte avant gauche
			{StartX: screenX - bodySpacing - 14, StartY: screenY, LongCuisse: 6, LongTibia: 4, LongPied: 2, ThighCuisseAngle: float32(7 * math.Pi / 12), ShinTibiaAngle: float32(7 * math.Pi / 4), FootAngle: float32(57 * math.Pi / 36), StrokeWidth: 1, OutlineColor: color.RGBA{R: 169, G: 169, B: 169, A: 255}, InnerColor: antColor},
		}

		// Dessiner les pattes arrière en premier (2e, 4e, 6e) pour qu'elles apparaissent derrière le corps
		for i, leg := range legParams {
			if i%2 == 1 { // Index impair, donc pattes 2, 4, 6
				leg.Draw(screen)
			}
		}

		// Dessiner les quatre cercles du corps de la fourmi
		for i := 0; i < 4; i++ {
			radius := bodyRadius
			if i >= 1 && i <= 2 {
				radius = middleRadius
			}
			vector.DrawFilledCircle(screen, screenX-float32(i)*bodySpacing, screenY, radius, antColor, false)
			vector.StrokeCircle(screen, screenX-float32(i)*bodySpacing, screenY, radius, 1, color.RGBA{0, 0, 0, 255}, false)
		}

		// Dessiner les pattes avant (1ère, 3e, 5e) pour qu'elles apparaissent devant le corps
		for i, leg := range legParams {
			if i%2 == 0 { // Index pair, donc pattes 1, 3, 5
				leg.Draw(screen)
			}
		}

		//////////////  ANTENNE

		// Dessiner les antennes
		antennaTopLeftX := float32(6)      // Écartement horizontal de l'antenne en haut à gauche
		antennaBottomLeftX := float32(10)  // Écartement horizontal de l'antenne en bas à gauche
		antennaTopRightX := float32(-4)    // Écartement horizontal de l'antenne en haut à droite
		antennaBottomRightX := float32(-9) // Écartement horizontal de l'antenne en bas à droite
		antennaTopLeftY := float32(20)     // Longueur verticale de l'antenne en haut à gauche
		antennaBottomLeftY := float32(3)   // Longueur verticale de l'antenne en bas à gauche
		antennaTopRightY := float32(20)    // Longueur verticale de l'antenne en haut à droite
		antennaBottomRightY := float32(3)  // Longueur verticale de l'antenne en bas à droite
		antennaShiftLeft := float32(0)     // Décalage des antennes à gauche
		antennaShiftRight := float32(3)    // Décalage des antennes à droite
		antennaTipShiftX := float32(5)     // Décalage horizontal du petit trait en haut (peut être ajusté individuellement)
		antennaTipShiftY := float32(5)     // Décalage vertical du petit trait en haut (peut être ajusté individuellement)

		// Dessiner l'antenne gauche
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftLeft-antennaBottomLeftX, screenY+antennaBottomLeftY,
			screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX, screenY-antennaTopLeftY, 2, antColor, false)
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftLeft-antennaBottomLeftX, screenY+antennaBottomLeftY,
			screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX, screenY-antennaTopLeftY, 0, color.RGBA{0, 0, 0, 255}, false)

		// Dessiner l'antenne droite
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftRight+antennaBottomRightX, screenY+antennaBottomRightY,
			screenX+bodyRadius+antennaShiftRight+antennaTopRightX, screenY-antennaTopRightY, 2, antColor, false)
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftRight+antennaBottomRightX, screenY+antennaBottomRightY,
			screenX+bodyRadius+antennaShiftRight+antennaTopRightX, screenY-antennaTopRightY, 0, color.RGBA{0, 0, 0, 255}, false)

		// Dessiner le petit trait en haut de chaque antenne
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX, screenY-antennaTopLeftY,
			screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX-antennaTipShiftX, screenY-antennaTopLeftY+antennaTipShiftY, 2, antColor, false)
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX, screenY-antennaTopLeftY,
			screenX+bodyRadius+antennaShiftLeft-antennaTopLeftX-antennaTipShiftX, screenY-antennaTopLeftY+antennaTipShiftY, 0, color.RGBA{0, 0, 0, 255}, false)
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftRight+antennaTopRightX, screenY-antennaTopRightY,
			screenX+bodyRadius+antennaShiftRight+antennaTopRightX+antennaTipShiftX, screenY-antennaTopRightY+antennaTipShiftY, 2, antColor, false)
		vector.StrokeLine(screen, screenX+bodyRadius+antennaShiftRight+antennaTopRightX, screenY-antennaTopRightY,
			screenX+bodyRadius+antennaShiftRight+antennaTopRightX+antennaTipShiftX, screenY-antennaTopRightY+antennaTipShiftY, 0, color.RGBA{0, 0, 0, 255}, false)

		// Afficher l'ID de la fourmi
		antIDStr := fmt.Sprintf("%d", ant.ID)
		ebitenutil.DebugPrintAt(screen, antIDStr, (int(screenX)-len(antIDStr)*3)-23, int(screenY)-8)

		// Dessiner les yeux et la bouche
		antFace := ".."
		text.Draw(screen, antFace, basicfont.Face7x13, (int(screenX)-len(antFace)*7)+9, int(screenY), color.RGBA{255, 255, 255, 255})
		text.Draw(screen, antFace, basicfont.Face7x13, (int(screenX)-len(antFace)*7)+10, int(screenY), color.RGBA{0, 0, 0, 255})
		antBouche := "-"
		text.Draw(screen, antBouche, basicfont.Face7x13, (int(screenX)-len(antFace)*7)+12, int(screenY)+7, color.RGBA{255, 255, 255, 255})
		text.Draw(screen, antBouche, basicfont.Face7x13, (int(screenX)-len(antFace)*7)+13, int(screenY)+8, color.RGBA{0, 0, 0, 255})
	}
}

// drawCelebrationLights dessine les lumières de célébration
func (a *App) drawCelebrationLights(screen *ebiten.Image) {
	for _, light := range a.celebrationLights {
		vector.DrawFilledCircle(screen,
			float32(light.X)*float32(a.gridSize)+float32(a.margin),
			float32(a.height)-(float32(light.Y)*float32(a.gridSize)+float32(a.margin)),
			float32(light.Lifetime)/10,
			light.Color,
			false)
	}
}
