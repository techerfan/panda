package panda

import (
	"encoding/json"
)

type MessageType int

const (
	Raw MessageType = iota
	Subscribe
	Unsubscribe
)

type messageStruct struct {
	MsgType MessageType `json:"msgType"`
	Channel string      `json:"channel"`
	Message string      `json:"message"`
}

func newMessage(channel string, message string, msgType MessageType) *messageStruct {

	msg := &messageStruct{
		MsgType: msgType,
		Channel: channel,
		Message: message,
	}
	return msg
}

func (m *messageStruct) marshal() ([]byte, error) {
	msgJSON, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}
	return msgJSON, nil
}

func unmarshalMsg(msg []byte) (*messageStruct, error) {
	message := &messageStruct{}
	if err := json.Unmarshal(msg, message); err != nil {
		return nil, err
	}
	return message, nil
}
