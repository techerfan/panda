package panda

import (
	"encoding/json"
	"log"
)

type messageStruct struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

// type incomingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

// type forwardingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

func newMessage(channel string, message string) *messageStruct {

	msg := &messageStruct{
		Channel: channel,
		Message: message,
	}
	return msg
}

func (m *messageStruct) marshal() []byte {
	msgJSON, err := json.Marshal(&m)
	if err != nil {
		log.Println(err)
		return []byte("")
	}
	return msgJSON
}

func unmarshalMsg(msg []byte) *messageStruct {
	message := &messageStruct{}
	if err := json.Unmarshal(msg, message); err != nil {
		log.Println(err)
		return nil
	}
	return message
}
