package services

import (
	"fmt"

	"github.com/gin-contrib/sessions"
)

type Message struct {
	Provider string
	Context  string
}

func GetMessageContainer(s sessions.Session) Message {
	message := Message{
		Provider: "",
		Context:  "",
	}
	flashes := s.Flashes()
	if len(flashes) == 1 {
		if flashes != nil {
			message = Message{
				Provider: "error",
				Context:  fmt.Sprint(flashes[0]),
			}
		}
	}

	return message
}
