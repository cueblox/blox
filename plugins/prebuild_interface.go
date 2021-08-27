package plugins

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Prebuild is the interface that we're exposing as a plugin.
type Prebuild interface {
	Process() error
}

// Here is an implementation that talks over RPC
type PrebuildRPC struct{ client *rpc.Client }

func (g *PrebuildRPC) Process() error {
	var resp error
	err := g.client.Call("Plugin.Process", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type PrebuildRPCServer struct {
	// This is the real implementation
	Impl Prebuild
}

func (s *PrebuildRPCServer) Process(args interface{}, resp *error) error {
	*resp = s.Impl.Process()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type PrebuildPlugin struct {
	// Impl Injection
	Impl Prebuild
}

func (p *PrebuildPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PrebuildRPCServer{Impl: p.Impl}, nil
}

func (PrebuildPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PrebuildRPC{client: c}, nil
}
