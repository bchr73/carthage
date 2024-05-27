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

	//if err := m.RPC.TxPoolSubscribe(ctx); err != nil {
	//	log.Error().Err(err).Msg(err.Error())
	//	return err
	//}

	return &NodeService{
		rpc: rpc,
	}, nil
}
