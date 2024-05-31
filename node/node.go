package node

import (
	"context"

	"github.com/bchr73/carthage"
	"github.com/ethereum/go-ethereum/rpc"
)

type NodeService struct {
	rpc *RPCClient

	req chan *carthage.RPCCall
}

func NewNodeService(ctx context.Context, config carthage.Config) (*NodeService, error) {
	log := carthage.LoggerFromContext(ctx)

	rpc, err := NewRPCClient(ctx, config)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	req := make(chan *carthage.RPCCall, 10000)

	return &NodeService{
		rpc: rpc,
		req: req,
	}, nil
}

func (ns *NodeService) Close() {
	ns.rpc.Close()
}

func (ns *NodeService) Call(call *carthage.RPCCall) error {
	if ok, err := call.Validate(); ok == false {
		return err
	}
	ns.req <- call

	return nil
}

func (ns *NodeService) Start(ctx context.Context, res chan *carthage.RPCResult) {
	log := carthage.LoggerFromContext(ctx)

	for {
		select {
		case call := <-ns.req:
			switch call.Method {
			case "eth_blockNumber":
				go func() {
					var result rpc.BlockNumber
					ns.rpc.EthBlockNumber(ctx, &result)
					res <- &carthage.RPCResult{
						Method: call.Method,
						Result: &result,
					}
				}()
			case "txpool_status":
				go func() {
					var result map[string]interface{}
					ns.rpc.TxPoolStatus(ctx, &result)
					res <- &carthage.RPCResult{
						Method: call.Method,
						Result: &result,
					}
				}()
			case "txpool_content":
				go func() {
					var result interface{}
					ns.rpc.TxPoolContent(ctx, &result)
					res <- &carthage.RPCResult{
						Method: call.Method,
						Result: &result,
					}
				}()
			default:
				log.Error().Msgf("method call '%s' unrecognized", call.Method)
			}
		}
	}
}
