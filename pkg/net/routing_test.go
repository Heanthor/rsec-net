package net

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Heanthor/rsec-net/internal/tcp"
	"github.com/stretchr/testify/require"
)

const addr = ":1235"

var ni *NetInterface

func init() {
	newNI, err := NewNetInterface(addr, NetInterfaceSettings{time.Second * 1})
	if err != nil {
		panic(err)
	}
	ni = newNI
}

func TestNetInterface_sortReceived(t *testing.T) {
	// send two messages, make sure they arrive in the correct channels
	n, err := tcp.NewNet(addr)
	require.NoError(t, err)

	// internal packet
	err = n.Write(AnnouncePacket{Packet{1}, "othernode"})
	require.NoError(t, err)

	// random packet for something else
	err = n.Write(Packet{2})
	require.NoError(t, err)

	adMsgIn := <-ni.ad.msgChan
	assert.Equal(t, AnnouncePacket{Packet{1}, "othernode"}, adMsgIn)

	otherMsgIn := <-ni.MessageChan
	assert.Equal(t, Packet{2}, otherMsgIn)
}
