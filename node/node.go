package node

import (
	"context"

	"github.com/bchr73/carthage"
)

type NodeService struct {
	rpc *RPCClient
}

func NewNodeService(ctx context.Context, config carthage.Config) (*NodeService, error) {
	log := carthage.LoggerFromContext(ctx)

	rpc, err := NewRPCClient(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return &NodeService{
		rpc: rpc,
	}, nil
}

func (ns *NodeService) Close() {
	ns.rpc.Close()
}

func (ns *NodeService) Start(ctx context.Context) {}
