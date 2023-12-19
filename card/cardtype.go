package card

type CardType uint8

const (
	Combat CardType = iota
	Healing
	FieldEffect
	Object
)
