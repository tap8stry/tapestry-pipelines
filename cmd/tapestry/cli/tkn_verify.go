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
	"github.com/xlab/treeprint"
	"github.ibm.scs.com/tapestry/pkg/common"
	"github.ibm.scs.com/tapestry/pkg/tkn"
)

//TknVerify :
func TknVerify() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("verify", flag.ExitOnError)
		key     = flagset.String("key", "", "path to the private key file, KMS URI or Kubernetes Secret")
		// force       = flagset.Bool("f", false, "skip warnings and confirmations")
		recursive   = flagset.Bool("r", false, "scan all pipeline resources recusively")
		imgRegistry = flagset.String("i", "", "oci image registry path")
		imgTag      = flagset.String("t", "", "oci image path to use")
		piplineDir  = flagset.String("d", "", "pipeline directory")
	)
	return &ffcli.Command{
		Name:       "verify",
		ShortUsage: "tapestry tkn verify -key <key path> [-r] <pipeline dir> [i] <oci registry path> [t] <image tag>",
		ShortHelp:  `verify all tekton pipeline resources`,
		LongHelp: `Verify all tekton pipeline resources
EXAMPLES
  # verify all pipeline resources
  tapestry tkn verify -k ./cosign.pub -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1
  # verify resources for a give pipeline
  tapestry tkn verify -k ./cosign.pub -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1 -p pr-pipeline
  `,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {
			ko := common.KeyVerifyOpts{
				KeyRef:       *key,
				PipelineDir:  *piplineDir,
				RegistryPath: *imgRegistry,
				ImageTag:     *imgTag,
			}

			if err := VerifyPipeline(ctx, ko, *recursive); err != nil {
				return errors.Wrapf(err, "pipeline verification %s", ko.PipelineDir)
			}
			// for _, img := range args {
			// 	if err := SignCmd(ctx, ko, annotations.annotations, img, *cert, *upload, *payloadPath, *force, *recursive); err != nil {
			// 		return errors.Wrapf(err, "signing %s", img)
			// 	}
			// }
			return nil
		},
	}
}

var (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

//VerifyPipeline :
func VerifyPipeline(ctx context.Context, ko common.KeyVerifyOpts, recursive bool) error {
	signCandidates, err := tkn.GenSignCandidates(ctx, ko.PipelineDir, "")
	if err != nil {
		return errors.Wrapf(err, "pipeline signing %s", ko.PipelineDir)
	}
	if err := signCandidates.Verify(ctx, ko); err != nil {
		return errors.Wrapf(err, "pipeline signing %s", ko.PipelineDir)
	}
	tree := treeprint.NewWithRoot(ko.PipelineDir)
	for _, pc := range signCandidates.PipelinesSC {
		pipelineResult := ""
		if pc.Verified {
			pipelineResult = fmt.Sprintf("%s %s (pipeline)", string(colorGreen), pc.Name)
		} else {
			pipelineResult = fmt.Sprintf("%s %s (pipeline)", string(colorRed), pc.Name)
		}
		p := tree.AddBranch(pipelineResult)
		for _, tc := range pc.TaskRefs {
			taskResult := ""
			if tc.Verified {
				taskResult = fmt.Sprintf("%s %s (pipeline)", string(colorGreen), tc.Name)
			} else {
				taskResult = fmt.Sprintf("%s %s (pipeline)", string(colorRed), tc.Name)
			}
			t := p.AddBranch(taskResult)
			for _, sc := range tc.Steps {
				s := t.AddBranch(fmt.Sprintf("%s (step)", sc.Name))
				imgResult := ""
				if tc.Verified {
					imgResult = fmt.Sprintf("%s %s (pipeline)", string(colorGreen), sc.ImageRef)
				} else {
					imgResult = fmt.Sprintf("%s %s (pipeline)", string(colorRed), sc.ImageRef)
				}
				s.AddNode(imgResult)
			}
		}
	}
	fmt.Println(tree.String())
	return nil
}
