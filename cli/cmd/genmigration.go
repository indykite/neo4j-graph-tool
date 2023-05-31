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
	"fmt"
	"strings"
	"time"

	"github.com/indykite/neo4j-graph-tool-core/migrator"
	"github.com/spf13/cobra"
)

// cmd arguments.
var (
	migrationUpCommand, migrationDownCommand bool
)

// genMigrationCmd represents the genmigration command.
var genMigrationCmd = &cobra.Command{
	Use:   "genmigration <folder> <version> <name> \n genmigration <version> <name>",
	Short: "Generates new migration file(s) based on provided config",
	Long: `If only 2 arguments are provided, folder will be set to 'schema'.
By default the created files are Cypher files. This can be changed with flags.
If version doesn't contain part after +, current timestamp is automatically added. Otherwise provided value is used`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		var folder, version, name string
		upFileType, downFileType := migrator.Cypher, migrator.Cypher
		if len(args) == 2 {
			folder = "schema"
			version = args[0]
			name = args[1]
		} else {
			folder = args[0]
			version = args[1]
			name = args[2]
		}

		if len(name) == 0 {
			er("migration name cannot be empty")
		}
		targetVersion, err := migrator.ParseTargetVersion(version)
		if err != nil {
			er("invalid version: " + err.Error())
		}
		if targetVersion.Revision <= 0 {
			targetVersion.Revision = time.Now().Unix()
		}
		if migrationUpCommand {
			upFileType = migrator.Command
		}
		if migrationDownCommand {
			downFileType = migrator.Command
		}

		cfg := loadPlannerConfig()
		p, err := migrator.NewPlanner(cfg)
		if err != nil {
			er(err)
		}

		scanner, err := p.NewScanner(cfg.Planner.BaseFolder)
		if err != nil {
			er(err)
		}
		var paths []string
		paths, err = scanner.GenerateMigrationFiles(folder, targetVersion, name, upFileType, downFileType)
		if err != nil {
			er(err)
		}
		fmt.Printf("Created:\n -%s \n", strings.Join(paths, "\n -"))
	},
}

func init() {
	genMigrationCmd.Flags().BoolVarP(&migrationUpCommand, "up_cmd", "u", false, "Create Up migration as Command file")
	genMigrationCmd.Flags().
		BoolVarP(&migrationDownCommand, "down_cmd", "d", false, "Create Down migration as Command file")

	tmpl := genMigrationCmd.UsageTemplate()
	tmpl = strings.ReplaceAll(tmpl, "{{.UseLine}}", "{{MyUseLine .}}\n{{.Long}}")
	cobra.AddTemplateFunc("MyUseLine", MyUseLine)
	genMigrationCmd.SetUsageTemplate(tmpl)
	genMigrationCmd.SetHelpTemplate("{{.UsageString}}")

	applyConfigFlag(genMigrationCmd)
	rootCmd.AddCommand(genMigrationCmd)
}

func MyUseLine(c *cobra.Command) string {
	originalLine := c.UseLine()
	parts := strings.Split(originalLine, c.Use)
	if len(parts) != 2 {
		return originalLine
	}

	b := &strings.Builder{}
	lines := strings.Split(c.Use, "\n")
	for i, v := range lines {
		if i == 0 {
			fmt.Fprintf(b, "    %s%s%s\n", parts[0], strings.TrimSpace(v), parts[1])
			continue
		}

		fmt.Fprintf(b, "  or: %s%s%s\n", parts[0], strings.TrimSpace(v), parts[1])
	}
	return b.String()
}
