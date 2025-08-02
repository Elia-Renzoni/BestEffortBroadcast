package dialer

import (
	"beb/cluster"
	"log"
	"encoding/json"
	"net"
	"io"
)

type Message struct {
	Endpoint string `json:"endpoint"`
	Key string `json:"key"`
	Value []byte `json:"value,omitempty"`
	ProcessIPAddr string `json:"ip_address,omitempty"`
	Broadcaster bool `json:"broadcasted,omitempty"`
}

func (m Message) Broadcast(pgroup []cluster.Process) {
	for _, val := range pgroup {	
		func(process cluster.Process) {
			log.Println(process.TcpAddr.String())
			conn, err := process.Conn.Dial(process.Protocol, process.TcpAddr.String())
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			defer conn.Close()
			
			jsonMessage, errMarshal := json.Marshal(m)
			if errMarshal != nil {
				log.Fatal(errMarshal.Error())
				return
			}
			conn.Write(jsonMessage)

			// read reply
			buf := make([]byte, 1<<24)
			var n int
			var errRead error
			for {
				n, errRead = conn.Read(buf)
				if errRead != nil {
					netErr, ok := errRead.(net.Error)
					if ok {
						if netErr.Timeout() {
							log.Fatal("Timeout Occured, Connection Aborted")
							return
						}
					}
					
					if errRead == io.EOF {
						break
					}
				}	
			}

			log.Printf("Reply: %s", string(buf[:n]))
		}(val)
	}
}
