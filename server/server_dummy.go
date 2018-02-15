package main

import (
	"net/http"
	"fmt"
	"time"
	"strings"
	"github.com/gorilla/websocket"
	"github.com/onrik/ethrpc"
	//"github.com/Pallinder/go-randomdata"
	"github.com/Pallinder/go-randomdata"
)
var dummy = false
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
type SocketInfo struct {
	Info_type string `json:"info_type"`
	Data interface{} `json:"data"`
	Block int `json:"block"`
}
//para cuando un usuario se conecta
type Client struct {
	ws   *websocket.Conn
	send chan SocketInfo
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
			go func() { c.send <- s.GetUncles(rpc,number) }()
			go func() { c.send <- s.GetTransactionCount(rpc,number) }()
			return sock
		}
		sock = SocketInfo{
			Info_type:"Syncing",
			Data:result,
		}
		fmt.Print(sock)
		fmt.Println("sinck pidioendo uncles")
		go func() { c.send <- s.GetUncles(rpc,result.CurrentBlock) }()
		go func() { c.send <- s.GetTransactionCount(rpc,result.CurrentBlock) }()
	}



	fmt.Println("sinck terminado")
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
	fmt.Println("fake generado")
	fmt.Print(sock)
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
	fmt.Println("fake generado")
	fmt.Print(sock)
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
	fmt.Println("fake generado")
	fmt.Print(sock)
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
	fmt.Println("fake generado")
	fmt.Print(sock)
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
	fmt.Println("fake generado")
	fmt.Print(sock)
	return sock
}
//Crea un nuevo objeto a mandar a las vista
func (s *Emisora) GetFake() SocketInfo {
	result := SendClass{
		Identificador:s.identificador,
		Best_Block: 1 + randomdata.Number(1, 90),
		Uncles: 2 + randomdata.Number(1,70)/100,
		Transactions: randomdata.Number(1, 80),
		Uncle_count: 4 + randomdata.Number(1, 60)/100,
	}
	sock := SocketInfo{
		Info_type:"Sendclass",
		Data:result,
	}
	fmt.Println("fake generado")
	fmt.Print(sock)
    return sock
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
	ticker := time.NewTicker(time.Second)
	emisora := FirstValues["eth1"]
	fmt.Println("iniciado write")
	ethclient := ethrpc.New("http://127.0.0.1:8545")
	for {
		select {
		case <-ticker.C:
			fmt.Println("Generando facke")
			if (dummy) {
				go func() { c.send <- emisora.GetFake() }()
				fmt.Println("Fake pedido")
			} else {
				go func() { c.send <- emisora.GetEthHashrate(ethclient) }()
				go func() { c.send <- emisora.GetSyncing(ethclient, c) }()
				go func() { c.send <- emisora.GetPeers(ethclient) }()
				go func() { c.send <- emisora.EthGasPrice(ethclient) }()
				fmt.Println("Info  pedido")
			}
		case message := <-c.send:
			fmt.Println("Escribiendo mensaje")
			fmt.Println(message)
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

func wsIndex(res http.ResponseWriter, req *http.Request){
	conn, _ := upgrader.Upgrade(res, req, nil)

	client := &Client{
		ws:   conn,
		send: make(chan SocketInfo),
		subs: []string{},
	}

	hub.addClient <- client
	fmt.Println("cliente recibido")
	go client.write()
	go client.read()
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("manejado por serveHome")
	fmt.Println(r.URL.Path)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	//if r.URL.Path == "/" {
	//	fmt.Println("retorna html index")
	//	http.ServeFile(w, r, "../client/build/index.html")
	//	return
	//}else{
		urlfile := "../client/build"+r.URL.Path
		fmt.Println("retorna fichero",urlfile)
		http.ServeFile(w, r, urlfile)
	//	return
	//}
	//http.ServeFile(w, r, "../client/build/index.html")
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
	http.ListenAndServe(":8080",nil)
}

