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
	"net"
	"errors"
	"strconv"
	"encoding/json"
)
type Server struct{
	socket        *websocket.Conn
	ServerInfo    *ServerInfo
	rpc           *ethrpc.EthRPC
	last_block    int
	last_peers    int
	pendingFilter string
	eth_coinbase  string
	pongCh        chan struct{}
}
type PoolStatus struct {
	pending int
	queued  int
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
		s.rpc.NetListening()
		s.ServerInfo.Block,s.ServerInfo.Err = s.rpc.EthGetBlockByNumber(s.ServerInfo.BlockNumber,false)
		s.ServerInfo.Peers,s.ServerInfo.Err = s.rpc.NetPeerCount()
		s.ServerInfo.IsMining, s.ServerInfo.Err = s.rpc.EthMining()
		s.ServerInfo.Pending, s.ServerInfo.Err = s.rpc.EthGetTransactionCount(s.eth_coinbase, "pending")
		fmt.Println("respuesta de txpool")
		response1, err1 := s.rpc.Call("txpool_status")
		data := &PoolStatus{}
		json.Unmarshal(response1, data)
		fmt.Println(response1)
		fmt.Println(data)
		fmt.Println(err1)

		s.ServerInfo.Transactions, s.ServerInfo.Err = s.rpc.EthGetBlockTransactionCountByNumber(s.ServerInfo.BlockNumber)

	}
	fmt.Println("Info  pedido de envio")
	if (s.ServerInfo.BlockNumber > s.last_block || s.ServerInfo.Peers != s.last_peers) {
		s.socket.WriteJSON(s.ServerInfo)
		fmt.Print(s.ServerInfo.BlockNumber)
		fmt.Println("Info  Enviada")
	}
	s.last_block = s.ServerInfo.BlockNumber
	s.last_peers = s.ServerInfo.Peers
}


func (server *Server) start() {
	go server.read()
	ticker := time.NewTicker(4 * time.Second)

	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			fmt.Print(t.String())
			server.reportLatency()
			server.write()
		}
	}
}
func (s *Server) read() {
	for {
		_, _, err := s.socket.ReadMessage()
		if err == nil {
			select {
			case s.pongCh <- struct{}{}:
				// Pong delivered, continue listening
				continue
			default:
				// Ping routine dead, abort
				fmt.Println("Stats server se murio")
				return
			}
		}
	}
}

func (s *Server) reportLatency() error {
	// Send the current time to the ethstats server
	start := time.Now()

	infoping := *s.ServerInfo
	infoping.Ping = start.String()
	if err := s.socket.WriteJSON(infoping); err != nil {
		fmt.Println(err)
		return err
	}
	// Wait for the pong request to arrive back
	select {
	case <-s.pongCh:
		// Pong delivered, report the latency
	case <-time.After(5 * time.Second):
		// Ping timeout, abort
		return errors.New("ping timed out")
	}
	latency := strconv.Itoa(int((time.Since(start) / time.Duration(2)).Nanoseconds() / 100000))
	s.ServerInfo.Ping = ""
	s.ServerInfo.Latency = latency
	return s.socket.WriteJSON(s.ServerInfo)
}

var addr = flag.String("addr", "35.227.84.238:80", "http service address")
var name = flag.String("name", "", "nombre a mostrar del servidor")
var rpc = flag.String("rpc", "http://localhost:8545", "nombre a mostrar del servidor")
func main() {

	flag.Parse()
	log.SetFlags(3)

	fmt.Println(*addr)
	fmt.Println(*name)
	fmt.Println(*rpc)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/api"}
	log.Printf("connecting to %s", u.String())
	ethclient := ethrpc.New(*rpc) //conectando al rpc local
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	fmt.Println(err)
	defer c.Close()
	if *name == "" {
		ifaces, err := net.Interfaces()
		fmt.Println(err)
		if err == nil {
			var ip net.IP
			for _, i := range ifaces {
				addrs, _ := i.Addrs()
				// handle err
				for _, addr := range addrs {

					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
				}
			}
			*name = ip.String()
		}
	}
	serverInfo := &ServerInfo{
		Server: *name,
		Ping:   "",
	}
	base, _ := ethclient.EthCoinbase()
	server := &Server{
		socket:       c,
		ServerInfo:   serverInfo,
		rpc:          ethclient,
		last_block:   0,
		last_peers:   0,
		eth_coinbase: base,
		pongCh:       make(chan struct{}),
	}

	server.pendingFilter, _ = server.rpc.EthNewPendingTransactionFilter()
	server.start()
}