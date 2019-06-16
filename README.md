# syso

[![godoc]](https://godoc.org/github.com/hallazzang/syso)
[![goreportcard]](https://goreportcard.com/report/github.com/hallazzang/syso)

A tool for generating .syso file(COFF object file) that can be consumed
by golang build toolchain to embed binary data into an executable.
See [golang wiki](https://github.com/golang/go/wiki/GcToolchainTricks) for details.

# Install

```
$ go get -u github.com/hallazzang/syso/...
```

# Usage

```
$ syso -ico icon.ico
```

This command will generate `out.syso` file into current working directory.
You can then `go build` to include the icon resource in your executable.

[godoc]: https://godoc.org/github.com/hallazzang/syso?status.svg
[goreportcard]: https://goreportcard.com/badge/github.com/hallazzang/syso
