package carthage

import "context"

type RPCCall struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
}

type NodeService interface {
	Start(ctx context.Context)
	Close()
}
