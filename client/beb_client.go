package main

import (
	"net"
	"log"
	"time"
	"encoding/json"
	"os"
	"io"
)

type Message struct {
	Endpoint string `json:"endpoint"`
	Key string `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}

func main() {
	log.Println("Client Ready...")

	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "7878"), 10 * time.Second)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	m := Message{
		Endpoint: "/set",
		Key: "foo.key",
		Value: []byte("bar.value"),
	}

	b, errM := json.Marshal(m)
	if errM != nil {
		log.Println(errM.Error())
		os.Exit(1)
	}

	conn.Write(b)

	// read the reply
	buffer := make([]byte, 1<<24)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Ack: %s", string(buffer[:n]))
				break
			}
			log.Println(err.Error())
			break
		}
	}
}
