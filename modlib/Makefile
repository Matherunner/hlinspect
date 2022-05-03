OUTPUT = hlinspect.dll
GOPATH = $(shell go env GOPATH)

all:
	mkdir -p internal/hlrpc/schema
	capnp compile ../capnp/hlrpc.capnp -o go:internal/hlrpc/schema --src-prefix ../capnp -I "$(GOPATH)/pkg/mod/capnproto.org/go/capnp/v3@v3.0.0-alpha.2/std/"
	go generate ./...
	go build -v -buildmode=c-shared -o $(OUTPUT)

clean:
	rm -f $(OUTPUT)
	rm -f hlinspect.h
	rm -f $(PROTOS)
	rm -f $(CAPNP)

.PHONY: all clean
