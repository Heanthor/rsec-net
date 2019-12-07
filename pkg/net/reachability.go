package net

import (
	"time"

	"github.com/rs/zerolog/log"
)

// StartAnnounceDaemon creates the announce daemon and starts its operation.
// The announce daemon does two things: periodically announces on the network, and listens for
// other announcements, updating the map of known nodes when found.
func (n *NetInterface) StartAnnounceDaemon() error {
	log.Info().Msg("Starting announce daemon...")
	announceTicker := time.NewTicker(5 * time.Second)

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

	recvChan, err := n.net.StartReceiving()
	if err != nil {
		return err
	}

	go func(r <-chan interface{}) {
		for {
			select {
			case <-n.announceDaemonStopChan:
				n.net.Close()
			case msgIn := <-r:
				log.Debug().Interface("msgIn", msgIn).Msg("Announce daemon got message")
			default:
			}
		}
	}(recvChan)

	log.Info().Msg("Announce daemon started")

	return nil
}

func (n *NetInterface) StopAnnounceDaemon() {
	n.announceDaemonStopChan <- true
}

func (n *NetInterface) doAnnounce() {

}
