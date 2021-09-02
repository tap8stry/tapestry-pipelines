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

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
	"github.com/tapestry-pipelines/pkg/tkn"
	"github.com/xlab/treeprint"
)

//KeyShowOpts :
type KeyShowOpts struct {
	PipelineDir  string
	PipelineName string
}

//TknShow :
func TknShow() *ffcli.Command {
	var (
		flagset    = flag.NewFlagSet("show", flag.ExitOnError)
		piplineDir = flagset.String("d", "", "pipeline directory")
		pipeline   = flagset.String("p", "", "pipeline name")
	)
	return &ffcli.Command{
		Name:       "show",
		ShortUsage: "tapestry-pipelines tkn show [-r] <pipeline dir>",
		ShortHelp:  `Show all tekton pipeline resources`,
		LongHelp: `Show all tekton pipeline resources
EXAMPLES
  # Show all pipeline resources from given repository
  tapestry-pipelines tkn show -d ./sample-pipeline

  # Show all resources for a pipeline from given repository
  tapestry-pipelines tkn show -d ./sample-pipeline -p pr-pipeline
  `,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {
			if *piplineDir == "" {
				return flag.ErrHelp
			}
			ko := KeyShowOpts{
				PipelineDir:  *piplineDir,
				PipelineName: *pipeline,
			}

			if err := ShowPipeline(ctx, ko); err != nil {
				return errors.Wrapf(err, "pipeline show %s", ko.PipelineDir)
			}
			return nil
		},
	}
}

//ShowPipeline :
func ShowPipeline(ctx context.Context, ko KeyShowOpts) error {
	signCandidates, _ := tkn.GenSignCandidates(ctx, ko.PipelineDir, ko.PipelineName)
	tree := treeprint.NewWithRoot(ko.PipelineDir)
	for _, pc := range signCandidates.PipelinesSC {
		p := tree.AddBranch(fmt.Sprintf("%s (pipeline)", pc.Name))
		for _, tc := range pc.TaskRefs {
			t := p.AddBranch(fmt.Sprintf("%s (task)", tc.Name))
			for _, sc := range tc.Steps {
				s := t.AddBranch(fmt.Sprintf("%s (step)", sc.Name))
				s.AddNode(sc.ImageRef)
			}
		}
	}
	fmt.Println(tree.String())
	return nil
}
