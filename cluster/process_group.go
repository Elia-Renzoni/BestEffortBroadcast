package cluster

import (
	"net"
	"sync"
)

type ClusterManager interface {
	AddProcess()
	InitCluster()
	GetProcessGroup() []Process
}

type Process struct {
	TcpAddr *net.TCPAddr
	Conn net.Dialer
	Protocol string
}

var processGroup []Process
var lock sync.Mutex

func (p Process) AddProcess() {
	lock.Lock()
	defer lock.Unlock()

	processGroup = append(processGroup, p)
}

func (p Process) InitCluster() {
	processGroup = make([]Process, 0)
}

func (p Process) GetProcessGroup() []Process {
	lock.Lock()
	defer lock.Unlock()

	return processGroup
}
