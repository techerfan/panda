package panda

import (
	"encoding/json"

	"github.com/techerfan/panda/logger"
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

// type incomingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

// type forwardingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

func newMessage(channel string, message string, msgType MessageType) *messageStruct {

	msg := &messageStruct{
		MsgType: msgType,
		Channel: channel,
		Message: message,
	}
	return msg
}

func (m *messageStruct) marshal() []byte {
	msgJSON, err := json.Marshal(&m)
	if err != nil {
		logger.GetLogger().Log(logger.Error, err.Error())
		return []byte("")
	}
	return msgJSON
}

func unmarshalMsg(msg []byte) *messageStruct {
	message := &messageStruct{}
	if err := json.Unmarshal(msg, message); err != nil {
		logger.GetLogger().Log(logger.Error, err.Error())
		return nil
	}
	return message
}
