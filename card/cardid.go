package card

import "fmt"

// CardId はカードの ID を表す型。
// 0 はカード（ユニット）が存在しないことを表すために用いるため定義しない。
type CardId uint8

func (c CardId) Name() (string, error) {
	switch c {
	case 1:
		return "Slime", nil
	case 2:
		return "Dragon", nil
	case 3:
		return "Fairy", nil
	default:
		return "", fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Spell() (string, error) {
	switch c {
	case 1:
		return "Slippery, slippery. Slick with slime.", nil
	case 2:
		return "Don't dare disturb the dragon's deep slumber.", nil
	case 3:
		return "Fragile fairy, swift and bright, bring your healing, end this plight!", nil
	default:
		return "", fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Type() (CardType, error) {
	switch c {
	case 1:
		return Combat, nil
	case 2:
		return Combat, nil
	case 3:
		return Healing, nil
	default:
		return 0, fmt.Errorf("unknown card id: %d", c)
	}
}

func (c CardId) Symbol() (rune, error) {
	switch c {
	case 1:
		return 'S', nil
	case 2:
		return 'D', nil
	case 3:
		return 'F', nil
	default:
		return 0, fmt.Errorf("unknown card id: %d", c)
	}
}
