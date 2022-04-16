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
Coming soon

## Usage
After installing, make sure your `$GOBIN` is on your path. You can also alias `kubectl safe` as `k` or `kubectl`.

Use `kubectl safe` just like you would `kubectl`. 

```shell
alias k="kubectl safe"
alias kubectl="kubectl safe"
```

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

## Configuration
`KUBECTL_SAFE_COMMANDS` is an environment variable that can either point to a file or be a csv of kubectl commands. 

### Default Commands
Kubectl-safe by default will only prompt on write commands. You can see default set of "safe" commands 
    [here](https://github.com/rumstead/kubectl-safe/blob/c1ce432104844b460044653020b54bee7a3fc9d1/pkg/cmd/safe/types.go#L9).

### CSV example
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
