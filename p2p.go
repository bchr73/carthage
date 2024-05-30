package carthage

type P2PMessage struct {
	Data []byte
}

type PeerService interface {
	Start()
	Close()
}
