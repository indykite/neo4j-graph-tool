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

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/indykite/neo4j-graph-tool-core/config"
	"github.com/indykite/neo4j-graph-tool-core/migrator"
	"github.com/spf13/cobra"
)

var (
	planTargetVersion = new(migrator.TargetVersion)
	initialBatch      string
	forceRunApplyFlag bool
)

func preparePlan(ctx context.Context) (*config.Config, *migrator.ExecutionSteps) {
	cfg := loadPlannerConfig()
	version := queryVersion(ctx, cfg)

	p, err := migrator.NewPlanner(cfg)
	if err != nil {
		er(err)
	}

	scanner, err := p.NewScanner(cfg.Planner.BaseFolder)
	if err != nil {
		er(err)
	}

	vf, err := scanner.ScanFolders()
	if err != nil {
		er(err)
	}
	execSteps := new(migrator.ExecutionSteps)

	err = p.Plan(vf, version, planTargetVersion, migrator.Batch(initialBatch), p.CreateBuilder(execSteps, true))
	if err != nil {
		er(err)
	}

	return cfg, execSteps
}

func executePlan(cfg *config.Config, execSteps *migrator.ExecutionSteps, forceRun bool) {
	if len(*execSteps) == 0 {
		fmt.Printf("%sThe plan is empty%s\n", colorCyan, colorReset)
		return
	}
	fmt.Println("Plan:")
	fmt.Print(execSteps.String())

	if !forceRun {
		fmt.Printf(
			"\n%sAre you sure you want to perform described cyphers?%s"+
				"\n    Only 'yes' will be accepted to approve."+
				"\n\n    %sEnter a value%s: ",
			colorRed, colorReset, colorYellow, colorReset)

		var confirmation string
		if _, err := fmt.Scanln(&confirmation); err != nil {
			er(err)
		}
		if confirmation != "yes" {
			fmt.Printf("\n%sI will not perform any actions. Good bye%s\n", colorGreen, colorReset)
			return
		}
	}

	// This is required by cypher-shell
	if os.Getenv("NEO4J_USERNAME") == "" {
		_ = os.Setenv("NEO4J_USERNAME", getUsername())
	}
	if os.Getenv("NEO4J_PASSWORD") == "" {
		_ = os.Setenv("NEO4J_PASSWORD", getPassword())
	}

	for _, step := range *execSteps {
		var cmd *exec.Cmd
		if step.IsCypher() {
			// #nosec G204
			cmd = exec.Command(
				"/app/cypher-shell/bin/cypher-shell",
				"-a", getAddress(),
				"--format", cfg.Planner.CypherShellFormat)
			cmd.Stdin = step.Cypher()
		} else {
			toExec := step.Command()
			if toExec[0] == "exit" {
				continue
			}
			// TODO check loops

			fmt.Println(">>> ", strings.Join(toExec, " "))
			// #nosec G204
			cmd = exec.Command(toExec[0], append(toExec[1:], "-a", getAddress())...)
		}
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			er(err)
		}
	}

	fmt.Printf("\n%sAll finished. Good bye%s\n", colorGreen, colorReset)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply changes Neo4j database according to cypher ",
	Long: `Generates and execute an execution plan inside cypher-shell.

  By default, apply scans the current directory for the cypher scripts
  and applies the changes appropriately. However, a path to another
  configuration or an execution plan can be provided. Execution plans can be
  used to only execute a pre-determined set of actions`,
	Run: func(cmd *cobra.Command, _ []string) {
		cfg, steps := preparePlan(cmd.Context())
		executePlan(cfg, steps, forceRunApplyFlag)
	},
}

// planCmd represents the plan command.
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generates and prints out an execution plan for cypher-shell.",
	Long: `Generates and prints out an execution plan for cypher-shell.

  This execution plan can be reviewed prior to running apply to get a
  sense for what cypher-shell will do. Optionally, the plan can be saved to
  a file, and apply can take this plan file to execute this plan exactly.`,
	Run: func(cmd *cobra.Command, _ []string) {
		_, steps := preparePlan(cmd.Context())
		fmt.Print(steps)
	},
}

func init() {
	if v, has := os.LookupEnv("NEO4J_TARGET_VERSION"); has && v != "" {
		// VarP uses default value from String() method. Set the environment as default
		if err := planTargetVersion.Set(v); err != nil {
			er(err)
		}
	}

	for _, c := range []*cobra.Command{applyCmd, planCmd} {
		rootCmd.AddCommand(c)
		applyConnectionFlags(c)

		c.Flags().VarP(planTargetVersion, "target", "t",
			"Target version to plan for (can be set with NEO4J_TARGET_VERSION env variable)")
		c.Flags().StringVarP(&initialBatch, "batch", "b", getEnvWithDefault("NEO4J_INITIAL_BATCH", "schema"),
			"Batch name to load based on configuration, (can be set with NEO4J_INITIAL_BATCH env variable)")
	}

	applyCmd.Flags().BoolVar(&forceRunApplyFlag, "force-run", false, "Skip confirmation before executing plan")
}
