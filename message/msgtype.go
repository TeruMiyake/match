package message

type MsgType uint8

const (
	Cast MsgType = iota
	Surrender
)
