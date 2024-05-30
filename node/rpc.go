package node

import (
	"context"
	"fmt"

	"github.com/bchr73/carthage"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
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

	if !client.SupportsSubscriptions() {
		err := fmt.Errorf("connection does not support subscriptions")
		return nil, err
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

func (r *RPCClient) Close() {
	r.c.Close()
}

func (r *RPCClient) EthBlockNumber(ctx context.Context, result *rpc.BlockNumber, args ...interface{}) error {
	if err := r.c.CallContext(ctx, result, "eth_blockNumber", args...); err != nil {
		return err
	}

	return nil
}

func (r *RPCClient) TxPoolStatus(ctx context.Context, result *map[string]interface{}, args ...interface{}) error {
	if err := r.c.CallContext(ctx, result, "txpool_status", args...); err != nil {
		return err
	}

	return nil
}

func (r *RPCClient) TxPoolContent(ctx context.Context, result interface{}, args ...interface{}) error {
	if err := r.c.CallContext(ctx, result, "txpool_content", args...); err != nil {
		return err
	}

	return nil
}

type FullPendingTransactionsSubscription struct {
	ch  chan *types.Transaction
	sub *rpc.ClientSubscription
}

func (r *RPCClient) TxPoolSubscribe(ctx context.Context, nm chan *carthage.P2PMessage) (*FullPendingTransactionsSubscription, error) {
	ch := make(chan *types.Transaction)

	client := gethclient.New(r.c)
	sub, err := client.SubscribeFullPendingTransactions(ctx, ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case tx := <-ch:
				if tx == nil {
					continue
				}
				m := fmt.Sprintf("hash: %s to: %s cost: %s, nonce: %d\n", tx.Hash(), tx.To(), tx.Cost(), tx.Nonce())
				nm <- &carthage.P2PMessage{Data: []byte(m)}
			case err := <-sub.Err():
				fmt.Println(err)
				return
			}
		}
	}()

	return &FullPendingTransactionsSubscription{
		ch:  ch,
		sub: sub,
	}, nil
}
