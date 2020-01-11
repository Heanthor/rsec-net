package udp

import (
	"net"

	"github.com/rs/zerolog/log"
)

// UniReader implements NetReader for unicast UDP communication
type UniReader struct {
	addr       *net.UDPAddr
	addrString string

	stopChan         chan bool
	doneStoppingChan chan bool

	stopListener func()
	ErrChan      <-chan error
}

// NewUniReader creates a new net struct used for receiving from the given address (hostname:port)
func NewUniReader(addr string) (*UniReader, error) {
	errChan := make(chan error)
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveUDPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)
	doneStoppingChan := make(chan bool)

	return &UniReader{
		addr:             udpAddr,
		addrString:       addr,
		stopChan:         stopChan,
		doneStoppingChan: doneStoppingChan,
		ErrChan:          errChan,
		stopListener:     func() {},
	}, nil
}

// StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func (n *UniReader) StartReceiving(tag string) (<-chan interface{}, error) {
	msgChan, resetFunc, err := startReceiving(n.addr, n.stopChan, n.doneStoppingChan, net.ListenUDP, tag)
	n.stopListener = resetFunc

	return msgChan, err
}

// StopReceiving closes channels and stops the receive loop
func (n *UniReader) StopReceiving() {
	n.stopChan <- true
	<-n.doneStoppingChan
	n.stopListener()
	log.Debug().Msg("unicast stopped receiving")
}

// ReadAddr returns the address being read from (host:port)
func (n *UniReader) ReadAddr() string {
	return n.addrString
}
