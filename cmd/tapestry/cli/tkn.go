package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

//KeyOpts :
// type KeyOpts struct {
// 	PipelineDir  string
// 	RegistryPath string
// 	ImageTag     string
// 	PipelineOpt  string
// 	KeyRef       string
// }

//Tkn :
func Tkn() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("tkn", flag.ExitOnError)
	)
	return &ffcli.Command{
		Name:       "tkn",
		ShortUsage: "tapestry tkn sign|show|verify",
		ShortHelp:  `manage all tekton pipeline resources`,
		LongHelp:   `manage all tekton pipeline resources`,
		FlagSet:    flagset,
		Subcommands: []*ffcli.Command{
			// Pipeline Options
			TknSign(),
			TknShow(),
			TknVerify()},
		Exec: func(ctx context.Context, args []string) error {
			if len(args) == 0 {
				return flag.ErrHelp
			}
			if err := flagset.Parse(args[1:]); err != nil {
				printErrAndExit(err)
			}
			return nil
		},
	}
}

func printErrAndExit(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
