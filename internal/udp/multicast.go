package udp

// import (
// 	"bytes"
// 	"encoding/gob"
// 	"net"

// 	"github.com/rs/zerolog/log"
// )

// // UniNet implements NetCommunicator for unicast UDP communication
// type MulticastNet struct {
// 	addr       *net.UDPAddr
// 	addrString string

// 	stopChan         chan bool
// 	doneStoppingChan chan bool

// 	stopListener func()
// 	ErrChan      <-chan error
// }

// func NewMulticastNet(addr string) (*MulticastNet, error) {
// 	errChan := make(chan error)
// 	tcpAddr, err := net.ResolveUDPAddr("udp4", addr)
// 	if err != nil {
// 		log.Error().Err(err).Msg("ResolveUDPAddr failure")
// 		return nil, err
// 	}

// 	stopChan := make(chan bool)
// 	doneStoppingChan := make(chan bool)

// 	return &MulticastNet{
// 		addr:             tcpAddr,
// 		addrString:       addr,
// 		stopChan:         stopChan,
// 		doneStoppingChan: doneStoppingChan,
// 		ErrChan:          errChan,
// 		stopListener:     func() {},
// 	}, nil
// }

// // Write opens a writes a UDP datagram to the configured address and port.
// func (n *UniNet) Write(data interface{}) error {
// 	conn, err := net.DialUDP("udp4", nil, n.addr)
// 	if err != nil {
// 		log.Error().Err(err).Msg("DialUDP failure")
// 		return err
// 	}
// 	defer conn.Close()
// 	conn.SetWriteBuffer(maxDatagramSize)

// 	// need to use a buffer instead of writing directly to the wire
// 	// since gob will attempt to send type information as separate
// 	// udp datagrams, which we don't want!
// 	var buf bytes.Buffer
// 	encoder := gob.NewEncoder(&buf)
// 	err = encoder.Encode(Message{data})
// 	if err != nil {
// 		log.Error().Err(err).Msg("Encode failure")
// 		return err
// 	}

// 	len, err := conn.Write(buf.Bytes())
// 	if err != nil {
// 		log.Error().Err(err).Msg("Write failure")
// 		return err
// 	}

// 	log.Debug().Int("len", len).Interface("data", data).Msg("wrote message")

// 	return nil
// }

// // StartReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
// func (n *UniNet) StartReceiving() (<-chan interface{}, error) {
// 	listener, err := net.ListenUDP("udp4", n.addr)
// 	if err != nil {
// 		log.Error().Err(err).Msg("ListenUDP failure")
// 		return nil, err
// 	}
// 	listener.SetReadBuffer(maxDatagramSize)

// 	dataChan := make(chan interface{})
// 	go func(chan interface{}) {
// 		for {
// 			select {
// 			case msg := <-n.stopChan:
// 				if msg {
// 					n.doneStoppingChan <- true
// 					return
// 				}
// 			default:
// 			}

// 			b := make([]byte, maxDatagramSize)
// 			len, src, err := listener.ReadFromUDP(b)
// 			if err != nil {
// 				log.Error().Err(err).Msg("Accept failure")
// 			}

// 			log.Debug().Interface("src", src).Int("len", len).Msg("got message")

// 			var data Message
// 			r := bytes.NewReader(b)
// 			decoder := gob.NewDecoder(r)
// 			err = decoder.Decode(&data)
// 			if err != nil {
// 				log.Error().Err(err).Msg("Read failure")
// 			}

// 			log.Debug().Interface("message", data).Msg("StartReceiving got message")
// 			dataChan <- data.Data
// 		}
// 	}(dataChan)

// 	n.stopListener = func() {
// 		listener.Close()
// 		close(dataChan)
// 	}

// 	return dataChan, err
// }

// func (n *UniNet) StopReceiving() {
// 	n.stopChan <- true
// 	<-n.doneStoppingChan
// 	n.stopListener()
// }
