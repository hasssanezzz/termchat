package main

import "fmt"

type Message struct {
	client  *Client
	content string
	room    string
}

func NewMessage(client *Client, msg string, room string) *Message {
	message := &Message{
		client:  client,
		content: "",
		room:    room,
	}

	if len(room) == 0 {
		message.content = fmt.Sprintf("[%s] >> %s\n", client.name, msg)
	} else {
		message.content = fmt.Sprintf("[%s]:[%s] >> %s\n", client.name, room, msg)
	}

	return message
}
