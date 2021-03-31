![units-tests-status](https://github.com/aws-controllers-k8s/dev-tools/actions/workflows/unit-tests.yaml/badge.svg)
![Repository license](https://img.shields.io/github/license/aws-controllers-k8s/dev-tools?style=flat)
![GitHub watchers](https://img.shields.io/github/watchers/aws-controllers-k8s/dev-tools?style=social)
![GitHub stars](https://img.shields.io/github/stars/aws-controllers-k8s/dev-tools?style=social)
![GitHub forks](https://img.shields.io/github/forks/aws-controllers-k8s/dev-tools?style=social)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/aws-controllers-k8s/dev-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/aws-controllers-k8s/dev-tools)](https://goreportcard.com/report/github.com/aws-controllers-k8s/dev-tools)

## ACK Development tools

A list of tools and binaries that will make your contributor journey easier and much more enjoyable.

### ackdev

> This tool is still a work in progress. Expect dragons :dragon: and fires :fire:.

`ackdev` is a command line that helps ACK contributors to manage, test and generate controllers. It also helps you manage dependencies, local repositories and github forks.


### Installation

#### By forking/cloning the repository

First fork the `aws-controllers-k8s/dev-tools` repo and rename it to `$GH_USERNAME/ack-dev-tools` then run the script below - after replacing `"A-Hilaly"` with your own GitHub username.

```bash
export GH_USERNAME="A-Hilaly"
cd $HOME/go/src/github.com/aws-controllers-k8s
git clone git@github.com:$GH_USERNAME/ack-dev-tools dev-tools
cd dev-tools
git remote add upstream git@github.com:aws-controllers-k8s/dev-tools
git fetch --all

# Make sure to a have a Go compiler >v1.9 installed locally
cd dev-tools && make install
```

#### Using `go get` (after merge)

```shell
go get github.com/aws-controllers-k8s/dev-tools/cmd/ackdev
```

### Usage

Call `ackdev help` for detailed usage instructions.

### Examples

#### Manage ackdev configuration

You can view the configuration used by ackdev by calling:

```bash
ackdev get config
```

The output will look like:

```yaml
rootDirectory: /home/amine/source/github.com/aws-controllers-k8s
git:
  sshKeyPath: /home/amine/.ssh/id_ed25519
github:
  token: somerandomtoken165489415631684131
  username: A-Hilaly
  forkPrefix: ack-
repositories:
  core:
  - runtime
  - community
  - code-generator
  - dev-tools
  - test-infra
  services:
  - s3
  - sns
  - dynamodb
  - ecr
  - elasticache
  - sagemaker
  - sqs
  - lambda
run:
  flags:
    aws-account-id: "000000000000"
    aws-region: eu-west-2
    enable-development-logging: "true"
    log-level: debug
```

To edit the configuration you can simply call:

```
ackdev edit config
```

by default this will open the configuration file using your OS default editor
which is stored in the `EDITOR` environment variable. If this variable is not
set `ackdev` will open the configuration using `vi`.

#### List dependencies

`ackdev` can help you manage dependencies and tools you will need in your ACK development journey.
To start you can run:

```bash
ackdev list deps # dep|dependency|dependencies [--short-path]
```

The output will look like:
```bash
NAME           STATUS    VERSION         PATH                     
go             OK        1.15.6          /usr/local/go/bin/go     
kind           OK        0.9.0           /usr/local/bin/kind      
helm           OK        v3.2.4+g0ad800e /usr/local/bin/helm      
mockery        NOT FOUND -                                        
kubectl        OK        v1.20.0         /usr/local/bin/kubectl   
kustomize      OK        v4.0.1          /usr/local/bin/kustomize 
controller-gen OK        v0.4.0          /usr/bin/controller-gen
```

## License

This project is licensed under the Apache-2.0 License.

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.