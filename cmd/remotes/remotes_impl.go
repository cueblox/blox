package main

import (
	"os"

	"github.com/cueblox/blox/plugins"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type RemoteManager struct {
	logger hclog.Logger
}

func (g *RemoteManager) Process() error {
	g.logger.Debug("message from ImageScanner.Process")
	return g.manageRemotes()
}

func (g *RemoteManager) manageRemotes() error {

	remotes, err := s.Cfg.GetList("remotes")
	if err != nil {
		return err
	}
	iter, err := remotes.List()
	if err != nil {
		return err
	}
	for iter.Next() {
		val := iter.Value()

		//nolint
		name, err := val.FieldByName("name", false)
		if err != nil {
			return err
		}
		n, err := name.Value.String()
		if err != nil {
			return err
		}
		//nolint
		version, err := val.FieldByName("version", false)
		if err != nil {
			return err
		}
		v, err := version.Value.String()
		if err != nil {
			return err
		}
		//nolint
		repository, err := val.FieldByName("repository", false)
		if err != nil {
			return err
		}
		r, err := repository.Value.String()
		if err != nil {
			return err
		}
		err = s.ensureRemote(n, v, r)
		if err != nil {
			return err
		}
	}
	return nil

}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BLOX_PLUGIN",
	MagicCookieValue: "remotes",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Info,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	remoteManager := &RemoteManager{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"remotes": &plugins.PrebuildPlugin{Impl: remoteManager},
	}

	logger.Info("initializing plugin", "name", "remotes_impl")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
