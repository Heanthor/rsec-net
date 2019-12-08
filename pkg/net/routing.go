package net

import (
	"time"

	"github.com/Heanthor/rsec-net/internal/tcp"
	cmap "github.com/orcaman/concurrent-map"
)

type NodeInfo struct {
}

type NetInterfaceSettings struct {
	AnnounceInterval time.Duration
}

type NetInterface struct {
	net            *tcp.Net
	connectedNodes cmap.ConcurrentMap
	settings       *NetInterfaceSettings

	announceDaemonMsgChan  chan AnnouncePacket
	announceDaemonStopChan chan bool
	ErrChan                chan<- error
	MessageChan            chan interface{}
}

// NewNetInterface creates a net interface.
// addr must be of form ip:port.
// returns error if a tcp connection cannot be established.
func NewNetInterface(addr string, settings NetInterfaceSettings) (*NetInterface, error) {
	errChan := make(chan error)
	msgChan := make(chan interface{})

	n, err := tcp.NewNet(addr)
	if err != nil {
		return nil, err
	}

	recvChan, err := n.StartReceiving()
	if err != nil {
		return nil, err
	}

	m := cmap.New()

	ni := &NetInterface{
		net:                    n,
		connectedNodes:         m,
		settings:               &settings,
		announceDaemonMsgChan:  make(chan AnnouncePacket),
		announceDaemonStopChan: make(chan bool),
		ErrChan:                errChan,
		MessageChan:            msgChan,
	}

	go ni.sortReceived(recvChan)

	return ni, nil
}

func (n *NetInterface) sortReceived(recvChan <-chan interface{}) {
	for {
		// filter messages being received
		// if it's a routing protocol specific message, keep it internal
		// otherwise, forward it to the general message channel
		msgIn := <-recvChan
		if p, ok := msgIn.(AnnouncePacket); ok {
			n.announceDaemonMsgChan <- p
		} else {
			n.MessageChan <- p
		}
	}
}
