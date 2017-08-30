package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client는 한명의 채팅 사용자
type client struct {
	// socket은 이 클라이언트의 웹 소켓
	socket *websocket.Conn
	// send는 메시지가 전송되는 채널
	send chan *message
	room *room
	// userData는 사용자 대한 정보를 보유
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarURL = c.room.avatar.GetAvatarURL(c)
			msg.AvatarURL = avatarURL.(string)
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}