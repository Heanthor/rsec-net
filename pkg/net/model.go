package net

import (
	"encoding/gob"
)

const (
	packetAnnounce = iota
)

type Packet struct {
	SequenceNum uint64
}

type AnnouncePacket struct {
	Packet
	RemoteName string
}

func init() {
	gob.Register(Packet{})
	gob.Register(AnnouncePacket{})
}
