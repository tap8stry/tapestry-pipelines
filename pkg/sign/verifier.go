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
	"io/ioutil"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	"github.com/sigstore/cosign/cmd/cosign/cli"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/k8s-manifest-sigstore/pkg/k8smanifest"
	k8ssigutil "github.com/sigstore/k8s-manifest-sigstore/pkg/util"
	"github.com/tapestry-pipelines/pkg/common"
)

//Verify :
func (candidates *TknSignCandidates) Verify(ctx context.Context, verifyOpts common.KeyVerifyOpts) error {
	for idx, pc := range candidates.PipelinesSC {
		imgRef := fmt.Sprintf("%s/%s:%s", verifyOpts.RegistryPath, pc.Name, verifyOpts.ImageTag)
		if verified, err := verifyYamlFile(imgRef, verifyOpts.KeyRef, pc.Filepath); err != nil || !verified {
			// fmt.Fprintln(os.Stderr, err.Error())
			candidates.PipelinesSC[idx].Verified = false
		} else {
			candidates.PipelinesSC[idx].Verified = true
		}
		for idy, tc := range pc.TaskRefs {
			imgRef := fmt.Sprintf("%s/%s:%s", verifyOpts.RegistryPath, tc.Name, verifyOpts.ImageTag)
			if verified, err := verifyYamlFile(imgRef, verifyOpts.KeyRef, tc.Filepath); err != nil || !verified {
				candidates.PipelinesSC[idx].TaskRefs[idy].Verified = false
				// fmt.Fprintln(os.Stderr, err.Error())
			} else {
				candidates.PipelinesSC[idx].TaskRefs[idy].Verified = true
			}
			for idz, sc := range tc.Steps {
				if verified, err := verifyTaskImage(ctx, sc.ImageRef, verifyOpts.KeyRef); err != nil || !verified {
					// fmt.Fprintln(os.Stderr, err.Error())
					candidates.PipelinesSC[idx].TaskRefs[idy].Steps[idz].Verified = false
				} else {
					candidates.PipelinesSC[idx].TaskRefs[idy].Steps[idz].Verified = true
				}
			}
		}
	}
	return nil
}

func verifyYamlFile(imgRef, pubkeyRef, filepath string) (bool, error) {
	var verified bool
	manifest, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return verified, errors.Wrapf(err, "error opening manifest file: %s", filepath)
	}

	annotations := k8ssigutil.GetAnnotationsInYAML(manifest)
	annoImageRef, annoImageRefFound := annotations[k8smanifest.ImageRefAnnotationKey]
	if imgRef == "" && annoImageRefFound {
		imgRef = annoImageRef
	}

	vo := &k8smanifest.VerifyManifestOption{}
	if imgRef != "" {
		vo.ImageRef = imgRef
	}
	if pubkeyRef != "" {
		vo.KeyPath = pubkeyRef
	}
	result, err := k8smanifest.VerifyManifest(manifest, vo)
	if err != nil {
		// fmt.Fprintln(os.Stderr, err.Error())
		return verified, errors.Wrapf(err, "error verifying manifest file: %s", filepath)
	}
	verified = result.Verified
	return verified, nil
}

func verifyTaskImage(ctx context.Context, imgRef, pubkeyRef string) (bool, error) {
	ref, err := name.ParseReference(imgRef)
	if err != nil {
		return false, fmt.Errorf("failed to parse image ref `%s`; %s", imgRef, err.Error())
	}
	co := &cosign.CheckOpts{
		RegistryClientOpts: []remote.Option{
			remote.WithAuthFromKeychain(authn.DefaultKeychain),
			remote.WithContext(ctx),
		},
	}

	pubkeyVerifier, err := cli.LoadPublicKey(context.Background(), pubkeyRef)
	if err != nil {
		return false, fmt.Errorf("error loading public key; %s", err.Error())
	}
	co.SigVerifier = pubkeyVerifier
	verified, err := cosign.Verify(context.Background(), ref, co)
	if err != nil {
		return false, fmt.Errorf("error occured while verifying image `%s`; %s", imgRef, err.Error())
	}
	if len(verified) == 0 {
		return false, fmt.Errorf("no verified signatures in the image `%s`; %s", imgRef, err.Error())
	}

	return true, nil
}
