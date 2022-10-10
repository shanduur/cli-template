# Installation

## Prerequisites

Ensure all tools are installed. The necessary tools are:

- git (`2.37.0` or newer)
- go (`1.19`)
- make (`3.81` or newer)
- dasel (`1.27.1` or newer)

On MacOS you can install them using brew:

```
brew install git go make dasel
```

On Linux, depending on distribution install *git*, *make* and *go*. Dasel can be installed from source:

```
go install github.com/tomwright/dasel/cmd/dasel@latest
```

## Installation

To build the app from source, run `make install`, it will install the app in `/usr/local/bin`.
