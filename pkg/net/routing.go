package net

import (
	"time"

	"github.com/Heanthor/rsec-net/internal/udp"
	cmap "github.com/orcaman/concurrent-map"
)

// NodeInfo contains information about a discovered network node
type NodeInfo struct {
	NodeName  string
	Addr      string
	Latency   latency
	lastSeqNo uint16
}

// InterfaceSettings contains settings for the net interface
type InterfaceSettings struct {
	AnnounceInterval time.Duration
}

// Interface maintains connectivity with the mesh network,
// and provides functions for sending and receiving on the network.
// TODO:
// function for sending to addr, send to node name
// receive from node, receive from all
type Interface struct {
	uni      *udp.UniNet
	settings *InterfaceSettings
	ad       *announceDaemon

	ErrChan     chan<- error
	MessageChan <-chan interface{}
}

// NewInterface creates a net interface.
// addr must be of form ip:port.
// returns error if udp address resolution fails.
func NewInterface(nodeName, addr, multicastAddr string, settings InterfaceSettings) (*Interface, error) {
	errChan := make(chan error)

	// create unicast and multicast communicators
	u, err := udp.NewUniNet(addr)
	if err != nil {
		return nil, err
	}

	recvChan, err := u.StartReceiving()
	if err != nil {
		return nil, err
	}

	mu, err := udp.NewMulticastNet(multicastAddr)
	if err != nil {
		return nil, err
	}

	mRecvChan, err := mu.StartReceiving()
	if err != nil {
		return nil, err
	}

	m := cmap.New()

	return &Interface{
		uni:         u,
		settings:    &settings,
		ErrChan:     errChan,
		MessageChan: recvChan,
		ad: &announceDaemon{
			identity:         Identity{nodeName, addr},
			mu:               mu,
			errChan:          errChan,
			announceInterval: settings.AnnounceInterval,
			msgChan:          mRecvChan,
			stopChan:         make(chan bool),
			doneStoppingChan: make(chan bool),
			connectedNodes:   m,
			acceptOwnPackets: false,
		},
	}, nil
}

// Close stops the announce daemon and closes all open connections and channels
func (n *Interface) Close() {
	n.uni.StopReceiving()
	n.ad.StopAnnounceDaemon()
	close(n.ErrChan)
}
