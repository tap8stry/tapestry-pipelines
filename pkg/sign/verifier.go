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
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/k8s-manifest-sigstore/pkg/k8smanifest"
	"github.ibm.scs.com/tapestry/pkg/common"
)

const (
	verifyYamlCosignCmd  = "kubectl sigstore verify -k {{.Key}} -f {{.YamlFile}} -i {{.Image}}"
	verifyImageCosignCmd = "cosign verify -key {{.Key}} {{.Image}}"
)

//Verify :
func (candidates *TknSignCandidates) Verify(ctx context.Context, verifyOpts common.KeyVerifyOpts) error {
	for idx, pc := range candidates.PipelinesSC {
		imgRef := fmt.Sprintf("%s/%s:%s", verifyOpts.RegistryPath, pc.Name, verifyOpts.ImageTag)
		if verified, err := verifyYamlFile(imgRef, verifyOpts.KeyRef, pc.Filepath); err != nil || !verified {
			fmt.Fprintln(os.Stderr, err.Error())
			candidates.PipelinesSC[idx].Verified = false
		} else {
			candidates.PipelinesSC[idx].Verified = true
		}
		for idy, tc := range pc.TaskRefs {
			imgRef := fmt.Sprintf("%s/%s:%s", verifyOpts.RegistryPath, tc.Name, verifyOpts.ImageTag)
			if verified, err := verifyYamlFile(imgRef, verifyOpts.KeyRef, tc.Filepath); err != nil || !verified {
				candidates.PipelinesSC[idx].TaskRefs[idy].Verified = false
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				candidates.PipelinesSC[idx].TaskRefs[idy].Verified = true
			}
			for idz, sc := range tc.Steps {
				if verified, err := verifyTaskImage(ctx, sc.ImageRef, verifyOpts.KeyRef); err != nil || !verified {
					fmt.Fprintln(os.Stderr, err.Error())
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
	vo := &k8smanifest.VerifyManifestOption{}
	if imgRef != "" {
		vo.ImageRef = imgRef
	}
	if pubkeyRef != "" {
		vo.KeyPath = pubkeyRef
	}
	result, err := k8smanifest.VerifyManifest(manifest, vo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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
		Claims: true,
		RegistryClientOpts: []remote.Option{
			remote.WithAuthFromKeychain(authn.DefaultKeychain),
			remote.WithContext(ctx),
		},
	}

	pubkeyVerifier, err := cosign.LoadPublicKey(context.Background(), pubkeyRef)
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
	// var cert *x509.Certificate
	// for _, vp := range verified {
	// 	ss := payload.SimpleContainerImage{}
	// 	err := json.Unmarshal(vp.Payload, &ss)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	cert = vp.Cert
	// 	break
	// }
	// signerName := "" // singerName could be empty in case of key-used verification
	// if cert != nil {
	// 	signerName = k8smnfutil.GetNameInfoFromCert(cert)
	// }

	return true, nil
}
