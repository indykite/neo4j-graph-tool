// Copyright (c) 2023 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains entrypoint for supervisor.
package main

import (
	"fmt"
	"os"

	"github.com/indykite/neo4j-graph-tool-core/config"
	"github.com/indykite/neo4j-graph-tool-core/supervisor"
)

func main() {
	cfg, err := config.LoadFile("/app/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config:\n%s\n", err.Error())
		os.Exit(1)
	}
	err = supervisor.Start(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start supervisor:\n%s\n", err.Error())
		os.Exit(2)
	}
}
