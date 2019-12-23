package udp

import (
	"net"

	"github.com/rs/zerolog/log"
)

// MulticastNet implements NetCommunicator for multicast UDP communication
type MulticastNet struct {
	addr       *net.UDPAddr
	addrString string

	stopChan         chan bool
	doneStoppingChan chan bool

	stopListener func()
	ErrChan      <-chan error
}

// NewMulticastNet creates a new net struct used for sending and receiving to and from the given address (hostname:port)
func NewMulticastNet(addr string) (*MulticastNet, error) {
	errChan := make(chan error)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveUDPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)
	doneStoppingChan := make(chan bool)

	return &MulticastNet{
		addr:             udpAddr,
		addrString:       addr,
		stopChan:         stopChan,
		doneStoppingChan: doneStoppingChan,
		ErrChan:          errChan,
		stopListener:     func() {},
	}, nil
}

// Write opens a writes a UDP datagram to the configured address and port.
func (n *MulticastNet) Write(data interface{}) error {
	return write(n.addr, data)
}

// StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func (n *MulticastNet) StartReceiving() (<-chan interface{}, error) {
	listenFunc := func(network string, gaddr *net.UDPAddr) (*net.UDPConn, error) {
		return net.ListenMulticastUDP(network, nil, gaddr)
	}

	msgChan, resetFunc, err := startReceiving(n.addr, n.stopChan, n.doneStoppingChan, listenFunc)
	n.stopListener = resetFunc

	return msgChan, err
}

// StopReceiving closes channels and stops the receive loop
func (n *MulticastNet) StopReceiving() {
	n.stopChan <- true
	<-n.doneStoppingChan
	n.stopListener()
	log.Debug().Msg("multicast stopped receiving")
}

func (n *MulticastNet) Addr() string {
	return n.addrString
}
