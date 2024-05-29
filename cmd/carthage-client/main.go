package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bchr73/carthage"
	"github.com/bchr73/carthage/node"
	"github.com/bchr73/carthage/p2p"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	m := NewMain()

	if err := m.ParseFlags(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		os.Exit(1)
	}

	if err := m.Run(ctx); err != nil {
		m.Close()
		os.Exit(1)
	}
	defer m.RPC.Close()

	<-ctx.Done()

	fmt.Printf("interrupt received, exiting\n")

	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	Config      carthage.Config
	Environment string
	RPC         *node.RPCClient
}

func NewMain() *Main {
	return &Main{
		Config: DefaultConfig(),
	}
}

func (m *Main) Run(ctx context.Context) (err error) {
	ctx = carthage.NewContextWithLogger(ctx, LogLevel(m.Config.LogLevel))
	log := carthage.LoggerFromContext(ctx)

	log.Info().Msgf("running  %s", m.Environment)

	peerService, err := p2p.NewPeerService(ctx, m.Config)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}
	peerService.Start(ctx)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			peerService.Send <- &carthage.NodeMessage{Data: []byte(s)}
		}
	}()

	return nil
}

func (m *Main) Close() error {
	return nil
}

func (m *Main) ParseFlags(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("carthage", flag.ContinueOnError)
	fs.StringVar(&m.Config.LogLevel, "log.level", "warning", "debug|info|warning|error")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return nil
}

func LogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.WarnLevel
	}
}

func DefaultConfig() carthage.Config {
	var config carthage.Config
	config.P2P.ListenAddr = "/ip4/0.0.0.0/tcp/0"
	return config
}
