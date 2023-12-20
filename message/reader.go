package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

type MessageReader struct {
	buf []byte
}

func NewMessageReader() *MessageReader {
	return &MessageReader{
		buf: make([]byte, 0, 1024),
	}
}

func (mr *MessageReader) Write(b []byte) {
	mr.buf = append(mr.buf, b...)
}

func (mr *MessageReader) Read() (Message, error) {
	msg, err := mr.tryRead()
	if err != nil {
		log.Println("Failed to parse message: ", err)
		return Message{}, err
	}
	return msg, nil
}

func (mr *MessageReader) tryRead() (Message, error) {
	br := bytes.NewBuffer(mr.buf) // Reader として利用

	// 1 バイト目が存在し、有効な msgType か
	b1, err := br.ReadByte()
	if err != nil {
		return Message{}, fmt.Errorf("invalid message format (required 1 msgType byte): %v", b1)
	}
	msgType, err := ParseMsgType(b1)
	if err != nil {
		return Message{}, err
	}
	switch msgType {
	case Cast:
		// 追加の 2 バイトが存在し、それが示す body の長さが足りるか
		bodyLenByte := br.Next(2)
		if len(bodyLenByte) < 2 {
			return Message{}, fmt.Errorf("invalid message format (Cast message must be more than 5 bytes): %v", bodyLenByte)
		}
		bodyLength := binary.BigEndian.Uint16(bodyLenByte)

		bodyByte := br.Next(int(bodyLength))
		if len(bodyByte) < int(bodyLength) {
			return Message{}, fmt.Errorf("invalid message format (lacking some bytes to fill bodyLength): %v", bodyByte)
		}
		body := string(bodyByte)

		// 読み込んだ分のバイト列を削除する
		// 1 + 2 + bodyLength バイトを読み込んだので、mr.buf にそれを詰め替えればいい
		mr.buf = br.Bytes()

		return Message{
			msgType: msgType,
			body:    body,
		}, nil
	case Surrender:
		// 読み込んだ分のバイト列を削除する
		// 1 バイトを読み込んだので、mr.buf にそれを詰め替えればいい
		mr.buf = br.Bytes()

		return NewSurrenderMessage(), nil
	default:
		return Message{}, fmt.Errorf("invalid message format (unknown msgType): %v", b1)
	}
}
