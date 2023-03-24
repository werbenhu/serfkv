package cluster

import "encoding/json"

type Message struct {
	Action string
	Key    string
	Val    any
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Decode(payload []byte) error {
	return json.Unmarshal(payload, m)
}
