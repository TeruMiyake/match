package card

import "fmt"

// CardId はカードの ID を表す型。
// 0 はカード（ユニット）が存在しないことを表すために用いるため定義しない。
type CardId uint8

func (c CardId) Name() (string, error) {
	switch c {
	case 1:
		return "Slime", nil
	default:
		return "", fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Spell() (string, error) {
	switch c {
	case 1:
		return "Slippery, slippery. Slick with slime.", nil
	default:
		return "", fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Type() (CardType, error) {
	switch c {
	case 1:
		return Combat, nil
	default:
		return 0, fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Symbol() (rune, error) {
	switch c {
	case 1:
		return 'S', nil
	default:
		return 0, fmt.Errorf("unknown card id: %d", c)
	}
}
