package card

import "math/rand"

type Card struct {
	Id      uint8
	Name    string
	Spell   string
	Quality uint8
}

func NewCard(id uint8) (*Card, error) {
	name, err := CardId(id).Name()
	if err != nil {
		return nil, err
	}
	spell, err := CardId(id).Spell()
	if err != nil {
		return nil, err
	}
	quality := uint8(rand.Intn(256))
	return &Card{
		Id:      id,
		Name:    name,
		Spell:   spell,
		Quality: quality,
	}, nil
}
