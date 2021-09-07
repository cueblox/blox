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
		pterm.Info.Println("Registering prebuild plugin", n)
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
		err = s.callPrePlugin(n, e)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *Service) callPrePlugin(name, executable string) error {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	executablePath, err := exec.LookPath(executable)
	if err != nil {
		pterm.Error.Println("plugin not found in path", executable)
		log.Fatal(err)
	}
	pterm.Info.Println("found plugin at path", executable, executablePath)
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.PrebuildHandshakeConfig,
		Plugins:         prePluginMap,
		Cmd:             exec.Command(executablePath),
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

	// We should have a Prebuild plugin now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	preplug := raw.(plugins.Prebuild)
	return preplug.Process(s.rawConfig)
}
