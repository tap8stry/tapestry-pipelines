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
