<!-- omit in toc -->
# yum2md

Generate Markdown table from yum or dnf check-update output

[![Latest Release](https://img.shields.io/github/release/atc0005/yum2md.svg?style=flat-square)](https://github.com/atc0005/yum2md/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/yum2md.svg)](https://pkg.go.dev/github.com/atc0005/yum2md)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/yum2md)](https://github.com/atc0005/yum2md)
[![Lint and Build](https://github.com/atc0005/yum2md/actions/workflows/lint-and-build.yml/badge.svg)](https://github.com/atc0005/yum2md/actions/workflows/lint-and-build.yml)
[![Project Analysis](https://github.com/atc0005/yum2md/actions/workflows/project-analysis.yml/badge.svg)](https://github.com/atc0005/yum2md/actions/workflows/project-analysis.yml)

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
- [Changelog](#changelog)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using release binaries](#using-release-binaries)
  - [Deployment](#deployment)
- [Configuration](#configuration)
  - [Command-line arguments](#command-line-arguments)
- [Examples](#examples)
- [License](#license)
- [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Overview

This repo is intended to provide various tools used to monitor memory usage.

| Tool Name | Overall Status | Description                                                                                                                        |
| --------- | -------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `yum2md`  | Beta           | CLI tool to convert `yum check-update` / `dnf check-update` output (RHEL 8+) from stdin into a Markdown-formatted table on stdout. |

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for the *target* version this project
    was developed against
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Red Hat Enterprise Linux 8+

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/yum2md`
   1. `cd yum2md`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - manually, explicitly specifying target OS and architecture
     - `GOOS=linux GOARCH=amd64 go build -mod=vendor ./cmd/yum2md/`
       - most likely this is what you want (if building manually)
       - substitute `amd64` with the appropriate architecture if using
         different hardware (e.g., `arm64` or `386`)
   - using Makefile `linux` recipe
     - `make linux`
       - generates x86 and x64 binaries
   - using Makefile `release-build` recipe
     - `make release-build`
       - generates the same release assets as provided by this project's
         releases
1. Locate generated binaries
   - if using `Makefile`
     - look in `/tmp/yum2md/release_assets/yum2md/`
   - if using `go build`
     - look in `/tmp/yum2md/`
1. Copy the applicable binaries to whatever systems needs to run them so that
   they can be deployed

**NOTE**: Depending on which `Makefile` recipe you use the generated binary
may be compressed and have an `xz` extension. If so, you should decompress the
binary first before deploying it (e.g., `xz -d yum2md-linux-amd64.xz`).

### Using release binaries

1. Download the [latest release][repo-url] binaries
1. Decompress binaries
   - e.g., `xz -d yum2md-linux-amd64.xz`
1. Copy the applicable binaries to whatever systems needs to run them so that
   they can be deployed

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

### Deployment

1. Place `yum2md` in a location where it can be easily accessed
   - Usually the same place as other custom tools installed outside of your
     package manager's control
   - e.g., `/usr/local/bin/yum2md`
   - This binary does not require elevated privileges

**NOTE**:

DEB and RPM packages are provided as an alternative to manually deploying
binaries.

## Configuration

### Command-line arguments

None. Pipe output from `yum check-update` or `dnf check-update` into `yum2md`
to create a Markdown-formatted table.

## Examples

Output below from various RHEL test WSL instances running on a Windows 11
system. The OS release is listed in the hostname.

```console
[root@atc0005-wsl-rhel8-test yum2md-prototype]# yum check-update | ./yum2md
| Package           | Release                       | Repo              |
| ----------------- | ----------------------------- | ----------------- |
| libnghttp2.x86_64 | rhel-8-for-x86_64-baseos-rpms | 1.33.0-6.el8_10.2 |
```

```console
[root@atc0005-wsl-rhel9-test yum2md-prototype]# yum check-update | ./yum2md
| Package           | Release                       | Repo             |
| ----------------- | ----------------------------- | ---------------- |
| libnghttp2.x86_64 | rhel-9-for-x86_64-baseos-rpms | 1.43.0-6.el9_7.1 |
```

```console
[root@atc0005-wsl-rhel10-test yum2md-prototype]# yum check-update | ./yum2md
| Package            | Release                        | Repo                 |
| ------------------ | ------------------------------ | -------------------- |
| libnghttp2.x86_64  | rhel-10-for-x86_64-baseos-rpms | 1.64.0-2.el10_1.1    |
| vim-data.noarch    | rhel-10-for-x86_64-baseos-rpms | 2:9.1.083-6.el10_1.3 |
| vim-minimal.x86_64 | rhel-10-for-x86_64-baseos-rpms | 2:9.1.083-6.el10_1.3 |
```

## License

See the [LICENSE](LICENSE) file for details.

## References

- <https://linux.die.net/man/8/yum>
- <https://stackoverflow.com/questions/23698638/how-to-get-just-a-list-of-yum-updates>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/yum2md>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"
