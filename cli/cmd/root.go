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

// Package cmd implements the CLI commands.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/indykite/neo4j-graph-tool-core/config"
	"github.com/indykite/neo4j-graph-tool-core/migrator"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/spf13/cobra"
)

type neo4jSecret struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	configPath string
	address    string
	username   string
	password   string

	neo4jSecretVal *neo4jSecret
)

const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "graph-tool",
	Short: "Tool to operate with Graph DB Neo4j",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func applyConfigFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.toml", "path to config file")
}

func applyConnectionFlags(cmd *cobra.Command) {
	applyConfigFlag(cmd)

	cmd.PersistentFlags().StringVarP(&address, "address", "a", getEnvWithDefault("NEO4J_HOST", "neo4j://localhost:7687"),
		"address and port to connect to (can be set with NEO4J_HOST env variable)")
	cmd.PersistentFlags().StringVarP(&username, "username", "u", getEnvWithDefault("NEO4J_USERNAME", "neo4j"),
		"username to connect as (can be set with NEO4J_USERNAME env variable)")
	cmd.PersistentFlags().StringVarP(&password, "password", "p", "",
		"password to connect with (can be set with NEO4J_PASSWORD env variable)")
}

func er(msg interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "%sError: %v%s\n", colorRed, msg, colorReset)
	os.Exit(1)
}

func loadPlannerConfig() *config.Config {
	c, err := config.LoadFile(configPath)
	if err != nil {
		er(err)
	}
	return c
}

func queryVersion(c *config.Config) migrator.DatabaseModel {
	d, err := newDriver()
	if err != nil {
		er(err)
	}
	defer func() { _ = d.Close() }()

	p, err := migrator.NewPlanner(c)
	if err != nil {
		er(err)
	}

	dbm, err := p.Version(d)
	if err != nil {
		er(err)
	}

	return dbm
}

func getAddress() string {
	// TODO: implement custom getter for address, like Secret Manager etc.
	return address
}

func getUsername() string {
	// TODO: implement custom getter for username, like Secret Manager etc.
	return username
}

func getPassword() string {
	// TODO: implement custom getter for password, like Secret Manager etc.
	if password != "" {
		return password
	}
	return os.Getenv("NEO4J_PASSWORD")
}

func newDriver() (neo4j.Driver, error) {
	if getUsername() == "" || getPassword() == "" {
		return nil, errors.New("missing username/password")
	}
	return neo4j.NewDriver(getAddress(), neo4j.BasicAuth(getUsername(), getPassword(), ""), func(config *neo4j.Config) {
		config.UserAgent = "neo4j-tool/1.0 (IndyKite)"
	})
}

func getEnvWithDefault(key, defaultValue string) string {
	if v, has := os.LookupEnv(key); has {
		return v
	}
	return defaultValue
}
