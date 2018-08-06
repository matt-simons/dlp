# DLP Example

## How to run

[Install golang](https://golang.org/doc/install)

Create your workspace:

```
export GOPATH=$HOME/go
mkdir -p $GOPATH/src/github.com
cd $GOPATH/src/github.com
git clone git@github.com:matt-simons/dlp.git
cd dlp
```

Finally run `go get -u`, this command may take a while as it's doing a `git clone` of all dependencies into `$GOPATH/src`.

Now you can run the example via one of the following:

* `go run *.go`: compiles and runs
* `go build`: creates a binary in local dir
* `go install`: create a binary in $GOPATH/bin
