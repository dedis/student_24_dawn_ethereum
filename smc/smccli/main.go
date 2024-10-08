package main

import (
	"fmt"
	"io"
	"os"

	"go.dedis.ch/dela/cli/node"

	db "go.dedis.ch/dela/core/store/kv/controller"
	dkg "go.dedis.ch/f3b/smc/dkg/controller"
	mino "go.dedis.ch/dela/mino/minogrpc/controller"

	_ "go.dedis.ch/f3b/smc/dkg/json"
)

func main() {
	err := run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	return runWithCfg(args, config{Writer: os.Stdout})
}

type config struct {
	Channel chan os.Signal
	Writer  io.Writer
}

func runWithCfg(args []string, cfg config) error {
	builder := node.NewBuilderWithCfg(
		cfg.Channel,
		cfg.Writer,
		db.NewController(),
		mino.NewController(),
		dkg.NewMinimal(),
	)

	app := builder.Build()

	err := app.Run(args)
	if err != nil {
		return err
	}

	return nil
}
