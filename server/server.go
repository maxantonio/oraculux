// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/onrik/ethrpc"
	// "math/rand"
	"strconv"
	"time"
)

type Hub struct {
	// Registered clients.go
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}


func (c *Client) writeInfo() {
	ticker := time.NewTicker(5 * time.Second)

	ethclient := ethrpc.New("http://127.0.0.1:8545")

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case <-ticker.C:
			result, err2 := ethclient.EthSyncing()
			log.Println("OBTENIDO RPC")
			if err2 != nil {
				log.Println("OBTENIDO RPC")
				if err2 != nil {
					log.Println("ERROR")
					log.Fatal(err2)
				}
				log.Println(result)
				log.Println("RESULTADO RPC")
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				//ESCRIBIENDO LA INFO A ENVIAR
				w.Write([]byte("Mayor: " + strconv.Itoa(result.HighestBlock) + " CurrentBlock: " + strconv.Itoa(result.CurrentBlock)))
				c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
	// serveWs handles websocket requests from the peer.
	func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
	log.Println(err)
	return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writeInfo()
	}

	var addr = flag.String("addr", ":80", "http service address")

	func serveHome(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		if r.URL.Path != "/" {
			http.Error(w, "Not found", 404)
			return
		}

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		http.ServeFile(w, r, "home.html")
	}

	func main() {
	flag.Parse()
	hub := NewHub()
	go hub.run()
	//iniciando servidor http
	log.Println("INICIANDO SERVIDORES")
	http.HandleFunc("/", serveHome)
	//iniciando websocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	log.Println("INICIANDO RPC")

}
