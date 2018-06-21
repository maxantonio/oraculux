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
	"flag"

)
import "../comon"
import s "strings"
var dummy = false

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//informacion que se envia al socket web preparada para enviar cualquier tipo de informacion
type SocketInfo struct {
	Server string `json:"server"`
	Info_type string `json:"info_type"`
	Data interface{} `json:"data"`
	Block int `json:"block"`
}

type ClientInfo struct {
	Action string            `json:"action"`
	Data   map[string]string `json:"data"`
}

//objeto de una conexion con las api
type Server struct{
	ws   *websocket.Conn
	send chan comon.ServerInfo
}

//para obtener info de una API

//para cuando un usuario se conecta
type Client struct {
	ws         *websocket.Conn
	sendServer chan SocketInfo
	subs       []string
}

//estructura encargada de gestionar los clientes web y las api de informacion
type Hub struct {
	clients      map[*Client]bool
	servers      map[*Server]bool
	broadcast    chan SocketInfo
	addClient    chan *Client
	removeClient chan *Client
	addServer    chan *Server
	removeServer chan *Server
	fullInfo     *comon.ServerInfo
}

//instancia unica del hub
var hub = Hub{
	broadcast:    make(chan SocketInfo),
	addClient:    make(chan *Client),
	addServer:    make(chan *Server),
	removeClient: make(chan *Client),
	removeServer: make(chan *Server),
	clients:      make(map[*Client]bool),
	servers:      make(map[*Server]bool),
	fullInfo: &comon.ServerInfo{
		Ping: "",
	},
}
//metodo encargado de obtener informacion del servidor local si esta activo
func (h *Hub) readSelfInfo() {
	ticker := time.NewTicker(3 * time.Second)
	rpc := ethrpc.New(*selfserver)
	fmt.Println("INICIALIZANDO TIMER HUB")
	for {
		select {
		case <-ticker.C:
				hashRate, err := rpc.EthHashrate()
				if (err != nil) {
					fmt.Println(err)
					continue
				}
				if hashRate != 0 {
					h.fullInfo.HashRate = hashRate
				}
				syncing, _ := rpc.EthSyncing()
				self_block := 0
				h.fullInfo.Sincing = syncing
				if (syncing.IsSyncing) {
					self_block = syncing.CurrentBlock
					//self_block = h.fullInfo.BlockNumber + 2;
					//fmt.Println(self_block)
				} else {
					self_block, err = rpc.EthBlockNumber()
					//self_block = h.fullInfo.BlockNumber + 2; //para uso local cuando no este online
					//fmt.Println(self_block)
				}
				if (self_block >= h.fullInfo.BlockNumber) {
					h.fullInfo.BlockNumber = self_block
					h.fullInfo.Uncles, _ = rpc.EthGetUncleCountByBlockNumber(self_block)
					h.fullInfo.Transactions, _ = rpc.EthGetBlockTransactionCountByNumber(self_block)
					h.fullInfo.Block, _ = rpc.EthGetBlockByNumber(self_block, false)
				}
				h.fullInfo.Peers, _ = rpc.NetPeerCount()
			gaspric, _ := rpc.EthGasPrice()
			h.fullInfo.GasPrice = gaspric.Int64()
				info_to_send := &SocketInfo{
					Info_type: "FullInfo",
					Data:      h.fullInfo,
					Server:    "selfi",
				}
				h.broadcast <- *info_to_send
			}
		}
	}


//administracion de los canales de comunicacion de los clientes y los servidores
func (h *Hub) start() {
	if (*mode == "merge") {
		fmt.Println("USANDO MODE MERGE")
		go hub.readSelfInfo()
	}
	for {
		select {
		case conn := <-h.addClient:
			h.clients[conn] = true
		case conn := <-h.addServer:
			h.servers[conn] = true
		case conn := <-h.removeServer:
			if _, ok := h.servers[conn]; ok {
				delete(h.servers, conn)
				//close(conn.read)  ver si es necesario cerrar el read del server
			}
		case conn := <-h.removeClient:
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				close(conn.sendServer)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.sendServer <- message:
				default:
					//close(client.sendServer)
					//delete(h.clients, client)
				}
			}
		}
	}
}

//encargado de enviar informacion a una instancia de cliente
func (c *Client) writeServers(){
	for {
		select {
		case message := <-c.sendServer:
			//fmt.Println("RECIBIDA INFO A ENVIAR")
			c.ws.WriteJSON(message)
		}
	}
}
func wsIndex(res http.ResponseWriter, req *http.Request) {
	conn, _ := upgrader.Upgrade(res, req, nil)

	client := &Client{
		ws:         conn,
		sendServer: make(chan SocketInfo),
		subs:       []string{},
	}
	fmt.Println("ABRIENDO WEBSOCKET")
	hub.addClient <- client
	go client.readClient()
	go client.writeServers() //mostrando info servidores conectados
}

