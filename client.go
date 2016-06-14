package main
import ( 
	"github.com/gorilla/websocket" 
	"time"
//	"fmt"
)

type client struct {
	socket *websocket.Conn
	//send chan []byte
	send chan *message
	room *room
	userData map[string]interface{}
}

func (c *client) read() {
	for { 
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
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
		//if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
