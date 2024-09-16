package core

type Room struct {
	Name           string
	X, Y           int
	Ants           []Ant
	Links          []string
	IsStart, IsEnd bool
}

type Ant struct {
	Index int
	Path  []string
}
