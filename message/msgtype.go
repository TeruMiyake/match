package message

import "fmt"

type MsgType uint8

const (
	Cast MsgType = iota
	Surrender
)

func ParseMsgType(b byte) (MsgType, error) {
	switch b {
	case 0:
		return Cast, nil
	case 1:
		return Surrender, nil
	default:
		return 0, fmt.Errorf("unknown message type: %d", b)
	}
}
