## ACK Development tools

A list of tools and binaries that will make your contributor journey easier and much more enjoyable.

### ackdev

!!! This tool is still a work in progress. Expect dragons :dragon: and fires :fire:.

`ackdev` is a command line that helps contributors to ACK manage, test and generate controllers. It also helps you manage dependencies, local repositories and github forks.


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

# Make sure to a have a Go compiler installed locally
cd dev-tools && make install
```

#### Using `go get` (after merge)


```shell
go get github.com/aws-controllers-k8s/dev-tools/cmd/ackdev
```


### Usage

Call `ackdev help` for detailed usage instructions.

#### Setup

To able to use `ackdev` you'll have to run `ackdev setup` before any other command.

For example
```bash
ackdev setup --root-directory $WORKDIR --services s3
```

The setup command will simply create a yaml file named `$HOME/.ackdev.yaml` (you can choose a different file path `--config-file`)

> If you are already a contributor you probably want to use `$GOPATH/github.com/aws-controllers-k8s` as `$WORKDIR`. 
If you only want to test the tool, 
I recommended that you to set your `$WORKDIR` to a different directory.

> Will my projects work/compile if i use a different directory than `$GOPATH/github.com/aws-controllers-k8s` ?
Yes, even if it's not inside `GOPATH`, it will work as expected. Since Go 1.11.4 the Go compiler doesn't rely **entirely** on GOPATH to build projects.

You can view your generated configuration using:

```bash
ackdev config view
```

The configuration file will look like:

``` yaml
rootDirectory: /home/amine/go/github.com/aws-controllers-k8s/dev-tools
git:
  sshKeyPath: ""
github:
  token: ""
  username: ""
  forkPrefix: ""
repositories:
  core:
  - runtime
  - dev-tools
  - community
  - code-generator
  services:
  - s3
run:
  flags: {}
