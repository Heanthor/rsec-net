package tcp

import (
	"encoding/gob"
	"net"

	"github.com/rs/zerolog/log"
)

type Net struct {
	addr       *net.TCPAddr
	addrString string

	stopChan chan bool
}

// NewNet creates a new net struct used for sending and receiving to and from the given address (hostname:port)
func NewNet(addr string) (*Net, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveTCPAddr failure")
		return nil, err
	}

	stopChan := make(chan bool)

	return &Net{tcpAddr, addr, stopChan}, nil
}

func (n *Net) Write(data interface{}) error {
	conn, err := net.DialTCP("tcp", nil, n.addr)
	if err != nil {
		log.Error().Err(err).Msg("DialTCP failure")
		return err
	}

	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(data)
	if err != nil {
		log.Error().Err(err).Msg("WriteString failure")
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
					return
				}
			default:
			}

			conn, err := listener.Accept()
			if err != nil {
				log.Error().Err(err).Msg("Accept failure")
			}

			log.Debug().Msg("got connection")

			var data interface{}
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
			dataChan <- data
		}
	}(dataChan)

	return dataChan, err
}

func (n *Net) Close() {
	n.stopChan <- true
}
