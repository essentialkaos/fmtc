<p align="center"><a href="#readme"><img src="https://gh.kaos.st/fmtc.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/fmtc"><img src="https://kaos.sh/r/fmtc.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/fmtc/ci"><img src="https://kaos.sh/w/fmtc/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/fmtc/codeql"><img src="https://kaos.sh/w/fmtc/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`fmtc` is a simple utility for rendering [fmtc](https://github.com/essentialkaos/ek/tree/master/fmtc#readme) formatted data. You can use it instead of the `echo` command to print colored messages.

```bash
# Simple formatted message
fmtc "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"

# Use option -E/--error to print message to stderr
fmtc -E "{r*}There is no user bob{!}"

# You can use stdin as a source of data
echo "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}" | fmtc
# or
fmtc <<< "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"

# You can print message without colors using -nc/--no-color option
fmtc -nc "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"

# Also, fmtc supports the NO_COLOR environment variable (https://no-color.org)
NO_COLOR=1 fmtc "{*}Done!{!} File {#87}$file{!} successfully uploaded to {g_}$host{!}"
```

### Installation

#### From source

To build the `fmtc` from scratch, make sure you have a working Go 1.20+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/fmtc@latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/fmtc/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) fmtc
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo fmtc --completion=bash 1> /etc/bash_completion.d/fmtc
```

ZSH:
```bash
sudo fmtc --completion=zsh 1> /usr/share/zsh/site-functions/fmtc
```

Fish:
```bash
sudo fmtc --completion=fish 1> /usr/share/fish/vendor_completions.d/fmtc.fish
```

### Man documentation

You can generate man page using next command:

```bash
fmtc --generate-man | sudo gzip > /usr/share/man/man1/fmtc.1.gz
```

### Usage

```

```

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/fmtc/ci.svg?branch=master)](https://kaos.sh/w/fmtc/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/fmtc/ci.svg?branch=develop)](https://kaos.sh/w/fmtc/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
