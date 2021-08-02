//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
