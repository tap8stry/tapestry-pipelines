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
	"github.com/tapestry-pipelines/pkg/common"
	"github.com/tapestry-pipelines/pkg/tkn"
)

var (
	privKeyPwd = ""
)

//TknSign :
func TknSign() *ffcli.Command {
	var (
		flagset     = flag.NewFlagSet("sign", flag.ExitOnError)
		key         = flagset.String("key", "", "path to the private key file, KMS URI or Kubernetes Secret")
		force       = flagset.Bool("f", false, "skip warnings and confirmations")
		recursive   = flagset.Bool("r", false, "scan all pipeline resources recusively")
		imgRegistry = flagset.String("i", "", "oci image registry path")
		imgTag      = flagset.String("t", "", "oci image path to use")
		piplineDir  = flagset.String("d", "", "pipeline directory")
	)
	return &ffcli.Command{
		Name:       "sign",
		ShortUsage: "tapestry-pipelines tkn sign -key <key path> [-f] [-r] <pipeline dir> [i] <oci registry path> [t] <image tag>",
		ShortHelp:  `Sign all tekton pipeline resources`,
		LongHelp: `Sign all tekton pipeline resources
EXAMPLES
  # sign all pipeline resources
  tapestry-pipelines tkn sign -k ./cosign.key -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1
  # sign resources for a give pipeline
  tapestry-pipelines tkn sign -k ./cosign.key -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1 -p pr-pipeline
  `,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {

			ko := common.KeySignOpts{
				KeyRef:       *key,
				PipelineDir:  *piplineDir,
				RegistryPath: *imgRegistry,
				ImageTag:     *imgTag,
			}

			if err := SignPipeline(ctx, ko, *force, *recursive); err != nil {
				return errors.Wrapf(err, "pipeline signing %s", ko.PipelineDir)
			}

			return nil
		},
	}
}

//SignPipeline :
func SignPipeline(ctx context.Context, ko common.KeySignOpts,
	force bool, recursive bool) error {
	signCandidates, err := tkn.GenSignCandidates(ctx, ko.PipelineDir, "")
	if err != nil {
		return errors.Wrapf(err, "pipeline signing %s", ko.PipelineDir)
	}
	if err := signCandidates.Sign(ctx, ko); err != nil {
		return errors.Wrapf(err, "pipeline signing %s", ko.PipelineDir)
	}
	return nil
}
