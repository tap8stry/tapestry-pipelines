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
	imageAnnotations := map[string]interface{}{}
	certPathStr := ""

	opt := cosigncli.KeyOpts{
		// Annotations: imageAnnotations,
		Sk:      sk,
		IDToken: idToken,
	}

	if keyRef != "" {
		opt.KeyRef = keyRef
		opt.PassFunc = cosigncli.GetPass
	}

	return cosigncli.SignCmd(context.Background(), opt, imageAnnotations, imgRef, certPathStr, true, "", false, false)
}
