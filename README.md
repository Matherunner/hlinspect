# hlinspect

## Building

Relies heavily on CGO. Suggested environmental variables:

```bash
export CC=i686-w64-mingw32-gcc
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=386
```

Run `make` to build the DLL.
