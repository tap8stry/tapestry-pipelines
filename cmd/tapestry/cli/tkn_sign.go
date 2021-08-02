package cli

import (
	"context"
	"flag"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
	"github.ibm.scs.com/tapestry/pkg/common"
	"github.ibm.scs.com/tapestry/pkg/tkn"
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
		ShortUsage: "tapestry tkn sign -key <key path> [-f] [-r] <pipeline dir> [i] <oci registry path> [t] <image tag>",
		ShortHelp:  `Sign all tekton pipeline resources`,
		LongHelp: `Sign all tekton pipeline resources
EXAMPLES
  # sign all pipeline resources
  tapestry tkn sign -k ./cosign.key -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1
  # sign resources for a give pipeline
  tapestry tkn sign -k ./cosign.key -d ./sample-pipeline-dir -i us.icr.io.tap8stry -t dev1 -p pr-pipeline
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
