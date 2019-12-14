package net

const (
	uniAddr = ":1145"
)

var iface *Interface

// func init() {
// 	fmt.Println("Running routing_test init")
// 	i, err := NewInterface("initNode", uniAddr, multiAddr, InterfaceSettings{time.Second * 1})
// 	if err != nil {
// 		panic(err)
// 	}

// 	iface = i
// }

// func TestInterface_Announce(t *testing.T) {
// 	testiface, err := NewInterface("testNode", uniAddr, multiAddr, InterfaceSettings{time.Second * 1})
// 	require.NoError(t, err)

// 	testiface.MessageChan
// }
