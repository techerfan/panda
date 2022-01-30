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
	msgByte, err := msg.marshal()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(msgByte))

	messageInstance, err := unmarshalMsg(msgByte)
	if err != nil {
		t.Error(err)
	}

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

	msg, err := unmarshalMsg([]byte(msgStr))
	if err != nil {
		t.Error(err)
	}

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
