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
	dataSend    udp.NetWriter
	dataReceive udp.NetReader

	announceReceive udp.NetReader

	settings *InterfaceSettings
	ad       *announceDaemon

	ErrChan     chan<- error
	MessageChan <-chan interface{}
}

// NewInterface creates a net interface.
// addr must be of form ip:port.
// returns error if udp address resolution fails.
func NewInterface(nodeName string, dataReceive udp.NetReader, announceSend udp.NetWriter, announceReceive udp.NetReader, settings InterfaceSettings) (*Interface, error) {
	errChan := make(chan error)

	// TODO create data sender when a recipient is determined

	recvChan, err := dataReceive.StartReceiving()
	if err != nil {
		return nil, err
	}

	mRecvChan, err := announceReceive.StartReceiving()
	if err != nil {
		return nil, err
	}

	m := cmap.New()

	return &Interface{
		dataReceive:     dataReceive,
		announceReceive: announceReceive,
		settings:        &settings,
		ErrChan:         errChan,
		MessageChan:     recvChan,
		ad: &announceDaemon{
			identity:         Identity{nodeName, dataReceive.ReadAddr()}, // TODO what is my external ip?
			w:                announceSend,
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

// StartAnnounce starts announcing the node to the network
func (n *Interface) StartAnnounce() {
	n.ad.StartAnnounceDaemon()
}

// Close stops the announce daemon and closes all open connections and channels
func (n *Interface) Close() {
	n.dataReceive.StopReceiving()
	n.announceReceive.StopReceiving()
	n.ad.StopAnnounceDaemon()
	close(n.ErrChan)
}
