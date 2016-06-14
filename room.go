package main

import ( 
	"log" 
	"github.com/gorilla/websocket" 
	"net/http"
	"fmt"
	"github.com/stretchr/objx"
//	"time"
)

type room struct {
	forward chan *message
	join chan *client
	leave chan *client
	clients map[*client]bool
	avatar Avatar
}

func (r *room) run() {
	for {
		//time.Sleep(1000*time.Millisecond)
		select {
		case client := <- r.join:
			// 参加
			r.clients[client] = true
		case client := <- r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
		case msg := <- r.forward:
			// 受け取った message をすべての client へ転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// message を送信
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		//default:
		// default があるとすべて default で処理しようとして失敗する。
		//fmt.Print( "default select ok\n")
		}// select
	} // for
} // func (r *room) run() 

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ ReadBufferSize: socketBufferSize, 
	WriteBufferSize: socketBufferSize }

func (r *room) ServeHTTP( w http.ResponseWriter, req *http.Request ){
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err )
		return
	}
	fmt.Print("in room ServeHTTP not err\n")

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("くっきーの取得に失敗しました:", err )
	}
	client := &client{ 
		socket: socket,
		//send: make(chan []byte, messageBufferSize),
		send: make(chan *message, messageBufferSize),
		room: r, 
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client } ()
	go client.write()
	client.read()
}

func newRoom(avatar Avatar) *room {
//func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		avatar: avatar,
	}
}
