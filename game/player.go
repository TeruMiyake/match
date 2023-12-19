package game

import (
	"github.com/TeruMiyake/match/card"
)

type Player struct {
	Id    uint8
	Name  string
	Cards []card.Card
}

var noPlayer = Player{
	Id:    0,
	Name:  "No player",
	Cards: []card.Card{},
}
