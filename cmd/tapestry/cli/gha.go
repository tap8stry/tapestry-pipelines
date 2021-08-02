package cli

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

//KeyOpts :
// type KeyOpts struct {
// 	PipelineDir  string
// 	RegistryPath string
// 	ImageTag     string
// 	PipelineOpt  string
// 	KeyRef       string
// }

//Gha :
func Gha() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("gha", flag.ExitOnError)
	)
	return &ffcli.Command{
		Name:        "gha",
		ShortUsage:  "tapestry gha sign|show|verify",
		ShortHelp:   `manage all github action pipeline resources`,
		LongHelp:    `manage all github action pipeline resources`,
		FlagSet:     flagset,
		Subcommands: []*ffcli.Command{
			// Pipeline Options
		},
		Exec: func(ctx context.Context, args []string) error {
			return errors.New("github action not yet supported")
		},
	}
}
