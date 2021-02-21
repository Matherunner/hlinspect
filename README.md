# hlinspect

Half-Life hooking mod for research.

## Features

- Fast compilation and set up
- Ability to hook custom compiled `hl.dll` and `client.dll` using debug information in PDB files

## Building

Relies heavily on CGO. Suggested environmental variables:

```bash
export CC=i686-w64-mingw32-gcc
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=386
```

Run the following to build:

```bash
cd hlinspect
go generate ./...
make
```

