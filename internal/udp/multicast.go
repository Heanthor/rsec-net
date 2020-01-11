package udp

import (
	"net"

	"github.com/rs/zerolog/log"
)

// MulticastReader implements NetReader for multicast UDP communication
type MulticastReader struct {
	addr       *net.UDPAddr
	addrString string

	stopChan         chan bool
	doneStoppingChan chan bool

	stopListener func()
	ErrChan      <-chan error
}

// NewMulticastReader creates a new net struct used for receiving from the given address (hostname:port)
func NewMulticastReader(addr string) (*MulticastReader, error) {
	errChan := make(chan error)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveUDPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)
	doneStoppingChan := make(chan bool)

	return &MulticastReader{
		addr:             udpAddr,
		addrString:       addr,
		stopChan:         stopChan,
		doneStoppingChan: doneStoppingChan,
		ErrChan:          errChan,
		stopListener:     func() {},
	}, nil
}

// StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func (n *MulticastReader) StartReceiving(tag string) (<-chan interface{}, error) {
	listenFunc := func(network string, gaddr *net.UDPAddr) (*net.UDPConn, error) {
		return net.ListenMulticastUDP(network, nil, gaddr)
	}

	msgChan, resetFunc, err := startReceiving(n.addr, n.stopChan, n.doneStoppingChan, listenFunc, tag)
	n.stopListener = resetFunc

	return msgChan, err
}

// StopReceiving closes channels and stops the receive loop
func (n *MulticastReader) StopReceiving() {
	n.stopChan <- true
	<-n.doneStoppingChan
	n.stopListener()
	log.Debug().Msg("multicast stopped receiving")
}

func (n *MulticastReader) ReadAddr() string {
	return n.addrString
}
