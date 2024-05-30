package carthage

type PeerMessage struct {
	Data []byte
}

type PeerService interface {
	Start()
	Close()
}
