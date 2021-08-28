package content

import (
	"log"
	"os"
	"os/exec"

	"github.com/cueblox/blox/plugins"
	"github.com/cueblox/blox/plugins/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/pterm/pterm"
)

func (s *Service) runPrePlugins() error {

	pre, err := s.Cfg.GetList("prebuild")
	if err != nil {
		return err
	}
	iter, err := pre.List()
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
		exec, err := val.FieldByName("executable", false)
		if err != nil {
			return err
		}
		e, err := exec.Value.String()
		if err != nil {
			return err
		}
		prePluginMap[n] = &plugins.PrebuildPlugin{}
		pterm.Info.Println("Calling Plugin", n, e)
		err = s.callPlugin(n, e)
		if err != nil {
			return err
		}

	}
	return nil

}
func (s *Service) callPlugin(name, executable string) error {
	pterm.Info.Println("calling the plugin")
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.PrebuildHandshakeConfig,
		Plugins:         prePluginMap,
		Cmd:             exec.Command(executable),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(name)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	imgs := raw.(plugins.Prebuild)
	return imgs.Process(s.rawConfig)
}
