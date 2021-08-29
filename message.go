package panda

type messageStruct struct {
	channel string `json:"channel"`
	message []byte `json:"message"`
}

// type incomingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

// type forwardingMessage struct {
// 	channel string `json:"channel"`
// 	message []byte `json:"message"`
// }

func newMessage(channel string, message []byte) *messageStruct {

	msg := &messageStruct{
		channel: channel,
		message: message,
	}
	return msg
}
