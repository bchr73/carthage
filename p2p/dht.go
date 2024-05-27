package p2p

import (
	"context"
	"sync"

	"github.com/bchr73/carthage"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func NewKademliaDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {
	log := carthage.LoggerFromContext(ctx)

	kad, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	err = kad.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerInfo); err != nil {
				log.Error().Err(err).Msg(err.Error())
			}
		}()

	}
	wg.Wait()

	return kad, nil
}
