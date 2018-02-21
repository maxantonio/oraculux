package main

import (
	"fmt"
	"time"
	"github.com/gorilla/websocket"
	"github.com/onrik/ethrpc"
	//"github.com/Pallinder/go-randomdata"
	//"github.com/Pallinder/go-randomdata"
	"flag"
	"log"
	"os"
	"net/url"
	"os/signal"
)
type Server struct{
	socket *websocket.Conn
	ServerInfo *ServerInfo
	rpc *ethrpc.EthRPC
}
type ServerInfo struct{
	Sincing *ethrpc.Syncing
	Block *ethrpc.Block
	BlockNumber int
	Peers int
	IsMining bool
	Err error
}

func (s *Server) write() {

	s.ServerInfo.Sincing,s.ServerInfo.Err = s.rpc.EthSyncing()
	s.ServerInfo.BlockNumber,s.ServerInfo.Err = s.rpc.EthBlockNumber()
	if(s.ServerInfo.Err!=nil){

	}else{
		if(s.ServerInfo.Sincing.IsSyncing){
			if(s.ServerInfo.Sincing.CurrentBlock>s.ServerInfo.BlockNumber){
				s.ServerInfo.BlockNumber = s.ServerInfo.Sincing.CurrentBlock
			}
		}
		s.ServerInfo.Block,s.ServerInfo.Err = s.rpc.EthGetBlockByNumber(s.ServerInfo.BlockNumber,false)
		s.ServerInfo.Peers,s.ServerInfo.Err = s.rpc.NetPeerCount()
	}
	fmt.Println("Info  pedido de envio")
	s.socket.WriteJSON(s.ServerInfo)
	fmt.Print(s.ServerInfo.BlockNumber)
	fmt.Println("Info  Enviada")
}

var addr = flag.String("addr", "localhost:80", "http service address")


func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api"}
	log.Printf("connecting to %s", u.String())
	ethclient := ethrpc.New("http://127.0.0.1:8545")//conectando al rpc local
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	serverInfo := &ServerInfo{}
	server := &Server{
		socket:   c,
		ServerInfo: serverInfo,
		rpc:ethclient,
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			fmt.Print(t.String())
			go server.write()
		}
	}
}