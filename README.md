# syso

[![godoc]](https://godoc.org/github.com/hallazzang/syso)
[![goreportcard]](https://goreportcard.com/report/github.com/hallazzang/syso)

**syso** - tool for embedding various resources in go executable

Table of contents:

- [Features](#Features)
- [Installation](#Installation)
- [Usage](#Usage)
- [License](#License)

## Features

| Feature                      | [rsrc] | [goversioninfo] | syso(this project) |
| :--------------------------- | :----: | :-------------: | :----------------: |
| Embedding icons              |   ✔    |        ✔        |         ✔          |
| Embedding manifest           |   ✔    |        ✔        |         ✔          |
| Configuration through a file |        |        ✔        |         ✔          |
| Embedding version info       |        |        ✔        |         ✔          |
| Fixed resource identifier    |        |                 |         ✔          |

## Installation

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

### Configuration

Here's all possible configurations:

```json
{
  "icons": [
    {
      "id": 1,
      "path": "icon.ico"
    },
    {
      "name": "Icon",
      "path": "icon2.ico"
    }
  ],
  "manifest": {
    "id": 1,
    "path": "App.exe.manifest"
  },
  "versioninfo": {
    "id": 1,
    "fixed": {
      "fileVersion": "1.2.3.4",
      "productVersion": "5.6.7.8"
    },
    "strings": {
      "comments": "Comments",
      "companyName": "My Company"
    }
  }
}
```

## License

MIT

[godoc]: https://godoc.org/github.com/hallazzang/syso?status.svg
[goreportcard]: https://goreportcard.com/badge/github.com/hallazzang/syso
[rsrc]: https://github.com/akavel/rsrc
[goversioninfo]: https://github.com/josephspurrier/goversioninfo
