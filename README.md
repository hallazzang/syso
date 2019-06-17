# syso

[![godoc]](https://godoc.org/github.com/hallazzang/syso)
[![goreportcard]](https://goreportcard.com/report/github.com/hallazzang/syso)

A tool for generating .syso file(COFF object file) that can be consumed
by golang build toolchain to embed binary data into an executable.
See [golang wiki](https://github.com/golang/go/wiki/GcToolchainTricks) for details.

# Features

- [x] Fixed resource identifier
- [x] Embed resource by integer id
- [ ] Embed resource by name (there's a bug now)

# Install

```
$ go get -u github.com/hallazzang/syso/...
```

# Usage

First, write a configuration file which tells syso tool how to generate .syso file.
I'll embed one icon resource, named `icon.ico` using integer id `1`.

```json
{
  "icons": [
    {
      "id": 1,
      "path": "icon.ico"
    }
  ]
}
```

Save it as `syso.json`, or whatever name you want. Then run this command in same directory:

```
$ syso
```

This command will read configuration file you wrote and generate `out.syso` file in current directory.
You can also specify different config file name using `-c` option.

You can then `go build` as usual to include the icon resource in your executable.

[godoc]: https://godoc.org/github.com/hallazzang/syso?status.svg
[goreportcard]: https://goreportcard.com/badge/github.com/hallazzang/syso
