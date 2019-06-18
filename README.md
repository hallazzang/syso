# syso

[![godoc]](https://godoc.org/github.com/hallazzang/syso)
[![goreportcard]](https://goreportcard.com/report/github.com/hallazzang/syso)

**syso** - tool for embedding various resources in go executable

## Features

- [x] Embed resources by **fixed** integer id
- [x] Embed resources by **fixed** string name

## Supported Resources

- [x] Icons
- [x] Manifest
- [ ] Version info

## Install

```
$ go get -u github.com/hallazzang/syso/...
```

## Usage

Write a configuration file in JSON, which tells syso what resources you want to embed.
Here's an example:

```json
{
  "icons": [
    {
      "id": 1,
      "path": "icon.ico"
    }
  ],
  "manifest": {
    "id": 2,
    "path": "App.exe.manifest"
  }
}
```

You can specify `name` instead of `id`:

```json
...
    {
      "name": "MyIcon",
      "path": "icon.ico"
    }
...
```

Save it as `syso.json` in project's directory and run the tool:

```
$ syso
```

This will generate `out.syso` in your current directory.
You can now `go build` to actually include the resources in your executable.

## License

MIT

[godoc]: https://godoc.org/github.com/hallazzang/syso?status.svg
[goreportcard]: https://goreportcard.com/badge/github.com/hallazzang/syso
