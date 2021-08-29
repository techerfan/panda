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
