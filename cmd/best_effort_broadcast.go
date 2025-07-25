package main

import (
	"flag"
	"beb/server"
	"beb/cluster"
	"log"
	"net"
	"strconv"
)

func main() {
	host := flag.String("host", "127.0.0.1", "a string")
	port := flag.Int("port", 7878, "an int")

	flag.Parse()

	var process = cluster.Process{}
	process.InitCluster()
	
	if *port != 7878 {
		joiner := cluster.NewSeed("127.0.0.1", 7878)
		joiner.PerformJoinRequest(net.JoinHostPort(*host, strconv.Itoa(*port)))
	}

	doneC := make(chan struct{})

	var serverReplica = server.NewReplica(*host, *port, doneC)
	serverReplica.StartListen()
	<- doneC
	log.Println("Done")
}
