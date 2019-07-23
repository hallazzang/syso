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
  "Icons": [
    {
      "ID": 1,
      "Path": "icon.ico"
    }
  ],
  "Manifest": {
    "ID": 2,
    "Path": "App.exe.manifest"
  }
}
```

You can specify `name` instead of `id`:

```json
...
    {
      "Name": "MyIcon",
      "Path": "icon.ico"
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
  "Icons": [
    {
      "ID": 1,
      "Path": "icon.ico"
    },
    {
      "Name": "Icon",
      "Path": "icon2.ico"
    }
  ],
  "Manifest": {
    "ID": 1,
    "Path": "App.exe.manifest"
  },
  "VersionInfos": [
    {
      "ID": 1,
      "Fixed": {
        "FileVersion": "10.0.14393.0",
        "ProductVersion": "10.0.14393.0"
      },
      "StringTables": [
        {
          "Language": "0409",
          "Charset": "04b0",
          "Strings": {
            "CompanyName": "Microsoft Corporation",
            "FileDescription": "Windows Command Processor",
            "FileVersion": "10.0.14393.0 (rs1_release.160715-1616)",
            "InternalName": "cmd",
            "LegalCopyright": "\u00a9 Microsoft Corporation. All rights reserved.",
            "OriginalFilename": "Cmd.Exe",
            "ProductName": "Microsoft\u00ae Windows\u00ae Operating System",
            "ProductVersion": "10.0.14393.0"
          }
        }
      ],
      "Translations": [
        {
          "Language": "0409",
          "Charset": "04b0"
        }
      ]
    }
  ]
}
```

Note that keys are case-insensitive.
You can use both `"companyName"` and `"CompanyName"`, or even `"companyname"` for key.

## License

MIT

[godoc]: https://godoc.org/github.com/hallazzang/syso?status.svg
[goreportcard]: https://goreportcard.com/badge/github.com/hallazzang/syso
[rsrc]: https://github.com/akavel/rsrc
[goversioninfo]: https://github.com/josephspurrier/goversioninfo
