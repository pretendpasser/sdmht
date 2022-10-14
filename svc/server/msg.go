package server

import (
	"bytes"
	"encoding/binary"
)

var (
	ByteOrder = binary.LittleEndian
)

type Message struct {
	key   int32
	value string
}

func (m *Message) Marshal() ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, ByteOrder, m.key); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, ByteOrder, m.value); err != nil {
		return nil, err
	}

	b := buf.Bytes()

	return b, nil
}
