package net

import (
	"github.com/Heanthor/rsec-net/internal/tcp"
	cmap "github.com/orcaman/concurrent-map"
)

type NodeInfo struct {
}

type NetInterface struct {
	net            *tcp.Net
	connectedNodes cmap.ConcurrentMap

	announceDaemonStopChan chan bool
}

func NewNetInterface(addr string) (*NetInterface, error) {
	n, err := tcp.NewNet(addr)
	if err != nil {
		return nil, err
	}

	m := cmap.New()

	return &NetInterface{
		net:                    n,
		connectedNodes:         m,
		announceDaemonStopChan: make(chan bool),
	}, nil
}

func (ni *NetInterface) Announce() error {
	err := ni.net.Write("abc")
	if err != nil {
		return err
	}

	return nil
}
