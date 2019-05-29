package message_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"testing"
)

func TestIsValid(t *testing.T) {
	msg := message.Message{
		Title: message.Ready,
	}
	msg.IsValid()
	msg.Title = message.GoTo
	msg.IsValid()
	msg.Title = message.ClientAnswer
	msg.IsValid()
	msg.Title = message.Leave
	msg.IsValid()
	msg.Title = message.Continue
	msg.IsValid()
	msg.Title = message.ChangeOpponent
	msg.IsValid()
	msg.Title = message.Quit
	msg.IsValid()
	msg.Title = message.State
	msg.IsValid()
	msg.Title = message.NotDesiredPack
	msg.IsValid()
	msg.Title = "ggggg"
	msg.IsValid()
}
