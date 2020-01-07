package udp

import (
	"bytes"
	"encoding/gob"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type resetFunc func()
type listenFunc func(network string, laddr *net.UDPAddr) (*net.UDPConn, error)

// UDPWriter contains methods to write to a udp address (host:port)
type UDPWriter struct {
	addr       *net.UDPAddr
	addrString string
}

// NewUDPWriter creates a new writer that writes to the given address (host:port)
func NewUDPWriter(addr string) (*UDPWriter, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ResolveUDPAddr failure")
		return nil, err
	}

	return &UDPWriter{
		addr:       udpAddr,
		addrString: addr,
	}, nil
}

// write opens a writes a UDP datagram to the configured address and port.
func (u *UDPWriter) Write(data interface{}) error {
	conn, err := net.DialUDP("udp4", nil, u.addr)
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

// WriteAddr returns the address written to
func (u *UDPWriter) WriteAddr() string {
	return u.addrString
}

// startReceiving starts listening on the Net, and returns a channel which will yield messages when they arrive.
func startReceiving(addr *net.UDPAddr, stopChan chan bool, doneStoppingChan chan bool, listenFunc listenFunc) (<-chan interface{}, resetFunc, error) {
	listener, err := listenFunc("udp4", addr)
	if err != nil {
		log.Error().Err(err).Msg("ListenUDP failure")
		return nil, nil, err
	}
	listener.SetReadBuffer(maxDatagramSize)

	dataChan := make(chan interface{}, 1)
	go func(chan interface{}) {
		for {
			select {
			case <-stopChan:
				doneStoppingChan <- true
				return
			default:
			}
			listener.SetReadDeadline(time.Now().Add(time.Second * 2))

			b := make([]byte, maxDatagramSize)
			len, src, err := listener.ReadFromUDP(b)
			if err != nil && strings.Index(err.Error(), "i/o timeout") < 0 {
				log.Error().Err(err).Msg("Accept failure")
				continue
			}

			if len == 0 {
				// TODO figure out why this happens
				continue
			}

			log.Debug().Interface("src", src).Int("len", len).Msg("got message")

			var data Message
			r := bytes.NewReader(b)
			decoder := gob.NewDecoder(r)
			err = decoder.Decode(&data)
			if err != nil {
				log.Error().Err(err).Msg("Read failure")
				continue
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
