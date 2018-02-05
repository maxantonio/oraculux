package main

import (
	"net/http"
	"fmt"
	"time"
	"strings"
	"github.com/gorilla/websocket"
	"github.com/Pallinder/go-randomdata"
)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Clase a mandar para React
type SendClass struct {
	Identificador  string `json:"symb"`
	Best_Block int `json:"best"`
	Uncles int `json:"uncles"`
	Transactions int `json:"transactions"`
	Uncle_count int `json:"uncle_count"`
}
//para cuando un usuario se conecta
type Client struct {
	ws   *websocket.Conn
	send chan SendClass
	subs []string
} 
type Emisora struct{
	identificador string
	clients int
	sockets map[*Client]bool
	addClient    chan *Client
	removeClient chan *Client
}
//inicia a mandar objetos, se adiciona y se eliminan cliente usando los "channels"
func (em *Emisora) start() {
	//cada 2 segundo se repite el for
	invervalo := time.Tick(2 *time.Second)
	for {
		select {
		case conn := <- em.addClient:
			em.sockets[conn] = true
		case conn := <- em.removeClient:
			if _, ok := em.sockets[conn]; ok {
				delete(em.sockets, conn)
			}
		case <-invervalo:
			temporal := em.GetFake()
		 	
		 	for cnn := range em.sockets {		
				cnn.send <- temporal
			}
			
		}
	}
	
}

//Crea un nuevo objeto a mandar a las vista
func (s *Emisora) GetFake() SendClass {
    return SendClass{
		Identificador:s.identificador,
    	Best_Block: 1 + randomdata.Number(1, 90),
    	Uncles: 2 + randomdata.Number(1,70)/100,
    	Transactions: randomdata.Number(1, 80),
    	Uncle_count: 4 + randomdata.Number(1, 60)/100,
    }
}

var FirstValues = map[string]Emisora{  
	"eth1": { "eth1",  0,  make(map[*Client]bool), make(chan *Client), make(chan *Client)},
	"eth2": { "eth2", 0,  make(map[*Client]bool), make(chan *Client), make(chan *Client)},
	"eth3": { "eth3",  0,  make(map[*Client]bool), make(chan *Client), make(chan *Client)},
	"eth4": { "eth4", 0, make(map[*Client]bool), make(chan *Client), make(chan *Client)},
	"eth5": { "eth5",  0,  make(map[*Client]bool), make(chan *Client), make(chan *Client)},
}

type Hub struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	addClient    chan *Client
	removeClient chan *Client
}

var hub = Hub{
	broadcast:    make(chan []byte),
	addClient:    make(chan *Client),
	removeClient: make(chan *Client),
	clients:      make(map[*Client]bool),
}

func (hub *Hub) start() {
	for {
		select {
		case conn := <-hub.addClient:
			hub.clients[conn] = true
		case conn := <-hub.removeClient:
			if _, ok := hub.clients[conn]; ok {
				delete(hub.clients, conn)
				close(conn.send)
			}
		}
	}
}

func (c *Client) write() {
	defer func() {
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.ws.WriteJSON(message)

		}
	}
}

func (c *Client) read() {
	defer func() {
		hub.removeClient <- c
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		str :=  string(message[:])
		if !strings.Contains(str, "cancelar_") {
			emisoraObj := FirstValues[str]
			emisoraObj.addClient <- c
			fmt.Println(str)
			c.subs = append(c.subs, str)
		}else {
			emisoraObj :=  FirstValues[strings.SplitAfter(str, "_")[1]]
			emisoraObj.removeClient <- c
			fmt.Println(str)
		}
		if err != nil {
			hub.removeClient <- c
			for i := 0; i < len(c.subs); i++ {
			  FirstValues[c.subs[i]].removeClient <- c
			}
			c.ws.Close()
			break
		}
	}
}

func wsIndex(res http.ResponseWriter, req *http.Request) {
	conn, _ := upgrader.Upgrade(res, req, nil)

	client := &Client{
		ws:   conn,
		send: make(chan SendClass),
		subs: []string{},
	}

	hub.addClient <- client

	go client.write()
	go client.read()
}

func main() {
	go hub.start()
	for v  := range FirstValues{
		w := FirstValues[v]
		go w.start()
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		wsIndex(w, r)
	})
	http.ListenAndServe(":6060",nil)
	
		
}

