[![codecov](https://codecov.io/gh/rumstead/kubectl-safe/branch/main/graph/badge.svg)](https://codecov.io/gh/rumstead/kubectl-safe)
[![Go Report Card](https://goreportcard.com/badge/github.com/rumstead/kubectl-safe)](https://goreportcard.com/report/github.com/rumstead/kubectl-safe)

# kubectl-safe
A kubectl plugin to prevent shooting yourself in the foot with edit commands

## Installation
### Go 
```
go install github.com/rumstead/kubectl-safe@latest
```
### Krew
```shell
❯ kubectl krew install safe
Updated the local copy of plugin index.
Installing plugin: safe
Installed plugin: safe
\
 | Use this plugin:
 | 	kubectl safe
 | Documentation:
 | 	https://github.com/rumstead/kubectl-safe
/
WARNING: You installed plugin "safe" from the krew-index plugin repository.
   These plugins are not audited for security by the Krew maintainers.
   Run them at your own risk.
```

## Usage
After installing, make sure your `$GOBIN` is on your path. 

You can also alias `kubectl safe` as `k`, `kubectl`, or `ks`.
```shell
# You should pick one :)
alias k="kubectl safe"
alias kubectl="kubectl safe"
alias ks="kubectl safe"
```

Use `kubectl safe` just like you would `kubectl`.

```shell
$ kubectl safe get pod -n kube-system
NAME                                     READY   STATUS    RESTARTS         AGE
coredns-78fcd69978-xwdt4                 1/1     Running   10 (2d4h ago)    57d
coredns-78fcd69978-zxj4q                 1/1     Running   10 (2d4h ago)    57d
etcd-docker-desktop                      1/1     Running   10 (2d4h ago)    57d
kube-apiserver-docker-desktop            1/1     Running   10 (2d4h ago)    57d
kube-controller-manager-docker-desktop   1/1     Running   10 (2d4h ago)    57d
kube-proxy-jr2wr                         1/1     Running   10 (2d4h ago)    57d
kube-scheduler-docker-desktop            1/1     Running   13 (2d4h ago)    57d
storage-provisioner                      1/1     Running   20 (2d4h ago)    57d
vpnkit-controller                        1/1     Running   1378 (16m ago)   57d

$ kubectl safe delete pod -n kube-system coredns-78fcd69978-xwdt4
You are running a delete against context docker-desktop, continue? [yY] n
I0416 14:40:50.966746   85123 root.go:52] Not running command.

$ kubectl safe delete pod -n kube-system coredns-78fcd69978-xwdt4
You are running a delete against context docker-desktop, continue? [yY] y
pod "coredns-78fcd69978-xwdt4" deleted
```

## Shell completion
You can read more the [issue](https://github.com/rumstead/kubectl-safe/issues/17)
Add the below script anywhere in your path with the executable bit set.
```shell
#!/usr/bin/env bash

# If we are completing a flag, use Cobra's builtin completion system.
# To know if we are completing a flag we need the last argument starts with a `-` and does not contain an `=`
args=("$@")
lastArg=${args[((${#args[@]}-1))]}
if [[ "$lastArg" == -* ]]; then
   if [[ "$lastArg" != *=* ]]; then
      kubectl safe __complete "$@"
   fi
else
   kubectl __complete "$@"
fi
```

## Configuration
`KUBECTL_SAFE_COMMANDS` is an environment variable that can either point to a file or be a csv of kubectl commands.
`KUBECTL_UNSAFE_COMMANDS` is an environment variable that can either point to a file or be a csv of kubectl commands.

*NOTE*: `KUBECTL_UNSAFE_COMMANDS` takes precedence 

### Default Commands
Kubectl-safe by default will only prompt on write commands. You can see default set of "safe" commands 
    [here](https://github.com/rumstead/kubectl-safe/blob/c1ce432104844b460044653020b54bee7a3fc9d1/pkg/cmd/safe/types.go#L9).

### CSV example
#### Kube Safe Commands
```shell
$ export KUBECTL_SAFE_COMMANDS=version,config
$ kubectl safe get pod
You are running a get against context docker-desktop, continue? [yY] n
I0416 15:10:12.967439   97368 root.go:52] Not running command.
$ kubectl safe version
Client Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.5", GitCommit:"c285e781331a3785a7f436042c65c5641ce8a9e9", GitTreeState:"clean", BuildDate:"2022-03-16T15:51:05Z", GoVersion:"go1.17.8", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"22", GitVersion:"v1.22.5", GitCommit:"5c99e2ac2ff9a3c549d9ca665e7bc05a3e18f07e", GitTreeState:"clean", BuildDate:"2021-12-16T08:32:32Z", GoVersion:"go1.16.12", Compiler:"gc", Platform:"linux/amd64"}
$ kubectl safe config current-context
docker-desktop
```
#### Kube Unsafe Commands
```shell
$ export KUBECTL_UNSAFE_COMMANDS=version,get
$ kubectl safe get pod
You are running a get against context rancher-desktop, continue? [yY] y
No resources found in default namespace.
$ kubectl safe version                         
You are running a version against context rancher-desktop, continue? [yY] y
Client Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.5", GitCommit:"c285e781331a3785a7f436042c65c5641ce8a9e9", GitTreeState:"clean", BuildDate:"2022-03-16T15:58:47Z", GoVersion:"go1.17.8", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"23", GitVersion:"v1.23.6+k3s1", GitCommit:"418c3fa858b69b12b9cefbcff0526f666a6236b9", GitTreeState:"clean", BuildDate:"2022-04-28T22:16:18Z", GoVersion:"go1.17.5", Compiler:"gc", Platform:"linux/amd64"}
$ kubectl safe delete pod -n kube-system  coredns-d76bd69b-4cngl
pod "coredns-d76bd69b-4cngl" deleted
```
### File example
```shell
$ cat /tmp/valid-commands.txt
list
version
$ export KUBECTL_SAFE_COMMANDS=/tmp/valid-commands.txt
$ kubectl safe get pod                                           
I0416 15:07:54.686263   96875 commands.go:50] reading commands from /tmp/valid-commands.txt.
I0416 15:07:54.686418   96875 commands.go:55] adding list command to the safe list.
I0416 15:07:54.686423   96875 commands.go:55] adding version command to the safe list.
You are running a get against context docker-desktop, continue? [yY] n
I0416 15:07:57.124902   96875 root.go:52] Not running command.
```


## Similar plugins
1. https://github.com/kubernetes-sigs/krew-index/blob/master/plugins/prompt.yaml
