package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	ServerMessagePrefix = ""
	HelpMessage         = `
list of available commands:

	:help                   display this help message
	:info                   display client's name and joined rooms
	:name <name>            set client's name
	:join <room>            join a room
	:leave <room>           leave a room
	:exit                   disconnect
	to:<room> <message>     send a message to a room
`
)

type Client struct {
	conn  net.Conn
	name  string
	rooms map[string]struct{}
}

func NewClient(conn net.Conn) *Client {
	newClient := &Client{
		conn:  conn,
		name:  conn.RemoteAddr().String(),
		rooms: make(map[string]struct{}),
	}

	newClient.rooms[""] = struct{}{}

	return newClient
}

func (c *Client) loop(s *Server) {
	reader := bufio.NewReader(c.conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error in reading from client: %s\n", c.conn.RemoteAddr())
			return
		}

		str = strings.TrimRight(str, "\r\n")

		switch {
		case strings.HasPrefix(str, ":exit"):
			return

		case strings.HasPrefix(str, ":help"):
			c.writeHelpMessage()

		case strings.HasPrefix(str, ":name "):
			c.name = str[6:]

		case strings.HasPrefix(str, ":info"):
			c.writeClientInfo()

		case strings.HasPrefix(str, ":join "):
			room := str[6:]
			if len(room) >= 1 {
				c.writeText(fmt.Sprintf("you have joined room %q", room))
				c.rooms[room] = struct{}{}
			}

		case strings.HasPrefix(str, ":leave "):
			room := str[6:]
			if len(room) >= 1 {
				// TODO fix this bug
				delete(c.rooms, room)
				c.writeText(fmt.Sprintf("you have left room %q", room))
			}

		case strings.HasPrefix(str, "to:"):
			words := strings.Split(str, " ")
			room := strings.Split(words[0], ":")[1]
			if len(words) >= 2 {
				s.broadcast <- NewMessage(c, strings.Join(words[1:], " "), room)
			}

		default:
			s.broadcast <- NewMessage(c, str, "")
		}
	}
}

func (c *Client) writeText(msg string) {
	_, err := c.conn.Write([]byte(ServerMessagePrefix + msg + "\n"))
	if err != nil {
		// TODO figure out what to do here
		log.Println(err)
	}
}

func (c *Client) writeMessage(message *Message) {
	n, err := c.conn.Write([]byte(message.content))
	if err != nil {
		// TODO figure out what to do here
		log.Println(err)
	}
	log.Printf("Wrote %d bytes\n", n)
}

func (c *Client) writeClientInfo() {
	// TODO clean this stuff:
	// use bufio writer instead of manual writing

	c.conn.Write([]byte(("============ Client info ==============\n")))
	c.conn.Write([]byte((fmt.Sprintf("Name: %q\n", c.name))))
	c.conn.Write([]byte(("Rooms:\n")))
	for room := range c.rooms {
		if len(room) == 0 {
			continue
		}
		c.conn.Write([]byte((room + "\n")))
	}
	c.conn.Write([]byte(("============ Client info ==============\n")))
}

func (c *Client) writeHelpMessage() {
	// TODO clean this stuff:
	// use bufio writer instead of manual writing

	c.conn.Write([]byte(("============ Help message ==============\n")))
	c.writeText(HelpMessage)
	c.conn.Write([]byte(("============ Help message ==============\n")))
}
