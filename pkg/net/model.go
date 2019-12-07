package net

const (
	packetAnnounce = iota
)

type Packet struct {
	PacketType  uint8
	SequenceNum uint64
}

type AnnouncePacket struct {
	Packet
}
