# hlinspect

Half-Life mod to faciliate Half-Life physics and NPC AI research.

## Features

- Easy and fast set up. No CMake or bespoke toolchains to mess with.
- Ability to hook custom compiled `hl.dll` and `client.dll` using debug information in PDB files and allow Visual Studio debugger to attach. Custom compiled libraries tend to break standard hooking mods due to the reliance on function signatures and symbol names.
- Writing features in the higher level and simpler language of Go.
- Ability to write logs to debug output viewable with Microsoft's [DebugView](https://docs.microsoft.com/en-us/sysinternals/downloads/debugview). No installation of custom or bespoke viewing tools.

## Missing

- Only supports Windows
- Speedrunning features

## Building

### Setup

Set up WSL in Windows or use a Linux system.

Install a relatively new version of Go.

Install `gcc-mingw-w64` to get `i686-w64-mingw32-gcc` or similar.

Install `protobuf-compiler` or similar to get `protoc`.

Run the following to get `protoc-gen-go`:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

Make sure `protoc-gen-go` can be found in `PATH`, which is needed by `protoc`. If not, set

```bash
export PATH="$PATH:"$(go env GOPATH)/bin
```

### Build

Suggested environmental variables in WSL:

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
