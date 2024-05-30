package p2p

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"

	"github.com/bchr73/carthage"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

var topicName = SHA256("carthagenodemessagetopic_main")

type PeerService struct {
	host    host.Host
	dht     *dht.IpfsDHT
	routing *drouting.RoutingDiscovery
	pbsb    *pubsub.PubSub

	Recv chan *carthage.PeerMessage
	Send chan *carthage.PeerMessage
}

func NewPeerService(ctx context.Context, config carthage.Config) (*PeerService, error) {
	host, err := libp2p.New(libp2p.ListenAddrStrings(config.P2P.ListenAddr))
	if err != nil {
		return nil, err
	}

	kad, err := NewKademliaDHT(ctx, host)
	if err != nil {
		return nil, err
	}

	routing := drouting.NewRoutingDiscovery(kad)

	pbsb, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	recv := make(chan *carthage.PeerMessage, 10000)
	send := make(chan *carthage.PeerMessage, 10000)

	return &PeerService{
		host:    host,
		dht:     kad,
		routing: routing,
		pbsb:    pbsb,
		Recv:    recv,
		Send:    send,
	}, nil
}

func (ps *PeerService) Close() {}

func (ps *PeerService) Start(ctx context.Context) {
	log := carthage.LoggerFromContext(ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := ps.discoverPeers(ctx, topicName)
		if err != nil {
			log.Error().Err(err).Msg(err.Error())
			return
		}
	}()
	wg.Wait()

	topic, err := ps.joinTopic(topicName)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	go func() {
		for m := range ps.Send {
			err := topic.Publish(ctx, m.Data)
			if err != nil {
				log.Error().Msgf("publish error %s", err)
			}
		}
	}()

	sub, err := ps.subscribeTopic(topic)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return
	}

	go func() {
		for {
			m, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}
			ps.Recv <- &carthage.PeerMessage{Data: m.Data}
		}
	}()
}

func (ps *PeerService) discoverPeers(ctx context.Context, topicName string) error {
	log := carthage.LoggerFromContext(ctx)

	dutil.Advertise(ctx, ps.routing, topicName)

	anyConnected := false
	for !anyConnected {
		log.Info().Msgf("searching peers for topic %s", topicName)
		ch, err := ps.routing.FindPeers(ctx, topicName)
		if err != nil {
			return err
		}
		for peer := range ch {
			if peer.ID == ps.host.ID() {
				continue // no self connection
			}
			if err := ps.host.Connect(ctx, peer); err != nil {
				log.Debug().Msgf("failed connecting to peer %s, error: %s", peer.ID, err)
				continue
			}
			log.Info().Msgf("connected to peer %s", peer.ID)
			anyConnected = true
		}
	}
	log.Info().Msg("peer discovery complete")

	return nil
}

func (p2 *PeerService) joinTopic(topicName string) (*pubsub.Topic, error) {
	topic, err := p2.pbsb.Join(topicName)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (p2 *PeerService) subscribeTopic(topic *pubsub.Topic) (*pubsub.Subscription, error) {
	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func SHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
