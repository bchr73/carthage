package node

import (
	"context"
	"fmt"

	"github.com/bchr73/carthage"
	"github.com/ethereum/go-ethereum/rpc"
)

type RPCClient struct {
	c *rpc.Client
}

func NewRPCClient(ctx context.Context, config carthage.Config) (*RPCClient, error) {
	var client *rpc.Client
	var err error

	log := carthage.LoggerFromContext(ctx)

	if config.RPC.RpcWs {
		client, err = rpc.DialWebsocket(ctx, config.RPC.RpcUri, "")
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("using rpc ws at %s", config.RPC.RpcUri)
	}
	if config.RPC.RpcIpc {
		client, err = rpc.DialIPC(ctx, config.RPC.RpcUri)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("using rpc ipc at %s", config.RPC.RpcUri)
	}

	modules, err := client.SupportedModules()
	if err != nil {
		return nil, err
	}

	if _, ok := modules["txpool"]; !ok {
		err := fmt.Errorf("connection does not support 'txpool' module")
		return nil, err
	}

	return &RPCClient{
		c: client,
	}, nil
}
