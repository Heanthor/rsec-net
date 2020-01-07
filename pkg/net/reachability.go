package net

import (
	"fmt"
	"time"

	"github.com/Heanthor/rsec-net/internal/maputils"

	"github.com/Heanthor/rsec-net/internal/udp"
	cmap "github.com/orcaman/concurrent-map"

	"github.com/rs/zerolog/log"
)

type announceDaemon struct {
	mu               udp.NetCommunicator
	announceInterval time.Duration
	errChan          chan error
	msgChan          <-chan interface{}
	stopChan         chan bool
	doneStoppingChan chan bool
	identity         Identity
	acceptOwnPackets bool

	connectedNodes cmap.ConcurrentMap

	// announce fields
	seqNo         uint16
	connNodesHash [16]byte
}

// StartAnnounceDaemon creates the announce daemon and starts its operation.
// The announce daemon does two things: periodically announces on the network, and listens for
// other announcements, updating the map of known nodes when found.
func (a *announceDaemon) StartAnnounceDaemon() {
	log.Info().Str("writeAddr", a.mu.WriteAddr()).Str("nodeName", a.identity.NodeName).Msg("Starting announce daemon...")
	a.startSending()

	time.Sleep(time.Second * 1)

	a.startReceiving()

	log.Info().Msg("Announce daemon started")
}

func (a *announceDaemon) startSending() {
	announceTicker := time.NewTicker(a.announceInterval)

	go func() {
		for {
			select {
			case <-a.stopChan:
				// don't care about waiting for this goroutine before doing other cleanup
				return
			case <-announceTicker.C:
				a.doAnnounce()
			}
		}
	}()
}

func (a *announceDaemon) startReceiving() {
	go func() {
		for {
			select {
			case <-a.stopChan:
				a.doneStoppingChan <- true

				return
			case msgIn := <-a.msgChan:
				if m, ok := msgIn.(AnnouncePacket); ok {
					if a.acceptOwnPackets || m.Identity.Addr != a.identity.Addr {
						a.handleAnnounceResponse(&m)
					}
				} else {
					log.Error().Interface("msgIn", msgIn).Msg("announce daemon got non-announce packet message")
					a.errChan <- fmt.Errorf("announce daemon got non-announce packet message")
				}
			}
		}
	}()
}

func (a *announceDaemon) StopAnnounceDaemon() {
	// two goroutines listen on this channel
	a.stopChan <- true
	a.stopChan <- true
	// wait to make sure receiving is done
	<-a.doneStoppingChan
	a.mu.StopReceiving()
	log.Debug().Msg("Announce daemon stopped")
}

func (a *announceDaemon) doAnnounce() {
	// send "i'm here" to anyone who will listen
	// if a response comes back, add them to the list of known neighbors
	// the response will be picked up in mu's receiving goroutine

	// we only update the sequence number if the message being sent is different
	// from the last sent message. we still send the message regardless
	// in case a new node has joined the network
	hash, items := maputils.ComputeHash(a.connectedNodes)
	if hash != a.connNodesHash {
		a.seqNo++
		a.connNodesHash = hash
	}

	log.Debug().Uint16("seqNo", a.seqNo).Msg("Announce daemon doing announce")
	if err := a.mu.Write(AnnouncePacket{
		Packet:         Packet{a.seqNo},
		Identity:       a.identity,
		ConnectedNodes: items,
	}); err != nil {
		a.errChan <- err
	}
}

func (a *announceDaemon) handleAnnounceResponse(ap *AnnouncePacket) {
	a.connectedNodes.SetIfAbsent(ap.NodeName, ap)
	log.Info().Interface("connectedNodes", a.connectedNodes).Msg("New connected nodes")
}
