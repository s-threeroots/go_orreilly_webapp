package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		// ReadJSONは基本的にはUnmarshalと同じ。JSON側のキーが小文字でも対応する値が入る
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			/*
				if avatarURL, ok := c.userData["avatar_url"]; ok {
					msg.AvatarURL = avatarURL.(string)
				}
			*/
			c.room.forward <- msg
		} else {
			break
		}
	}
	defer c.socket.Close()
}

/*
func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	defer c.socket.Close()
}
*/

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	defer c.socket.Close()
}
