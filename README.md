# Kar [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/kargar/kar) 
<p align="center">
<img width="100%" src="https://talks.golang.org/2012/waza/gophercomplex5.jpg">
</p>
Kar is a simple utility that helps you create run [Kargar](https://github.com/omeid/kargar) builds from CLI. 

To allow, co-existence of _build file_ and your Go project code and to insure it 
is properly isolated, kar is  _guarded_ with a build build constraint, `kar`;  to
properly leverage Go tooling's caching mechanism and avoid any issues with stale
cache objects, the task runner, `cmd/kar` uses it's own `pkgdir` (`$GOPATH/pkg/$GOOS_$GOARCH_kar`).


## Install
_Requires Go 1.5_.

```go
go install github.com/omeid/kar/cmd/kar
```

# Example

```sh
$ cat demo_kar.go
```

```go
// +build kar

package main

import (
  "github.com/omeid/gonzo/context"
  "github.com/omeid/kargar"
  "github.com/omeid/kargar/kar"
)

func init() {

  kar.Run(func(build *kargar.Build) error {

    return build.Add(
      kargar.Task{

        Name:  "say-hello",
        Usage: "This tasks is self-documented, it says hello for every second.",

        Action: func(ctx context.Context) error {
          ctx.Info("Hello!")
          return nil
        },
      })
  })
}

```

```sh
$ kar install    # build and  cache the dependencies.
$ kar say-hello #run task say-hello
INFO[0000] [say-hello]                                  
INFO[0000] Hello!                                        task=say-hello
```

#TODO:
 - Finish this doc.
