package udp

// NetCommunicator is a struct that can be used to send and receive messages
type NetCommunicator interface {
	Write(data interface{}) error
	StartReceiving() (<-chan interface{}, error)
	StopReceiving()
	Addr() string
}

// Message is a basic serializable message
type Message struct {
	Data interface{}
}

const maxDatagramSize = 8192
