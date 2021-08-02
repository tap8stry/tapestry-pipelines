# tapestry-pipelines

**Disclaimer:**
---------------

```
"The strategy is definitely: crawl, walk and run"
```
The current state of this project is in the middle of "crawl". 


## CICD Pipeline signing and verification

With growing open-source ecosystem around CICD technologies, (specially with tektoncd and Github Actions) it is becoming critical to ensure that the pipeline compositions are securely signed and verifiable. `tapestry-pipelines` is the first step in that direction to allow signing and verification of every individual pipeline resources. It performs signing/verification with [sigstore/cosign](https://github.com/sigstore/cosign) for images and [sigstore/k8s-manifest-sigstore](https://github.com/sigstore/k8s-manifest-sigstore) for yamls.

More detail discussion is recorded in this blog post: [It's time we start securing our CICD pipelines](https://nadgowdas.github.io/blog/2021/pipeline-security/)

## Installation

Currently, `tapestry-pipelines` is available as a stand-alone CLI. You can follow following procedure to install CLI:

```
git clone 
cd tapestry-pipelines 
make
```

This should create an executable binary `tapestry-pipelines`. You can add this binary to your local PATH.

## Quick start

Currently `tapestry-pipelines` supports signing/verification of tekton pipelines only. Support for Github Actions is underway. 

### Show pipeline resources
This command shows all pipeline resources from a given definition directory that are candidates for signing/verification. You can also pass an filtering option (with `-p <pipeline-name>`) to particular pipeline.

```
% tapestry-pipelines tkn show -d ./sample-pipeline
./sample-pipeline
└── pr-pipeline2 (pipeline)
    ├── git-clone-repo (task)
    │   ├── fetch-git-token (step)
    │   │   └── us.icr.io/tap8stry/pipeline-base-image:2.6
    │   └── clone-repo (step)
    │       └── us.icr.io/tap8stry/ubi-base:8.1
    ├── generate-bom (task)
    │   └── bom (step)
    │       └── us.icr.io/tap8stry/cra-bom:v1
    ├── cosign-verify-deps (task)
    │   ├── verify (step)
    │   │   └── us.icr.io/tap8stry/ssc-verifier:0.1.0
    │   ├── fetch-git-information (step)
    │   │   └── us.icr.io/tap8stry/pipeline-base-image:2.6
    │   └── comment-editor (step)
    │       └── us.icr.io/tap8stry/cra-comm-editor:main.1260
    └── git-set-commit-status (task)
        ├── fetch-git-information (step)
        │   └── us.icr.io/tap8stry/pipeline-base-image:2.6
        └── set-status (step)
            └── us.icr.io/tap8stry/ubi-base:8.1
```

### Sign pipeline resources
First, you need to generate a key-pair to be use for signing/verification. You can use `cosign` to generate it following these [Instructions](https://github.com/sigstore/cosign#generate-a-keypair). Then use your private key to sign the resources. To avoid entering password for your private key multiple times, you can set set the password as environment variable `COSIGN_PASSWORD`. You also need to provide an OCI registry `path` and `tag` use for uploading signed artifacts. 

```
% tapestry-pipelines tkn sign -key ./cosign.key -d ./sample-pipeline -i us.icr.io/tap8stry -t dev6
```

### Verify pipeline resources
You can use `verify` subcommand to verify signatures for all resources. In the `green` (pass) vs `red` (fail) color coding it reports the verification status for individual resources.

```
% tapestry-pipelines tkn verify -key ./cosign.pub -d ./sample-pipeline -i us.icr.io/tap8stry -t dev6
```

## Demo

![Demo](media/tapestry-demo.gif?)

## WIP

We have active Work In Progress (WIP) on following fronts and welcome anyone who is willing to join the forces. 

1. Support for signing/verification of Github Actions
2. Support for KMS for key handling
3. Tekton Pipeline Admission Controller that validate signatures for all pipeline resources before admitting pipeline
4. Extending sign verification to event triggers


