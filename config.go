package carthage

type Config struct {
	LogLevel string

	RPC struct {
		RpcWs  bool
		RpcIpc bool
		RpcUri string
	}

	P2P struct {
		ListenAddr string
	}
}
