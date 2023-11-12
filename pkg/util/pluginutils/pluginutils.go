/*
Copyright 2023 CARBONAUT AUTHORS

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	return fmt.Sprintf(" -i \"%s\" -p %d", args.Identifier, args.Port)
}

// Validate checks if the plugin args are valid
func (args PluginArgs) Validate() error {
	v := validate.Struct(args)
	if !v.Validate() {
		return v.Errors
	}
	return nil
}
