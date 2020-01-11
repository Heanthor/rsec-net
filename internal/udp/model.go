package udp

// NetReader binds to and continuously reads from the given host and port
type NetReader interface {
	StartReceiving(string) (<-chan interface{}, error)
	StopReceiving()
	ReadAddr() string
}

// NetWriter writes packets to the network
type NetWriter interface {
	Write(data interface{}) error
	WriteAddr() string
}

// Message is a basic serializable message
type Message struct {
	Data interface{}
}

const maxDatagramSize = 8192
