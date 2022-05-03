import subprocess

go_cmd_res = subprocess.run(['go', 'env', 'GOPATH'], capture_output=True)
go_path = go_cmd_res.stdout.decode('utf-8').strip()

import os
import sys
sys.path.append('../capnp')
# FIXME: should not hardcode this?
sys.path.append(os.path.join(go_path, 'pkg/mod/capnproto.org/go/capnp/v3@v3.0.0-alpha.2/std'))

import capnp
import hlrpc_capnp

# NOTE: For this to work on WSL, must use the default route, and the server must not listen to just localhost:32002

client = capnp.TwoPartyClient('172.26.112.1:32002')
halflife = client.bootstrap().cast_as(hlrpc_capnp.HalfLife)

promise = halflife.getFullPlayerState()

response = promise.wait()
print('response', response)
