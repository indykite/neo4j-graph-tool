GO111MODULE=on
.PHONY: schema_migration local-run neo4j-test-build graph-tool-build upgrade fmt

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))
current_graph_version := $(shell cat .env-file.list | grep GT_DEFAULT_GRAPH_VERSION | sed 's/GT_DEFAULT_GRAPH_VERSION=//')

# Helper to generate migration file
# Usage:
# > make schema_migration NAME=my_migration_name
# If you need different version than current one
# > make schema_migration NAME=my_migration_name VERSION=vX.Y.Z
schema_migration:
	go run cli/main.go genmigration --config=./config.toml schema $(or $(VERSION),$(current_graph_version)) "$(NAME)" $(PARAMS)

local-run:
	@docker run -it --rm -p 7474:7474 -p 7687:7687 -p 8080:8080   \
	--mount type=bind,source=$(current_dir)/import,target=/initial-data/import,readonly \
	--mount type=bind,source=$(current_dir)/config.toml,target=/app/config.toml,readonly \
	--env-file "./.env-file.list" neo4j-test

neo4j-test-build:
	@docker build -t neo4j-test:latest .
	@docker image prune --force --filter label=stage=supervisor_builder

graph-tool-build:
	@docker build -t graph-tool:latest -f cli/Dockerfile .
	@docker image prune --force --filter label=stage=graphtool_builder

upgrade:
	@GO111MODULE=on go get -u all && go mod tidy

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w $$(go list -f "{{.Dir}}" ./... | grep -vE /proto-gen)
