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

package main

import (
	"flag"

	"log/slog"

	"carbonaut.cloud/internal/core"
)

var configFullPath string

func init() {
	flag.StringVar(&configFullPath, "c", "config.yaml", "Full path of the Carbonaut configuration file")
	flag.Parse()
}

func main() {
	c, err := core.New(&core.ConfigCore{
		ConfigurationFilePath: configFullPath,
	})
	if err != nil {
		slog.Error("error while creating core", slog.String("error", err.Error()))
		return
	}
	if err := c.Run(); err != nil {
		slog.Error("error while running core", slog.String("error", err.Error()))
		return
	}
}
