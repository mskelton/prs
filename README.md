# prs

Better version of gh pr list

## Installation

You can install `prs` by running the install script which will download
the [latest release](https://github.com/mskelton/prs/releases/latest).

```bash
curl -LSfs https://go.mskelton.dev/prs/install | sh
```

Or you can build from source.

```bash
git clone git@github.com:mskelton/prs.git
cd prs
go install .
```

## Usage

All flags passed to `prs` will be passed through to the underlying `gh pr list`
command.

```bash
prs
prs -R mskelton/prs
```
