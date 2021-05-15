package main

import (
	"log"
	"net/http"
	"oreilly/web_app/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// join
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			// leave
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case msg := <-r.forward:
			// send message to all client
			r.tracer.Trace("メッセージを受信しました:", string(msg.Message))
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace(" -- クライアントに送信されました")
					// send message
				default:
					// fault sending
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗しました。クライアントをクリーンナップします")
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		println(err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("クッキーの取得に失敗しました:", err)
	}

	client := &client{
		socket: socket,
		send:   make(chan *message, messageBufferSize),
		room:   r,
		// ここでクッキー詰めてる
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
