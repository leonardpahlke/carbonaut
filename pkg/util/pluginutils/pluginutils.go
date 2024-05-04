package pluginutils

import (
	"flag"
	"fmt"

	"github.com/gookit/validate"
)

type PluginArgs struct {
	Identifier string `validate:"required"`
	Port       int    `validate:"required"`
}

// Retrieves the plugin args specified as flags
func GetPluginArgs() *PluginArgs {
	// parse flags
	var pluginIdentifier string
	var port int
	flag.StringVar(&pluginIdentifier, "i", "", "set a unique carbonaut identifier")
	flag.IntVar(&port, "p", 0, "set the carbonaut plugin port")
	flag.Parse()
	args := PluginArgs{
		Identifier: pluginIdentifier,
		Port:       port,
	}
	if err := args.Validate(); err != nil {
		panic(err)
	}
	return &args
}

// SetPluginArgs sets the plugin args as flags
func SetPluginArgs(args PluginArgs) string {
	if err := args.Validate(); err != nil {
		panic(err)
	}
	return fmt.Sprintf(" -i %q -p %d", args.Identifier, args.Port)
}

// Validate checks if the plugin args are valid
func (args PluginArgs) Validate() error {
	v := validate.Struct(args)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}
