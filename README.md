# API

## How to run

[Install golang](https://golang.org/doc/install)

Create your workspace:

```
export GOPATH=$HOME/go
mkdir -p $GOPATH/src/github.com
cd $GOPATH/src/github.com
git clone git@github.com:pi-top/utils-go.git
cd utils-go
```

Finally run `go get -u`, this command may take a while as it's doing a `git clone` of all dependencies into `$GOPATH/src`.

You may need to install libxml2-devel or libxml2-dev

Now you can run the API via one of the following:

* `go run *.go`: compiles and runs
* `go build`: creates a binary in local dir
* `go install`: create a binary in $GOPATH/bin

## Setting log level

By default the API only logs requests. To show all log messages `export LOG=debug` or choose info, warn, error, fatal, panic.
