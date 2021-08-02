package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.ibm.scs.com/tapestry/cmd/tapestry/cli"
)

var (
	rootFlagSet    = flag.NewFlagSet("tapestry", flag.ExitOnError)
	debug          = rootFlagSet.Bool("d", false, "log debug output")
	outputFilename = rootFlagSet.String("output-file", "", "log output to a file")
)

func main() {
	root := &ffcli.Command{
		ShortUsage: "tapestry [flags] <subcommand>",
		FlagSet:    rootFlagSet,
		Subcommands: []*ffcli.Command{
			// Pipeline Options
			cli.Tkn(),
			cli.Gha(),
			// Version
			cli.Version()},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.Parse(os.Args[1:]); err != nil {
		printErrAndExit(err)
	}

	if *debug {
		logs.Debug.SetOutput(os.Stderr)
	}

	if err := root.Run(context.Background()); err != nil {
		printErrAndExit(err)
	}
}

func printErrAndExit(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
