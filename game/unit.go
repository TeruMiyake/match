package game

import (
	"github.com/TeruMiyake/match/card"
)

type Unit struct {
	cardId  card.CardId
	ownerId uint8
	quality uint8
	health  float32
}

func NewUnit(ownerId uint8, cardId card.CardId, quality uint8) *Unit {
	return &Unit{
		ownerId: ownerId,
		cardId:  cardId,
		quality: quality,
		health:  1.0,
	}
}

func (u Unit) SymbolStr() string {
	r, err := card.CardId(u.cardId).Symbol()
	if err != nil {
		panic(err)
	}
	return string(r)
}
