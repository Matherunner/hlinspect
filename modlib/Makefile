OUTPUT = hlinspect.dll
GOPATH = $(shell go env GOPATH)

all: dll
	
dll: schema
	go build -v -buildmode=c-shared -o $(OUTPUT)

schema:
	mkdir -p internal/hlrpc/schema
	capnp compile ../capnp/hlrpc.capnp -o go:internal/hlrpc/schema --src-prefix ../capnp -I "$(GOPATH)/pkg/mod/capnproto.org/go/capnp/v3@v3.0.0-alpha.2/std/"
	go generate ./...

clean:
	rm -f $(OUTPUT)
	rm -f hlinspect.h
	rm -f $(PROTOS)
	rm -f $(CAPNP)

.PHONY: all dll schema clean
