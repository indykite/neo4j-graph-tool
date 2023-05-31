# Neo4j Graph Tool

Schema migration tool that helps you automate DB migration across environments and in CI/CD pipelines.

This boilerplate is dependent [Neo4j Graph Tool core](https://github.com/indykite/neo4j-graph-tool-core), which contains the codebase.
And consist of 2 main parts:

- Supervisor - is Docker image entrypoint replacement. Expose HTTP server that can manage underlying Neo4j instance together with running migrations.
- Graph Tool CLI - is binary, that can be used to connect to any Neo4j DB server and perform migrations there.
