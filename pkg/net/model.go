package net

const (
	packetAnnounce = iota
)

type Packet struct {
	SequenceNum uint64
}

type AnnouncePacket struct {
	Packet
}
