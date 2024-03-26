package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bchr73/carthage"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

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
}

func NewMain() *Main {
	return &Main{
		Config: DefaultConfig(),
	}
}

func (m *Main) Run(ctx context.Context) (err error) {
	fmt.Printf("running %s\n", m.Environment)
	return nil
}

func (m *Main) Close() error {
	return nil
}

func (m *Main) ParseFlags(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("carthage", flag.ContinueOnError)
	fs.StringVar(&m.Environment, "env", "testnet", "environment")

	if err := fs.Parse(args); err != nil {
		return err
	}

	return nil
}

func DefaultConfig() carthage.Config {
	return carthage.Config{}
}
