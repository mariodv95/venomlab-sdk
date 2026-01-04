package example_plugin

import (
	"github.com/hashicorp/go-plugin"
	sdk "github.com/mariodv95/venomlab-sdk"
)

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: sdk.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"venomlab_plugin": &sdk.VenomlabPlugin{
				Impl: NewMyAwesomeNode(),
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
