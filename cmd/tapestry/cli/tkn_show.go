package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
	"github.com/xlab/treeprint"
	"github.ibm.scs.com/tapestry/pkg/tkn"
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
		ShortUsage: "tapestry tkn show [-r] <pipeline dir>",
		ShortHelp:  `Show all tekton pipeline resources`,
		LongHelp: `Show all tekton pipeline resources
EXAMPLES
  # Show all pipeline resources from given repository
  tapestry tkn show -d ./sample-pipeline

  # Show all resources for a pipeline from given repository
  tapestry tkn show -d ./sample-pipeline -p pr-pipeline
  `,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {
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
