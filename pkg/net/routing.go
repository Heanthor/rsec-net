package net

import (
	"time"

	"github.com/Heanthor/rsec-net/internal/udp"
	cmap "github.com/orcaman/concurrent-map"
)

type NodeInfo struct {
}

type NetInterfaceSettings struct {
	AnnounceInterval time.Duration
}

type NetInterface struct {
	uni            *udp.UniNet
	connectedNodes cmap.ConcurrentMap
	settings       *NetInterfaceSettings
	ad             *announceDaemon

	ErrChan     chan<- error
	MessageChan <-chan interface{}
}

// NewNetInterface creates a net interface.
// addr must be of form ip:port.
// returns error if a tcp connection cannot be established.
func NewNetInterface(addr, multicastAddr string, settings NetInterfaceSettings) (*NetInterface, error) {
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

	return &NetInterface{
		uni:            u,
		connectedNodes: m,
		settings:       &settings,
		ErrChan:        errChan,
		MessageChan:    recvChan,
		ad: &announceDaemon{
			mu:               mu,
			errChan:          errChan,
			announceInterval: settings.AnnounceInterval,
			msgChan:          mRecvChan,
			stopChan:         make(chan bool),
			doneStoppingChan: make(chan bool),
		},
	}, nil
}

// Close stops the announce daemon and closes all open connections and channels
func (n *NetInterface) Close() {
	n.uni.StopReceiving()
	n.ad.StopAnnounceDaemon()
	close(n.ErrChan)
}
