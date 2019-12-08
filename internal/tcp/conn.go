package tcp

import (
	"encoding/gob"
	"net"

	"github.com/rs/zerolog/log"
)

type Net struct {
	addr       *net.TCPAddr
	addrString string

	stopChan         chan bool
	doneStoppingChan chan bool

	stopListener func()
	ErrChan      <-chan error
}

type Message struct {
	Data interface{}
}

// NewNet creates a new net struct used for sending and receiving to and from the given address (hostname:port)
func NewNet(addr string) (*Net, error) {
	errChan := make(chan error)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveTCPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)
	doneStoppingChan := make(chan bool)

	return &Net{
		addr:             tcpAddr,
		addrString:       addr,
		stopChan:         stopChan,
		doneStoppingChan: doneStoppingChan,
		ErrChan:          errChan,
		stopListener:     func() {},
	}, nil
}

// Write opens a connection to the saved addr, sends a message, then closes the connection.
func (n *Net) Write(data interface{}) error {
	conn, err := net.DialTCP("tcp", nil, n.addr)
	if err != nil {
		log.Error().Err(err).Msg("DialTCP failure")
		return err
	}

	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(Message{data})
	if err != nil {
		log.Error().Err(err).Msg("Write failure")
		return err
	}

	return nil
}

// StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func (n *Net) StartReceiving() (<-chan interface{}, error) {
	listener, err := net.ListenTCP("tcp", n.addr)
	if err != nil {
		log.Error().Err(err).Msg("Listen failure")
		return nil, err
	}

	dataChan := make(chan interface{})
	go func(chan interface{}) {
		for {
			select {
			case msg := <-n.stopChan:
				if msg {
					n.doneStoppingChan <- true
					return
				}
			default:
			}

			conn, err := listener.Accept()
			if err != nil {
				log.Error().Err(err).Msg("Accept failure")
			}

			log.Debug().Msg("got connection")

			var data Message
			decoder := gob.NewDecoder(conn)
			err = decoder.Decode(&data)
			if err != nil {
				log.Error().Err(err).Msg("Read failure")
			}

			err = conn.Close()
			if err != nil {
				log.Error().Err(err).Msg("Close failure")
			}

			log.Debug().Interface("message", data).Msg("StartReceiving got message")
			dataChan <- data.Data
		}
	}(dataChan)

	n.stopListener = func() {
		listener.Close()
		close(dataChan)
	}

	return dataChan, err
}

func (n *Net) StopReceiving() {
	n.stopChan <- true
	<-n.doneStoppingChan
	n.stopListener()
}
