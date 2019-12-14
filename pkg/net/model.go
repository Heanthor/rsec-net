package net

import (
	"encoding/gob"
)

const (
	packetAnnounce = iota
)

// Packet is the basic packet struct
type Packet struct {
	SequenceNum uint16
}

// Identity contains information to identify a struct
type Identity struct {
	NodeName string
	Addr     string
}

// AnnouncePacket contains information about the current node to send to other nodes
type AnnouncePacket struct {
	Packet
	Identity
}

func init() {
	gob.Register(Packet{})
	gob.Register(AnnouncePacket{})
}
