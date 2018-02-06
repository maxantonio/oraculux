package main
////peers
////mining
////hasrate
////gasprice
////sincing
//package main
//
//import (
//	"flag"
//	"log"
//	"net/http"
//	"github.com/gorilla/websocket"
//	"github.com/onrik/ethrpc"
//	// "math/rand"
//	"strconv"
//	"time"
//	"github.com/Pallinder/go-randomdata"
//)
//
//type Hub struct {
//	// Registered clients.go
//	clients map[*Client]bool
//
//	// Register requests from the clients.
//	register chan *Client
//
//	// Unregister requests from clients.
//	unregister chan *Client
//}
//
//func NewHub() *Hub {
//	return &Hub{
//		register:   make(chan *Client),
//		unregister: make(chan *Client),
//		clients:    make(map[*Client]bool),
//	}
//}
//
//func (h *Hub) run() {
//	for {
//		select {
//		case client := <-h.register:
//			h.clients[client] = true
//		case client := <-h.register:
//			h.clients[client] = true
//		case client := <-h.unregister:
//			if _, ok := h.clients[client]; ok {
//				delete(h.clients, client)
//				close(client.send)
//			}
//		}
//	}
//}
//
//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool { return true },
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}
//
////Clase a mandar para React
//type SendClass struct {
//	Identificador  string `json:"symb"`
//	Best_Block int `json:"best"`
//	Uncles int `json:"uncles"`
//	Transactions int `json:"transactions"`
//	Uncle_count int `json:"uncle_count"`
//}
//// Client is a middleman between the websocket connection and the hub.
//type Client struct {
//	hub *Hub
//	// The websocket connection.
//	conn *websocket.Conn
//	// Buffered channel of outbound messages.
//	send chan SendClass
//}
//
//
//func (c *Client) writeInfo() {
//	ticker := time.NewTicker(5 * time.Second)
//
//	ethclient := ethrpc.New("http://127.0.0.1:8545")
//
//	defer func() {
//		ticker.Stop()
//		c.conn.Close()
//	}()
//	for {
//		select {
//		case <-ticker.C:
//			result, err2 := ethclient.EthSyncing()
//
//				log.Println("OBTENIDO RPC")
//				if err2 != nil {
//					log.Println("ERROR")
//					log.Fatal(err2)
//				}
//				log.Println(result)
//				log.Println("RESULTADO RPC")
//				//w, err := c.conn.NextWriter(websocket.TextMessage)
//				//if err != nil {
//				//	return
//				//}
//				//ESCRIBIENDO LA INFO A ENVIAR
//				c.conn.WriteJSON(result)
//				//w.Write(result)
//				//c.conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
//				//c.ws.WriteJSON(message)
//				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//					return
//				}
//			}
//		}
//	}
//
////Crea un nuevo objeto a mandar a las vista
//func (s *Emisora) GetFake() SendClass {
//	result, err2 := ethclient.EthSyncing()
//	return SendClass{
//		Identificador:s.identificador,
//		Best_Block: 1 + randomdata.Number(1, 90),
//		Uncles: 2 + randomdata.Number(1,70)/100,
//		Transactions: randomdata.Number(1, 80),
//		Uncle_count: 4 + randomdata.Number(1, 60)/100,
//	}
//}
//	// serveWs handles websocket requests from the peer.
//	func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//	log.Println(err)
//	return
//	}
//	client := &Client{hub: hub, conn: conn, send: make(chan SendClass)}
//	client.hub.register <- client
//		go client.writeInfo()
//	}
//
//	var addr = flag.String("addr", ":8080", "http service address")
//
//	func serveHome(w http.ResponseWriter, r *http.Request) {
//		log.Println(r.URL)
//		//if r.URL.Path != "/" {
//		//	http.Error(w, "Not found", 404)
//		//	return
//		//}
//
//		if r.Method != "GET" {
//			http.Error(w, "Method not allowed", 405)
//			return
//		}
//		http.ServeFile(w, r, "../../home.html")
//	}
//
//	func main() {
//	flag.Parse()
//	hub := NewHub()
//	go hub.run()
//	//iniciando servidor http
//	log.Println("INICIANDO SERVIDORES")
//	http.HandleFunc("/", serveHome)
//	//iniciando websocket
//	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//		serveWs(hub, w, r)
//	})
//	err := http.ListenAndServe(*addr, nil)
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//	log.Println("INICIANDO RPC")
//
//}
