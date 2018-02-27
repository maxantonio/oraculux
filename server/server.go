package main

import (
	"net/http"
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"github.com/onrik/ethrpc"
	//"github.com/Pallinder/go-randomdata"
	//"github.com/Pallinder/go-randomdata"
	"encoding/json"
)
var dummy = false
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
const (
	// Time allowed to write a message to the peer.
	writeWait =  time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 2 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)
//Clase a mandar para React
type SendClass struct {
	Identificador  string `json:"symb"`
	Best_Block int `json:"best"`
	Uncles int `json:"uncles"`
	Transactions int `json:"transactions"`
	Uncle_count int `json:"uncle_count"`
}

type SocketInfo struct {
	Server string `json:"server"`
	Info_type string `json:"info_type"`
	Data interface{} `json:"data"`
	Block int `json:"block"`
}
//para obtener info de una API
type Server struct{
	ws   *websocket.Conn
	send chan ServerInfo
}
type ServerInfo struct{
	Server       string
	Sincing      *ethrpc.Syncing
	Block        *ethrpc.Block
	BlockNumber  int
	Peers        int
	IsMining     bool
	Transactions int
	Pending      int
	Ping         string
	Latency      string
	Err          error
}

//para cuando un usuario se conecta
type Client struct {
	ws         *websocket.Conn
	send       chan SocketInfo
	sendServer chan SocketInfo
	subs       []string
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
			//temporal := em.GetFake()

			//for cnn := range em.sockets {
			//	//cnn.send <- temporal
			//}

		}
	}

}

func (s *Emisora) GetSyncing(rpc *ethrpc.EthRPC,c *Client) SocketInfo {
	result,error := rpc.EthSyncing()
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		if result.CurrentBlock==0{
			number,err2 := rpc.EthBlockNumber()
			if err2 !=nil {
				sock = SocketInfo{
					Info_type:"Error",
					Data:error,
				}
			}
			sock = SocketInfo{
				Info_type:"FullBlock",
				Data:number,
			}
			go func() { c.send <- s.GetBlockByNumber(rpc,number)}()
			go func() { c.send <- s.GetUncles(rpc,number) }()
			go func() { c.send <- s.GetTransactionCount(rpc,number) }()
			return sock
		}
		sock = SocketInfo{
			Info_type:"Syncing",
			Data:result,
		}
		go func() { c.send <- s.GetUncles(rpc,result.CurrentBlock) }()
		go func() { c.send <- s.GetTransactionCount(rpc,result.CurrentBlock) }()
		go func() { c.send <- s.GetBlockByNumber(rpc,result.CurrentBlock) }()
	}
	return sock
}

func (s *Emisora) GetEthHashrate(rpc *ethrpc.EthRPC) SocketInfo {
	result,error := rpc.EthHashrate()
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"Hashrate",
			Data:result,
		}
	}
	return sock
}
func (s *Emisora) GetTransactionCount(rpc *ethrpc.EthRPC,currentBlock int) SocketInfo {
	result,error := rpc.EthGetBlockTransactionCountByNumber(currentBlock)
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"Transactions",
			Data:result,
			Block:currentBlock,
		}
	}
	return sock
}
func (s *Emisora) GetBlockByNumber(rpc *ethrpc.EthRPC,currentBlock int) SocketInfo {
	result,error := rpc.EthGetBlockByNumber(currentBlock,false)
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"Block",
			Data:result,
			Block:currentBlock,
		}
	}
	return sock
}
func (s *Emisora) GetUncles(rpc *ethrpc.EthRPC,currentBlock int) SocketInfo {
	result,error := rpc.EthGetUncleCountByBlockNumber(currentBlock)
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"Uncles",
			Data:result,
			Block:currentBlock,
		}
	}
	return sock
}
func (s *Emisora) EthGasPrice(rpc *ethrpc.EthRPC) SocketInfo {
	result,error := rpc.EthGasPrice()
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"GasPrice",
			Data:result.String(),
		}
	}

	return sock
}
func (s *Emisora) GetPeers(rpc *ethrpc.EthRPC) SocketInfo {
	result,error := rpc.NetPeerCount()
	var sock SocketInfo
	if error != nil {
		sock = SocketInfo{
			Info_type:"Error",
			Data:error,
		}
	}else{
		sock = SocketInfo{
			Info_type:"Peers",
			Data:result,
		}
	}

	return sock
}

