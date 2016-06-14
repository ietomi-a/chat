package main
import ( 
	"github.com/gorilla/websocket" 
	"fmt"
)

type client struct {
	socket *websocket.Conn
	send chan []byte
	room *room
}

func (c *client) read() {
	for { 
		//fmt.Print("in client read\n")
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			//fmt.Print("client read\n")
			c.room.forward <- msg
		}else{
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send { 
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
