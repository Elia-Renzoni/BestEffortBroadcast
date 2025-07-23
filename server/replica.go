package server

import (
	"net"
	"strconv"
)

const minPortAvaliable = 0x00
const maxPortAvaliable = 0xFFFF

type ReplicaStmt interface {
	StartListen()
	Shutdown()
}

type Replica struct {
	ipAddr *net.TCPAddr
	ls net.Listener
	socketError error
	socketHandling chan struct{}
}

func NewReplica(host string, port int) *Replica {
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
		socketHandling: make(chan struct{}),
	}
}

func (r *Replica) StartListen() {
	r.ls, r.socketError = net.Listen("tcp", r.ipAddr.String())	

	for {
		conn, err := r.ls.Accept()
		if err != nil {
			r.socketHandling <- struct{}{}	
		}
		go r.handleConnection(conn)
	}
}

func (r *Replica) handleConnection(conn net.Conn) {
	// TODO
}

func (r *Replica) Shutdown() {
	// TODO
}

