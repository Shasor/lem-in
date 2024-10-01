package visualizer

import (
	"image/color"
	"lem-in/pkg/core"
)

// GameState représente l'état du jeu
type GameState struct {
	NumAnts   int
	Rooms     []core.Room
	Movements [][]AntMove
}

// Ant représente une fourmi dans le jeu
type Ant struct {
	ID               int
	X, Y             float32
	TargetX, TargetY float32
	Room             string
	Color            color.RGBA
	Direction        float32
	LaunchAngle      float64
	LaunchVelocity   float64
	LaunchTime       float64
	IsLaunched       bool
	LegPhase         float64
	LegSpeed         float64
}

// AntMove représente un mouvement de fourmi
type AntMove struct {
	AntID int
	Room  string
}

// Structure pour les pattes
type Leg struct {
	StartX, StartY   float32     // Point de départ (où la patte est attachée au corps)
	LongCuisse       float32     // Longueur de la cuisse
	LongTibia        float32     // Longueur du tibia
	LongPied         float32     // Longueur du pied
	ThighCuisseAngle float32     // Angle de la cuisse (en radians)
	ShinTibiaAngle   float32     // Angle du tibia (en radians)
	FootAngle        float32     // Angle du pied (en radians)
	StrokeWidth      float32     // Largeur du trait
	OutlineColor     color.Color // Couleur du contour
	InnerColor       color.Color // Couleur interne de la patte
}

// CelebrationLight représente une lumière de célébration
type CelebrationLight struct {
	X, Y     float32
	Color    color.RGBA
	Lifetime int
}
