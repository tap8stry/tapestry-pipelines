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
	"github.com/tapestry-pipelines/pkg/common"
	"github.com/tapestry-pipelines/pkg/tkn"
	"github.com/xlab/treeprint"
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
		ShortUsage: "tapestry-pipelines tkn verify -key <key path> [-r] <pipeline dir> [i] <oci registry path> [t] <image tag>",
		ShortHelp:  `verify all tekton pipeline resources`,
		LongHelp: `Verify all tekton pipeline resources
EXAMPLES
  # verify all pipeline resources
  tapestry-pipelines tkn verify -k ./cosign.pub -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1
  # verify resources for a give pipeline
  tapestry-pipelines tkn verify -k ./cosign.pub -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1 -p pr-pipeline
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
				return errors.Wrapf(err, "pipeline verification failed %s", ko.PipelineDir)
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
	colorRed   = "\033[0;31m"
	colorGreen = "\033[0;32m"
	reset      = "\033[0m"
)

//VerifyPipeline :
func VerifyPipeline(ctx context.Context, ko common.KeyVerifyOpts, recursive bool) error {
	verifyStatus := true
	signCandidates, err := tkn.GenSignCandidates(ctx, ko.PipelineDir, "")
	if err != nil {
		return errors.Wrapf(err, "pipeline verification error %s", ko.PipelineDir)
	}
	if err := signCandidates.Verify(ctx, ko); err != nil {
		return errors.Wrapf(err, "pipeline verification error %s", ko.PipelineDir)
	}
	tree := treeprint.NewWithRoot(ko.PipelineDir)
	for _, pc := range signCandidates.PipelinesSC {
		pipelineResult := ""
		if pc.Verified {
			pipelineResult = fmt.Sprintf("%s %s (pipeline)", string(colorGreen), pc.Name)
		} else {
			verifyStatus = false
			pipelineResult = fmt.Sprintf("%s %s (pipeline)", string(colorRed), pc.Name)
		}
		p := tree.AddBranch(pipelineResult)
		for _, tc := range pc.TaskRefs {
			taskResult := ""
			if tc.Verified {
				taskResult = fmt.Sprintf("%s %s (task)", string(colorGreen), tc.Name)
			} else {
				verifyStatus = false
				taskResult = fmt.Sprintf("%s %s (task)", string(colorRed), tc.Name)
			}
			t := p.AddBranch(taskResult)
			for _, sc := range tc.Steps {
				s := t.AddBranch(fmt.Sprintf("%s %s (step)", string(reset), sc.Name))
				imgResult := ""
				if sc.Verified {
					imgResult = fmt.Sprintf("%s %s (image)", string(colorGreen), sc.ImageRef)
				} else {
					verifyStatus = false
					imgResult = fmt.Sprintf("%s %s (image)", string(colorRed), sc.ImageRef)
				}
				s.AddNode(imgResult)
			}
		}
	}
	fmt.Println(tree.String())
	fmt.Println(string(reset))

	if !verifyStatus { //failed
		return errors.New("pipeline resource verification error ")
	}
	return nil
}
