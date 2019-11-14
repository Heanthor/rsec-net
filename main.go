package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"sync"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"
)

var (
	addr = "127.0.0.1:8090"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("starting...")

	var wg sync.WaitGroup
	wg.Add(1)

	go write()
	go listen(&wg)

	wg.Wait()
}

func listen(wg *sync.WaitGroup) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panic().Err(err).Msg("Listen failure")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Accept failure")
		}

		log.Info().Msg("got connection")

		var b bytes.Buffer

		n, err := io.Copy(&b, conn)
		if err != nil {
			log.Error().Err(err).Msg("Read failure")
		}

		err = conn.Close()
		if err != nil {
			log.Error().Err(err).Msg("Close failure")
		}

		log.Info().Int64("bytes", n).Str("msgIn", b.String()).Msg("Got message")
		wg.Done()
	}
}

func write() {
	addr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Panic().Err(err).Msg("ResolveTCPAddr failure")
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Panic().Err(err).Msg("DialTCP failure")
	}

	defer conn.Close()

	writer := bufio.NewWriter(conn)

	n, err := writer.WriteString("hello world")
	if err != nil {
		log.Panic().Err(err).Msg("WriteString failure")
	}

	log.Info().Int("bytes", n).Msg("Wrote bytes")

	err = writer.Flush()
	if err != nil {
		log.Panic().Err(err).Msg("Flush failure")
	}

}
