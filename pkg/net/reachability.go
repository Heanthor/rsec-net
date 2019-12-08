package net

import (
	"time"

	"github.com/rs/zerolog/log"
)

type announceDaemon struct {
	NetInterface
	announceInterval time.Duration
	msgChan          chan AnnouncePacket
	stopChan         chan bool
	doneStoppingChan chan bool
}

// StartAnnounceDaemon creates the announce daemon and starts its operation.
// The announce daemon does two things: periodically announces on the network, and listens for
// other announcements, updating the map of known nodes when found.
func (a *announceDaemon) StartAnnounceDaemon() {
	log.Info().Msg("Starting announce daemon...")
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

	go func() {
		for {
			select {
			case <-a.stopChan:
				a.doneStoppingChan <- true

				return
			case msgIn := <-a.msgChan:
				a.HandleAnnounceResponse(&msgIn)
				log.Debug().Interface("msgIn", msgIn).Msg("Announce daemon got message")
			default:
			}
		}
	}()

	log.Info().Msg("Announce daemon started")
}

func (a *announceDaemon) StopAnnounceDaemon() {
	a.stopChan <- true
	// wait to make sure receiving is done
	<-a.doneStoppingChan
	close(a.msgChan)
	close(a.ErrChan)
}

func (a *announceDaemon) doAnnounce() {
	// send "i'm here" to anyone who will listen
	// if a response comes back, add them to the list of known neighbors
	// the response will be picked up in n's receiving goroutine
	err := a.net.Write(AnnouncePacket{Packet: Packet{packetAnnounce}})
	if err != nil {
		a.ErrChan <- err
	}
}

func (n *NetInterface) HandleAnnounceResponse(a *AnnouncePacket) {

}
