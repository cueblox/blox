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

func (s *Service) runPostPlugins() error {

	post, err := s.Cfg.GetList("postbuild")
	if err != nil {
		return err
	}
	iter, err := post.List()
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
		pterm.Info.Println("Registering postbuild plugin", n)

		//nolint
		exec, err := val.FieldByName("executable", false)
		if err != nil {
			return err
		}
		e, err := exec.Value.String()
		if err != nil {
			return err
		}
		prePluginMap[n] = &plugins.PostbuildPlugin{}
		pterm.Info.Println("Calling Plugin", n, e)
		err = s.callPostPlugin(n, e)
		if err != nil {
			return err
		}

	}
	return nil

}
func (s *Service) callPostPlugin(name, executable string) error {
	pterm.Info.Println("calling the plugin")
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
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.PostbuildHandshakeConfig,
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

	// We should have a Postbuild plugin now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	postplug := raw.(plugins.Postbuild)
	return postplug.Process(s.rawConfig)
}
