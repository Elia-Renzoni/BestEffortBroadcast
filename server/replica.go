package server

import (
	"beb/cluster"
	"beb/dialer"
	"beb/store"
)

import (
	"net"
	"strconv"
	"time"
	"log"
	"io"
	"encoding/json"
)

const minPortAvaliable = 0x00
const maxPortAvaliable = 0xFFFF

type ReplicaStmt interface {
	NewReplica() *Replica
	StartListen()
}

type Replica struct {
	ipAddr *net.TCPAddr
	ls net.Listener
	socketError error
	socketHandling chan <- struct{}
	deadlineSetter func(net.Conn)
	cache store.Storage
}

func NewReplica(host string, port int, socketChan chan <- struct{}) *Replica {
	if (host == "") && (port < minPortAvaliable || port > maxPortAvaliable) {
		return nil
	}

	castedPort := strconv.Itoa(port)

	ip, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, castedPort))
	if err != nil {
		return nil
	}
	return &Replica{
		ipAddr: ip,
		socketHandling: socketChan,
		deadlineSetter: func(conn net.Conn) {
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		},
		cache: store.NewInMemoryKV(),
	}
}

func (r *Replica) StartListen() {
	r.ls, r.socketError = net.Listen("tcp", r.ipAddr.String())	
	if r.socketError != nil {
		r.socketHandling <- struct{}{}
		log.Fatal(r.socketError.Error())
		return
	}

	log.Println("Server Ready...")
	for {
		conn, err := r.ls.Accept()
		r.deadlineSetter(conn)
		if err != nil {
			r.socketHandling <- struct{}{}	
			log.Fatal(r.socketError.Error())
			return
		}
		go r.handleConnection(conn)
	}
}

func (r *Replica) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {conn.Close()}(conn)
	 
	buffer := make([]byte, 1<<24)
	var (
		n int
		msg dialer.Message
		err error
	)
	
	n, err = conn.Read(buffer)
	netErr, ok := err.(net.Error)
	if ok {		
		if netErr.Timeout() {
			log.Fatal("Timeout Error!")
			return
		}
	}

	json.Unmarshal(buffer[:n], &msg)
	r.printIncomingMessage(msg)
	val := r.messageRouter(msg)
	encodedRes, err := json.Marshal(map[string]string{
		"ack": val,
	})

	if err != nil {
		r.sendBack(conn, []byte(err.Error()))	
	} else {
		r.sendBack(conn, encodedRes)
	}
}

func (r *Replica) messageRouter(msg dialer.Message) string {
	switch msg.Endpoint {
	case "/join":
		process := r.createNewProcess(msg.ProcessIPAddr)
		process.AddProcess()
	case "/set":
		r.cache.Set(msg.Key, msg.Value)

		// performing a broadcast action using Best Effort Broadcast
		// abstraction
		msg.Broadcast(cluster.Process{}.GetProcessGroup())
	case "/get":
		return string(r.cache.Get(msg.Key))
	case "/delete":
		r.cache.Delete(msg.Key)

		// broadcast message to all the correct processes
		msg.Broadcast(cluster.Process{}.GetProcessGroup())
	}

	return "ok"
}

func (r *Replica) createNewProcess(address string) cluster.ClusterManager {
	ipAddr, err := net.ResolveIPAddr("tcp", address)
	if err != nil {
		return cluster.Process{}
	}
	return cluster.Process{
		IpAddr: ipAddr,
		Conn: net.Dialer{
			KeepAlive: time.Second * 10,
			Timeout: time.Second * 10,
		},
		Protocol: "tcp",
	}
}

func (r *Replica) sendBack(w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (r *Replica) printIncomingMessage(msg dialer.Message) {
	log.Println(msg.Endpoint)
	log.Println(msg.Key)
	log.Println(msg.Value)
	log.Println(msg.ProcessIPAddr)
}