```

To use all `ackdev` features you will need to fill the `git.sshKeyPath` and `github` sections.
You can do that using the `ackdev config` command,

For example to set the Github username you can run:

```bash
ackdev config set github.username a-hilaly
```

The `git.sshKeyPath` should point to the private key you use to push code to Github.
The `github.token` should contain a token that provides `fork/renaming` rights (`repo/*` policies). You can create one by following these [instructions][github-token]

[github-token]: https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token

### Dependencies

`ackdev` can help you manage dependencies and tools you will need in your ACK development journey.
To start you can run:

```bash
ackdev list dependencies # deps|dep|dependency [--short-path]
```

The output will look like:
```bash
NAME           STATUS    VERSION         PATH                     
kind           OK        0.9.0           /usr/local/bin/kind      
helm           OK        v3.2.4+g0ad800e /usr/local/bin/helm      
mockery        NOT FOUND -                                        
kubectl        OK        v1.20.0         /usr/local/bin/kubectl   
kustomize      NOT FOUND -
controller-gen OK        v0.4.0          /usr/bin/controller-gen
```

To install the missing dependencies you can run

```bash
ackdev ensure deps # dependencies||dep|dependency
```

```bash
Installing mockery into bin/mockery ... ok.
Installing kustomize from https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh ... ok.
```

### Repositories

```bash
ack list repositories # repo|repos|repository --show-url
```

The output should look like this
```bash
NAME                TYPE       BRANCH 
runtime             core         
community           core       docs
code-generator      core         
dev-tools           core        
s3-controller       controller main
sns-controller      controller bug-fix
dynamodb-controller controller
```

To fork/clone all the listed repositories run:

```bash
ackdev ensure repos
```

## Examples

### Running controllers

Set run flags config
```
ackdev config set run.flags.aws-account-id $AWS_ACCOUNT_ID
ackdev config set run.flags.aws-region $AWS_REGION
ackdev config set run.flags.enable-development-logging true
ackdev config set run.flags.log-level debug
```

Run controller
```
ackdev run s3
```

Example output
```
I0224 20:42:01.061764   55069 request.go:621] Throttling request took 1.023440638s, request: GET:https://D82FF33364465F4917BD7F1ZZA56FC092.yl4.eu-west-2.eks.amazonaws.com/apis/networking.k8s.io/v1?timeout=32s
2021-02-24T20:42:01.434+0100	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": "0.0.0.0:8080"}
2021-02-24T20:42:01.434+0100	INFO	setup	initializing service controller	{"aws.service": "s3"}
2021-02-24T20:42:01.435+0100	INFO	setup	starting manager	{"aws.service": "s3"}
2021-02-24T20:42:01.435+0100	INFO	controller-runtime.manager	starting metrics server	{"path": "/metrics"}
2021-02-24T20:42:01.435+0100	INFO	controller-runtime.controller	Starting EventSource	{"controller": "bucket", "source": "kind source: /, Kind="}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.account	created account config map	{"name": "ack-role-account-map"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "default"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "capi-system"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "capi-webhook-system"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "ack-system-test-helm"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "cert-manager"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "kube-node-lease"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "demo"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "testing"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "capa-system"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "capi-kubeadm-bootstrap-system"}
2021-02-24T20:42:01.465+0100	DEBUG	ackrt.cache.namespace	created namespace	{"name": "capi-kubeadm-control-plane-system"}
2021-02-24T20:42:01.535+0100	INFO	controller-runtime.controller	Starting Controller	{"controller": "bucket"}
2021-02-24T20:42:01.535+0100	INFO	controller-runtime.controller	Starting workers	{"controller": "bucket", "worker count": 1}
2021-02-24T20:42:01.535+0100	DEBUG	ackrt	starting reconciliation	{"kind": "Bucket", "namespace": "default", "name": "test-ack-s3-bucket", "generation": 22, "account": "111174500800", "role": "", "region": "eu-west-2"}
2021-02-24T20:42:02.151+0100	INFO	ackrt	deleted resource	{"kind": "Bucket", "namespace": "default", "name": "test-ack-s3-bucket", "generation": 22}
2021-02-24T20:42:02.184+0100	DEBUG	ackrt	removed resource from management	{"kind": "Bucket", "namespace": "default", "name": "test-ack-s3-bucket", "generation": 22}
2021-02-24T20:42:02.184+0100	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "bucket", "request": "default/test-ack-s3-bucket"}
2021-02-24T20:42:02.185+0100	DEBUG	ackrt	starting reconciliation	{"kind": "Bucket", "namespace": "testing", "name": "test-ack-s3-bucket-2", "generation": 2, "account": "118455800300", "role": "arn:aws:iam::118455800300:role/s3FullAccess", "region": "eu-west-2"}
2021-02-24T20:42:02.764+0100	INFO	ackrt	deleted resource	{"kind": "Bucket", "namespace": "testing", "name": "test-ack-s3-bucket-2", "generation": 2}
2021-02-24T20:42:02.805+0100	DEBUG	ackrt	removed resource from management	{"kind": "Bucket", "namespace": "testing", "name": "test-ack-s3-bucket-2", "generation": 2}
2021-02-24T20:42:02.805+0100	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "bucket", "request": "testing/test-ack-s3-bucket-2"}
2021-02-24T20:42:02.805+0100	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "bucket", "request": "default/test-ack-s3-bucket"}
2021-02-24T20:42:02.805+0100	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "bucket", "request": "testing/test-ack-s3-bucket-2"}
```

### Testing controllers


```bash
ackdev test unit s3
```

```bash
ackdev test unit s3
go test -v ./...
?   	github.com/aws-controllers-k8s/s3-controller/apis/v1alpha1	[no test files]
?   	github.com/aws-controllers-k8s/s3-controller/cmd/controller	[no test files]
?   	github.com/aws-controllers-k8s/s3-controller/pkg/resource	[no test files]
?   	github.com/aws-controllers-k8s/s3-controller/pkg/resource/bucket	[no test files]
?   	github.com/aws-controllers-k8s/s3-controller/pkg/resource/version	[no test files]
?   	github.com/aws-controllers-k8s/s3-controller/pkg/version	[no test files]
```

```bash
ackdev test unit dynamodb ecr s3
```

```bash
SERVICE     UNIT-TESTS 
elasticache ERROR      
ecr         ERROR      
s3          PASS       
dynamodb    PASS
```

### [Re]Generating controllers

```bash
ackdev generate controller s3 dynamodb --stream-logs
```

```bash                                                                                                                          
Building Kubernetes API objects for s3
Generating deepcopy code for s3
Generating custom resource definitions for s3
Building service controller for s3
Generating RBAC manifests for s3
Running gofmt against generated code for s3
Building Kubernetes API objects for dynamodb
Generating deepcopy code for dynamodb
Generating custom resource definitions for dynamodb
Building service controller for dynamodb
Generating RBAC manifests for dynamodb
Running gofmt against generated code for dynamodb
```

## License

This project is licensed under the Apache-2.0 License.

