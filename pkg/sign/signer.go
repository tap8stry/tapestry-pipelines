package sign

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	cosigncli "github.com/sigstore/cosign/cmd/cosign/cli"
	"github.com/sigstore/k8s-manifest-sigstore/pkg/k8smanifest"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.ibm.scs.com/tapestry/pkg/common"
)

//TknSignCandidates :
type TknSignCandidates struct {
	PipelinesSC []PipelineSC
}

//PipelineSC :
type PipelineSC struct {
	Filepath    string
	Name        string
	PipelineObj *v1beta1.Pipeline
	Verified    bool
	TaskRefs    []TaskSC
}

//TaskSC :
type TaskSC struct {
	Filepath string
	Name     string
	TaskObj  *v1beta1.Task
	Verified bool
	Steps    []StepSC
}

//StepSC :
type StepSC struct {
	Name     string
	ImageRef string
	Verified bool
}

//Sign :
func (candidates *TknSignCandidates) Sign(ctx context.Context, signOpts common.KeySignOpts) error {
	for _, pc := range candidates.PipelinesSC {
		imgRef := fmt.Sprintf("%s/%s:%s", signOpts.RegistryPath, pc.Name, signOpts.ImageTag)
		if err := signYamlManifest(imgRef, signOpts.KeyRef, pc.Filepath); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return errors.Wrapf(err, "pipeline signing pipeline yaml %s", pc.Name)
		}
		fmt.Printf("signed manifest generated at :%s \n", fmt.Sprintf("%s.signed", pc.Filepath))
		for _, tc := range pc.TaskRefs {
			imgRef := fmt.Sprintf("%s/%s:%s", signOpts.RegistryPath, tc.Name, signOpts.ImageTag)
			if err := signYamlManifest(imgRef, signOpts.KeyRef, tc.Filepath); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return errors.Wrapf(err, "pipeline signing task yaml %s", tc.Name)
			}
			fmt.Printf("signed manifest generated at :%s \n", fmt.Sprintf("%s.signed", tc.Filepath))
			for _, sc := range tc.Steps {
				if err := signImage(sc.ImageRef, signOpts.KeyRef); err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					return errors.Wrapf(err, "pipeline signing image %s", sc.ImageRef)
				}
			}
		}
	}
	return nil
}

func signYamlManifest(imgRef, keyRef, filepath string) error {
	outputFilepath := fmt.Sprintf("%s.signed", filepath)
	so := &k8smanifest.SignOption{
		ImageRef:         imgRef,
		KeyPath:          keyRef,
		Output:           outputFilepath,
		UpdateAnnotation: true,
		// ImageAnnotations: anntns,
	}

	_, err := k8smanifest.Sign(filepath, so)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return nil
	}
	return nil
}

func signImage(imgRef, keyRef string) error {
	sk := false
	idToken := ""

	opt := cosigncli.SignOpts{
		// Annotations: imageAnnotations,
		Sk:      sk,
		IDToken: idToken,
	}

	if keyRef != "" {
		opt.KeyRef = keyRef
		opt.Pf = cosigncli.GetPass
	}

	return cosigncli.SignCmd(context.Background(), opt, imgRef, true, "", false, false)
}
