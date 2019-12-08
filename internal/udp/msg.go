package udp

import (
	"bytes"
	"encoding/gob"
	"net"

	"github.com/rs/zerolog/log"
)

type resetFunc func()
type listenFunc func(network string, laddr *net.UDPAddr) (*net.UDPConn, error)

// write opens a writes a UDP datagram to the configured address and port.
func write(addr *net.UDPAddr, data interface{}) error {
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Error().Err(err).Msg("DialUDP failure")
		return err
	}
	defer conn.Close()
	conn.SetWriteBuffer(maxDatagramSize)

	// need to use a buffer instead of writing directly to the wire
	// since gob will attempt to send type information as separate
	// udp datagrams, which we don't want!
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(Message{data})
	if err != nil {
		log.Error().Err(err).Msg("Encode failure")
		return err
	}

	len, err := conn.Write(buf.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("Write failure")
		return err
	}

	log.Debug().Int("len", len).Interface("data", data).Msg("wrote message")

	return nil
}

// startReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func startReceiving(addr *net.UDPAddr, stopChan chan bool, doneStoppingChan chan bool, listenFunc listenFunc) (<-chan interface{}, resetFunc, error) {
	listener, err := listenFunc("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ListenUDP failure")
		return nil, nil, err
	}
	listener.SetReadBuffer(maxDatagramSize)

	dataChan := make(chan interface{})
	go func(chan interface{}) {
		for {
			select {
			case msg := <-stopChan:
				if msg {
					doneStoppingChan <- true
					return
				}
			default:
			}

			b := make([]byte, maxDatagramSize)
			len, src, err := listener.ReadFromUDP(b)
			if err != nil {
				log.Error().Err(err).Msg("Accept failure")
			}

			log.Debug().Interface("src", src).Int("len", len).Msg("got message")

			var data Message
			r := bytes.NewReader(b)
			decoder := gob.NewDecoder(r)
			err = decoder.Decode(&data)
			if err != nil {
				log.Error().Err(err).Msg("Read failure")
			}

			log.Debug().Interface("message", data).Msg("StartReceiving got message")
			dataChan <- data.Data
		}
	}(dataChan)

	stopListener := func() {
		listener.Close()
		close(dataChan)
	}

	return dataChan, stopListener, err
}