//encargado de procesar las peticiones por socket del cliente
func (c *Client) readClient() {
	for {
		_, message, err := c.ws.ReadMessage()
		data := &ClientInfo{}
		fmt.Println("tratando de desparsear")
		err2 := json.Unmarshal(message, data)
		fmt.Println(data)
		fmt.Println("SALIO ?")
		fmt.Println(data.Data["account"])
		if err != nil || err2 != nil {
			fmt.Println("NOP NO SALIO, DIO ERROR")
			fmt.Println(err)
			fmt.Println(err2)
			fmt.Println(string(message))
			hub.removeClient <- c
			c.ws.Close()
			break
		} else {
			fmt.Println("NO DIO ERROR")
			rpc := ethrpc.New(*selfserver)
			balance, err3 := rpc.EthGetBalance(data.Data["account"], "latest")
			if err3 != nil {

			}
			fmt.Println(balance)
			fmt.Println(err3)
		}
	}
}
//envia la informacion resumida entre todos
func (h *Hub) sendFullInfo(data *comon.ServerInfo) {
	if (data.BlockNumber > hub.fullInfo.BlockNumber) {
		hub.fullInfo = data
	}
	Fullinfo := &SocketInfo{
		Info_type: "FullInfo",
		Data:      hub.fullInfo,
		Server:    "por definir",
	}
	hub.broadcast <- *Fullinfo
}
//proceso de lectura del socket de una instancia de servidor api
func (s *Server) read() {
	defer func() {
		hub.removeServer <- s
		s.ws.Close()
	}()
	for {
		_, message, err := s.ws.ReadMessage()
		data := &comon.ServerInfo{}
		fmt.Println("mensaje recibido de servidor")
		err2 := json.Unmarshal(message, data)

		if (err2 != nil) {
			fmt.Println("ERROR RECIBIENDO")
			fmt.Println(err2)
		} else {
			if data.Ping != "" {
				//DOING PONG
				fmt.Println("recibio un ping")
				s.ws.WriteJSON(data)
			} else {
				fmt.Println("RECIBIDA INFORMACION DE SERVIDOR")
				info := &SocketInfo{
					Info_type: "Server",
					Data:      data,
					Server:    "por definir",
				}
				go func() {
					hub.broadcast <- *info
					//hub.sendFullInfo(data)
				}()



			}
		}
		if err != nil {
			hub.removeServer <- s
			s.ws.Close()
			break
		}
	}
}

//handler de las peticiones web (renderizado de la pagina


//handler de las peticiones sockets de los servidores api
func wsApi(res http.ResponseWriter, req *http.Request){
	conn, _ := upgrader.Upgrade(res, req, nil)

	server := &Server{
		ws:   conn,
		send: make(chan comon.ServerInfo),
	}

	hub.addServer <- server
	fmt.Println("server conectado")
	go server.read()
}

//handler de las peticiones sockets de los clientes web
func wsApic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	urlfile := "../client/build" + r.URL.Path
	fmt.Println("procesando bien la url", urlfile)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, _ := json.Marshal(r.URL.Query().Get(":id"))
	fmt.Println(data)
	fmt.Println(r.URL.Path)
	w.Write(data)
	return
}
//handler de las peticiones sockets de los clientes web
func serveHome(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
	if s.Contains(r.URL.Path, "apic") {
		fmt.Println("revisado por home pero obviado")
		wsApic(w, r)
		return
	}
		urlfile := "../client/build"+r.URL.Path
		fmt.Println("retorna fichero",urlfile)

	http.ServeFile(w, r, urlfile)
	}

// fariables que se pueden recibir por parametros
var mode = flag.String("mode", "merge", "modo del servidor(self,soloapi,merge[default]")
//
var selfserver = flag.String("selfserver", "http://127.0.0.1:8545", "direccion servidor rpc de la info propia del stat")

//func main() {
//	flag.Parse()
//
//	go hub.start()
//	//manejando las peticiones por websockets de los servidores API (LISTADO)
//
//	//manejando laspeticiones web http
//
//	//manejando las peticiones por websockets de los clientes web
//	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
//		wsIndex(w, r)
//	})
//
//	//manejando las peticiones por websockets de los servidores API (LISTADO)
//	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
//		wsApi(w, r)
//	})
//	http.HandleFunc("/apic/{method}/{id}", func(w http.ResponseWriter, r *http.Request) {
//		wsApic(w, r)
//	})
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		serveHome(w, r)
//	})
//	//iniciando el servidor por el puerto
//	http.ListenAndServe(":8080", nil)
//}