package game

type GameStatus uint8

const (
	Waiting GameStatus = iota
	Playing
	Finished
)
