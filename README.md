# gogrip
![Coverage](https://img.shields.io/badge/Coverage-85.1%25-brightgreen)

A cli tool similar to grep but it show the entired block of lines around the matched line.

```bash
gogrip "func"  gogrip.go
```
[![asciicast](https://asciinema.org/a/tnsIbLPJfeJwr3mfCtaTDze4q.svg)](https://asciinema.org/a/tnsIbLPJfeJwr3mfCtaTDze4q)

## Installation

first you need in your .bashrc GOPATH as part of your PATH
```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Then you can install it with:
```bash
go install .
```