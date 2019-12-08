package udp

import (
	"net"

	"github.com/rs/zerolog/log"
)

// UniNet implements NetCommunicator for unicast UDP communication
type UniNet struct {
	addr       *net.UDPAddr
	addrString string

	stopChan         chan bool
	doneStoppingChan chan bool

	stopListener func()
	ErrChan      <-chan error
}

// NewUniNet creates a new net struct used for sending and receiving to and from the given address (hostname:port)
func NewUniNet(addr string) (*UniNet, error) {
	errChan := make(chan error)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveUDPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)
	doneStoppingChan := make(chan bool)

	return &UniNet{
		addr:             udpAddr,
		addrString:       addr,
		stopChan:         stopChan,
		doneStoppingChan: doneStoppingChan,
		ErrChan:          errChan,
		stopListener:     func() {},
	}, nil
}

// Write opens a writes a UDP datagram to the configured address and port.
func (n *UniNet) Write(data interface{}) error {
	return write(n.addr, data)
}

// StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func (n *UniNet) StartReceiving() (<-chan interface{}, error) {
	msgChan, resetFunc, err := startReceiving(n.addr, n.stopChan, n.doneStoppingChan, net.ListenUDP)
	n.stopListener = resetFunc

	return msgChan, err
}

// StopReceiving closes channels and stops the receive loop
func (n *UniNet) StopReceiving() {
	n.stopChan <- true
	<-n.doneStoppingChan
	n.stopListener()
}
