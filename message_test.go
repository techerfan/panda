package panda

import (
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	msg := &messageStruct{
		Channel: "sth",
		Message: "My message",
	}

	fmt.Println(string(msg.marshal()))

	messageInstance := unmarshalMsg(msg.marshal())

	if msg.Channel != messageInstance.Channel || msg.Message != messageInstance.Message {
		t.Fail()
	}
}

func TestUnmarshal(t *testing.T) {
	msgStr := `{
		"msgType": 0,
		"channel": "chat",
		"message": "This is a test message."
	}`

	msgType := 0
	m := "This is a test message."
	ch := "chat"

	msg := unmarshalMsg([]byte(msgStr))

	if msg.Message != m {
		t.Error("Message did not match")
	}

	if msg.Channel != ch {
		t.Error("Channel did not match")
	}

	if msg.MsgType != MessageType(msgType) {
		t.Error("Message type did not match")
	}
}
