package game

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"math/rand"

	"github.com/TeruMiyake/match/card"
)

type Game struct {
	status  GameStatus
	players []Player
	width   uint8
	height  uint8
	cells   [][]Cell
}

func NewGame(width, height uint8) *Game {
	// 処理の都合上のみの理由で int へ変換
	// 無駄な型変換を書かずに済ますためだけであり、値の変更はなし
	w, h := int(width), int(height)

	cells := make([][]Cell, height)
	for i := 0; i < h; i++ {
		cells[i] = make([]Cell, width)
		for j := 0; j < w; j++ {
			cells[i][j] = Cell{
				ownerId: 0,
				unit:    nil,
			}
		}
	}

	g := &Game{
		status:  Waiting,
		players: []Player{},
		width:   width,
		height:  height,
		cells:   cells,
	}

	// セル内に 3 個までランダムなユニットを配置する
	for i := 0; i < 3; i++ {
		x := rand.Intn(w)
		y := rand.Intn(h)
		u := NewUnit(1, 1, 255)
		g.deployUnit(x, y, u)
	}

	return g
}

func (g Game) Illustrate() {
	for _, row := range g.cells {
		for _, cell := range row {
			if cell.unit == nil {
				print(".")
			} else {
				print(cell.unit.SymbolStr())
			}
		}
		println()
	}
}

func (g *Game) deployUnit(x int, y int, unit *Unit) {
	g.cells[y][x].unit = unit
}

// Bytes はゲーム情報をバイト列に変換する
// 先頭にはゲームデータのサイズ (uint32 = 4 bytes) が付与される
// サイズの値は自身のサイズ (4 bytes) を含まない
func (g Game) Bytes() []byte {
	buf := make([]byte, 0, 1024)

	// status: 1 byte = uint8
	buf = append(buf, byte(g.status))

	// playersNum: 1 byte = uint8
	buf = append(buf, byte(len(g.players)))
	// []Player
	for _, player := range g.players {
		// playerID: 1 byte = uint8
		buf = append(buf, byte(player.Id))
		// playerNameLen: 1 byte = uint8
		buf = append(buf, byte(len(player.Name)))
		// playerName: n byte = []byte
		buf = append(buf, []byte(player.Name)...)
		// playerCardsNum: 1 byte = uint8
		buf = append(buf, byte(len(player.Cards)))
		// []Card
		for _, card := range player.Cards {
			// cardID: 1 byte = uint8
			buf = append(buf, byte(card.Id))
			// cardQuality: 1 byte = uint8
			buf = append(buf, byte(card.Quality))
			// spellLen: 1 byte = uint8
			buf = append(buf, byte(len(card.Spell)))
			// spell: n byte = []byte
			buf = append(buf, []byte(card.Spell)...)
		}
	}

	// width: 1 byte = uint8
	buf = append(buf, byte(g.width))
	// height: 1 byte = uint8
	buf = append(buf, byte(g.height))

	// cells
	for _, row := range g.cells {
		for _, cell := range row {
			if cell.unit == nil {
				// nil: seen as cardId 0
				// cardId: 1 byte = uint8
				buf = append(buf, 0)
				// cardId 0 has no additional bytes
			} else {
				// cardId: 1 byte = uint8
				buf = append(buf, byte(cell.unit.cardId))
				// ownerID: 1 byte = uint8
				buf = append(buf, byte(cell.ownerId))
				// health: 4 bytes = float32
				hu32 := math.Float32bits(cell.unit.health)
				hbu32 := make([]byte, 4)
				binary.BigEndian.PutUint32(hbu32, hu32)
				buf = append(buf, []byte(hbu32)...)
			}
		}
	}

	// ゲームデータのサイズを先頭に追加する
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(buf)))

	ret := make([]byte, 0, len(size)+len(buf))
	ret = append(ret, size...)
	ret = append(ret, buf...)

	return ret
}

// FromBytes はバイト列からゲーム情報を復元する
func FromBytes(b []byte) *Game {
	log.Println("Re-constructing game data from bytes... bytes: ", b)
	g := &Game{}
	buf := bytes.NewBuffer(b)

	// status: 1 byte = uint8
	g.status = GameStatus(buf.Next(1)[0])

	// playersNum: 1 byte = uint8
	playersNum := int(buf.Next(1)[0])
	g.players = make([]Player, playersNum)

	// []Player
	for j := 0; j < int(playersNum); j++ {
		// playerID: 1 byte = uint8
		g.players[j].Id = buf.Next(1)[0]
		// playerNameLen: 1 byte = uint8
		nameLen := buf.Next(1)[0]
		// playerName: n byte = []byte
		g.players[j].Name = string(buf.Next(int(nameLen)))
		// playerCardsNum: 1 byte = uint8
		cardsNum := buf.Next(1)[0]
		g.players[j].Cards = make([]card.Card, cardsNum)
		// []Card
		for k := 0; k < int(cardsNum); k++ {
			// cardID: 1 byte = uint8
			g.players[j].Cards[k].Id = buf.Next(1)[0]
			// cardQuality: 1 byte = uint8
			g.players[j].Cards[k].Quality = buf.Next(1)[0]
			// spellLen: 1 byte = uint8
			spellLen := int(buf.Next(1)[0])
			// spell: n byte = []byte
			g.players[j].Cards[k].Spell = string(buf.Next(spellLen))
		}
	}

	// width: 1 byte = uint8
	g.width = buf.Next(1)[0]
	// height: 1 byte = uint8
	g.height = buf.Next(1)[0]

	// cells
	g.cells = make([][]Cell, g.height)
	for j := 0; j < int(g.height); j++ {
		g.cells[j] = make([]Cell, g.width)
		for k := 0; k < int(g.width); k++ {
			// cardId: 1 byte = uint8
			cardId := buf.Next(1)[0]
			if cardId == 0 {
				// nil: seen as cardId 0
				// cardId 0 has no additional bytes
				g.cells[j][k].unit = nil
			} else {
				// ownerID: 1 byte = uint8
				ownerId := buf.Next(1)[0]
				// health: 4 bytes = float32
				hu32 := binary.BigEndian.Uint32(buf.Next(4))
				g.cells[j][k].unit = &Unit{cardId: card.CardId(cardId), health: math.Float32frombits(hu32)}
				g.cells[j][k].ownerId = ownerId
			}
		}
	}

	return g
}

func (g *Game) AddPlayer() uint8 {
	id := uint8(len(g.players) + 1)
	g.players = append(g.players, Player{
		Id:    id,
		Name:  "Player " + string(id),
		Cards: []card.Card{},
	})
	return id
}
