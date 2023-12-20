package message

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	// メッセージの種類
	msgType MsgType
	// メッセージの内容
	// utf-8 でエンコードしてバイト長が uint16 に収まる必要がある
	body string
}

func NewCastMessage(body string) Message {
	return Message{
		msgType: Cast,
		body:    body,
	}
}

func NewSurrenderMessage() Message {
	return Message{
		msgType: Surrender,
	}
}

// Bytes はメッセージをバイト列に変換する。
// バイト列のフォーマットは以下の通り。
// MsgType: uint8 (1 byte)
// bodyLength: uint32 (4 bytes)
// body: []byte (bodyLength bytes of utf-8 encoded string)
func (msg Message) Bytes() []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteByte(byte(msg.msgType))

	bodyBytes := []byte(msg.body)

	bodyLen := uint16(len(bodyBytes))

	binary.Write(buf, binary.BigEndian, bodyLen)
	buf.Write(bodyBytes)

	return buf.Bytes()
}

func (msg Message) String() string {
	s := ""
	switch msg.msgType {
	case Cast:
		s = "Cast: " + msg.body
	case Surrender:
		s = "Surrender"
	default:
		panic("Unknown message type")
	}
	return s
}
