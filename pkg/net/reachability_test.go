package net

import (
	"fmt"
	"testing"
	"time"

	"github.com/Heanthor/rsec-net/internal/udp"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnnounceDaemonSuite struct {
	suite.Suite
	addr   string
	daemon *announceDaemon
}

func (suite *AnnounceDaemonSuite) SetupSuite() {
	suite.addr = ":1146"
	// suite.daemon = initNewAnnounceDaemon("suiteDaemon", suite.addr, time.Second*1)
	// suite.daemon.StartAnnounceDaemon()
	fmt.Println("suite setup")
}

func (suite *AnnounceDaemonSuite) TearDownSuite() {
	// suite.daemon.StopAnnounceDaemon()
}

func (suite *AnnounceDaemonSuite) Test_DaemonAnnounce() {
	writeDaemon := initWriteOnlyNewAnnounceDaemon("writeDaemon", suite.addr, time.Second*1)
	fakeConnNodes := cmap.New()
	fakeConnNodes.Set("unknownNode", AnnouncePacket{
		Packet:   Packet{0},
		Identity: Identity{"unknownNode", ":2222"},
	})
	writeDaemon.connectedNodes = fakeConnNodes

	// receiving end
	testDaemon := initNewAnnounceDaemon("testDaemon", suite.addr, time.Second*1)
	testDaemon.startReceiving()

	time.Sleep(time.Second * 1)

	// manually trigger announce, expect testDaemon to receive the message
	writeDaemon.doAnnounce()

	time.Sleep(time.Second * 2)

	assert.Contains(suite.T(), testDaemon.connectedNodes.Keys(), "writeDaemon")
}

func TestAnnounceDaemonSuite(t *testing.T) {
	suite.Run(t, new(AnnounceDaemonSuite))
}

func initNewAnnounceDaemon(nodeName, addr string, announceInterval time.Duration) *announceDaemon {
	errChan := make(chan error)

	// even though reachability is through multicast, test with unicast
	w, err := udp.NewUDPWriter(addr)
	if err != nil {
		panic(err)
	}

	rd, err := udp.NewUniReader(addr)
	if err != nil {
		panic(err)
	}
	mRecvChan, err := rd.StartReceiving()
	if err != nil {
		panic(err)
	}

	m := cmap.New()

	return &announceDaemon{
		identity:         Identity{nodeName, addr},
		w:                w,
		errChan:          errChan,
		announceInterval: announceInterval,
		msgChan:          mRecvChan,
		stopChan:         make(chan bool),
		doneStoppingChan: make(chan bool),
		connectedNodes:   m,
		acceptOwnPackets: true,
	}
}

func initWriteOnlyNewAnnounceDaemon(nodeName, addr string, announceInterval time.Duration) *announceDaemon {
	errChan := make(chan error)

	w, err := udp.NewUDPWriter(addr)
	if err != nil {
		panic(err)
	}

	fakeRecvChan := make(chan interface{})

	m := cmap.New()

	return &announceDaemon{
		identity:         Identity{nodeName, addr},
		w:                w,
		errChan:          errChan,
		announceInterval: announceInterval,
		msgChan:          fakeRecvChan,
		stopChan:         make(chan bool),
		doneStoppingChan: make(chan bool),
		connectedNodes:   m,
	}
}
