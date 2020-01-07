package udp

import (
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

const addr = ":1145"

var (
	writer     *UDPWriter
	testReader *UniReader
	recvChan   <-chan interface{}
)

type s struct {
	St string
}

func init() {
	w, err := NewUDPWriter(addr)
	if err != nil {
		panic(err)
	}
	n, err := NewUniReader(addr)
	if err != nil {
		panic(err)
	}

	recv, err := n.StartReceiving()
	if err != nil {
		panic(err)
	}

	writer = w
	testReader = n
	recvChan = recv

	// needed to send the type through gob
	gob.Register(s{})
}

func TestNet_SendReceive(t *testing.T) {
	// use another Net to send to the running server
	w, err := NewUDPWriter(addr)
	require.NoError(t, err)

	err = w.Write(s{"hello"})
	require.NoError(t, err)

	msgIn := <-recvChan
	sIn, ok := msgIn.(s)
	assert.True(t, ok)
	assert.Equal(t, "hello", sIn.St)
}

func TestNet_SendReceiveMultiple(t *testing.T) {
	// use another Net to send to the running server
	w, err := NewUDPWriter(addr)
	require.NoError(t, err)

	err = w.Write(s{"hello"})
	require.NoError(t, err)

	err = w.Write(s{"goodbye"})
	require.NoError(t, err)

	msgIn := <-recvChan
	sIn, ok := msgIn.(s)
	assert.True(t, ok)
	assert.Equal(t, "hello", sIn.St)

	msgIn2 := <-recvChan
	sIn2, ok := msgIn2.(s)
	assert.True(t, ok)
	assert.Equal(t, "goodbye", sIn2.St)
}
