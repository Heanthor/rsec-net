package net

import (
	"time"

	"github.com/rs/zerolog/log"
)

// StartAnnounceDaemon creates the announce daemon and starts its operation.
// The announce daemon does two things: periodically announces on the network, and listens for
// other announcements, updating the map of known nodes when found.
func (n *NetInterface) StartAnnounceDaemon() {
	log.Info().Msg("Starting announce daemon...")
	announceTicker := time.NewTicker(n.settings.AnnounceInterval)

	go func() {
		for {
			select {
			case <-n.announceDaemonStopChan:
				return
			case <-announceTicker.C:
				n.doAnnounce()
			}
		}
	}()

	go func(r <-chan interface{}) {
		for {
			select {
			case <-n.announceDaemonStopChan:
				return
			case msgIn := <-r:
				if p, ok := msgIn.(AnnouncePacket); ok {
					n.HandleAnnounceResponse(&p)
				}
				log.Debug().Interface("msgIn", msgIn).Msg("Announce daemon got message")
			default:
			}
		}
	}(n.MessageChan)

	log.Info().Msg("Announce daemon started")
}

func (n *NetInterface) StopAnnounceDaemon() {
	n.announceDaemonStopChan <- true
}

func (n *NetInterface) doAnnounce() {
	// send "i'm here" to anyone who will listen
	// if a response comes back, add them to the list of known neighbors
	// the response will be picked up in n's receiving goroutine
	err := n.net.Write(AnnouncePacket{Packet: Packet{packetAnnounce}})
	if err != nil {
		n.ErrChan <- err
	}
}

func (n *NetInterface) HandleAnnounceResponse(a *AnnouncePacket) {

}