var FirstValues = map[string]Emisora{
	"eth1": { "eth1",  0,  make(map[*Client]bool), make(chan *Client), make(chan *Client)},
}

type Hub struct {
	clients      map[*Client]bool
	servers		map[*Server]bool
	broadcast    chan SocketInfo
	addClient    chan *Client
	removeClient chan *Client
	addServer	 chan *Server
	removeServer chan *Server
}

var hub = Hub{
	broadcast:    make(chan SocketInfo),
	addClient:    make(chan *Client),
	addServer:    make(chan *Server),
	removeClient: make(chan *Client),
	removeServer: make(chan *Server),
	clients:      make(map[*Client]bool),
	servers:      make(map[*Server]bool),
}

func (hub *Hub) start() {
	for {
		select {
		case conn := <-hub.addClient:
			hub.clients[conn] = true
		case conn := <-hub.addServer:
			hub.servers[conn] = true
		case conn := <-hub.removeServer:
			if _, ok := hub.servers[conn]; ok {
				delete(hub.servers, conn)
				//close(conn.read)  ver si es necesario cerrar el read del server
			}
		case conn := <-hub.removeClient:
			if _, ok := hub.clients[conn]; ok {
				delete(hub.clients, conn)
				close(conn.send)
			}

		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.sendServer <- message:
				default:

					close(client.sendServer)
					delete(hub.clients, client)
				}
			}
		}
	}
}

func (c *Client) write() {
	defer func() {
		c.ws.Close()
	}()
	ticker := time.NewTicker(2 * time.Second)
	emisora := FirstValues["eth1"]
	ethclient := ethrpc.New("http://127.0.0.1:8545")
	for {
		select {
		case <-ticker.C:
			if (dummy) {
				//go func() { c.send <- emisora.GetFake() }()
				fmt.Println("Fake pedido")
			} else {
				go func() { c.send <- emisora.GetEthHashrate(ethclient) }()
				go func() { c.send <- emisora.GetSyncing(ethclient, c) }()
				go func() { c.send <- emisora.GetPeers(ethclient) }()
				go func() { c.send <- emisora.EthGasPrice(ethclient) }()
				fmt.Println("Info  pedido")
			}
		case message := <-c.send:
			c.ws.WriteJSON(message)
		}
	}
}

func (c *Client) writeServers(){
	for {
		select {
		case message := <-c.sendServer:
			c.ws.WriteJSON(message)
		}
	}
}

func (s *Server) read() {
	defer func() {
		hub.removeServer <- s
		s.ws.Close()
	}()
	for {
		_, message, err := s.ws.ReadMessage()
		data := &ServerInfo{}
		fmt.Println(string(message[:]))
		err2 := json.Unmarshal(message, data)
		if (err2 != nil) {
			fmt.Println(err2)
		} else {
			if data.Ping != "" {
				//DOING PONG
				fmt.Println("recibio un ping")
				fmt.Println(data)
				s.ws.WriteJSON(data)
			} else {
				info := &SocketInfo{
					Info_type: "Server",
					Data:      data,
					Server:    "por definir",
				}
				hub.broadcast <- *info
			}
		}

		if err != nil {
			hub.removeServer <- s
			s.ws.Close()
			break
		}
	}
}

func wsIndex(res http.ResponseWriter, req *http.Request){
	conn, _ := upgrader.Upgrade(res, req, nil)

	client := &Client{
		ws:         conn,
		send:       make(chan SocketInfo),
		sendServer: make(chan SocketInfo),
		subs:       []string{},
	}

	hub.addClient <- client
	go client.write() //mostrando info servidor local
	go client.writeServers() //mostrando info servidores conectados
}
func wsApi(res http.ResponseWriter, req *http.Request){
	conn, _ := upgrader.Upgrade(res, req, nil)

	server := &Server{
		ws:   conn,
		send: make(chan ServerInfo),
	}

	hub.addServer <- server
	fmt.Println("server conectado")
	go server.read()
}
func serveHome(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}

		urlfile := "../client/build"+r.URL.Path
		fmt.Println("retorna fichero",urlfile)
		http.ServeFile(w, r, urlfile)
	}
func main() {
	go hub.start()
	//for v  := range FirstValues{
	//	w := FirstValues[v]
	//	go w.start()
	//}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		serveHome(w,r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		wsIndex(w, r)
	})
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request){
		wsApi(w, r)
	})
	http.ListenAndServe(":80",nil)
	//este si
}

