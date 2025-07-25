package cluster

import (
	"net"
	"strconv"
	"time"
	"encoding/json"
	"log"
)

type Seed struct {
	seedIpAddr *net.TCPAddr
}

func NewSeed(host string, port int) *Seed {
	conv := strconv.Itoa(port)
	ip, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, conv))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &Seed{
		seedIpAddr: ip,
	}
}

func (s *Seed) PerformJoinRequest(myAddress string) {
	conn, err := net.DialTimeout("tcp", s.seedIpAddr.String(), 10 * time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	b, e := json.Marshal(map[string]string{
		"endpoint": "/join",
		"ip_address": myAddress,
	})

	if e != nil {
		return
	}

	_, errW := conn.Write(b)
	if errW != nil {
		return
	}

	buf := make([]byte, 1<<24)
	n, errR := conn.Read(buf)
	if errR != nil {
		return
	}
	log.Printf("Ack: %s\n", string(buf[:n]))
}
