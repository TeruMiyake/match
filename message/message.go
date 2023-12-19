package message

type Message struct {
	// メッセージの種類
	msgType MsgType
	// メッセージの内容
	content string
}

func (msg Message) Bytes() []byte {
	buf := make([]byte, 0, 1024)
	buf = append(buf, byte(msg.msgType))

	contentBytes := []byte(msg.content)

	buf = append(buf, byte(len(contentBytes)))
	buf = append(buf, contentBytes...)

	return buf
}

func (msg Message) String() string {
	s := ""
	switch msg.msgType {
	case Cast:
		s = "Cast: " + msg.content
	case Surrender:
		s = "Surrender"
	default:
		panic("Unknown message type")
	}
	return s
}
