package plugins

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Postbuild is the interface that we're exposing as a plugin.
type Postbuild interface {
	Process() error
}

// Here is an implementation that talks over RPC
type PostbuildRPC struct{ client *rpc.Client }

func (g *PostbuildRPC) Process() error {
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
type PostbuildRPCServer struct {
	// This is the real implementation
	Impl Postbuild
}

func (s *PostbuildRPCServer) Process(args interface{}, resp *error) error {
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
type PostbuildPlugin struct {
	// Impl Injection
	Impl Postbuild
}

func (p *PostbuildPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PostbuildRPCServer{Impl: p.Impl}, nil
}

func (PostbuildPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PostbuildRPC{client: c}, nil
}
