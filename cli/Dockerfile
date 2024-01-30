# This Dockerfile must be executed from parent folder with `make graph-tool-build`

FROM golang:1.20-bullseye AS builder
# Regularly update from https://neo4j.com/download-center/#cypher-shell
ARG CYPHER_SHELL_VERSION=5.6.0

# Label is used in makefile to delete intermediate images from multistage build
LABEL stage=graphtool_builder
WORKDIR /go/src/github.com/indykite/neo4j-graph-tool
COPY . .
ENV GO111MODULE=on
# Ensure ca-certficates are up to date
RUN update-ca-certificates && \
    apt-get update && \
    apt-get install -y --no-install-recommends unzip && \
    # Build Graph Tool
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -buildvcs=false -ldflags "-w -s -extldflags \"-static\"" -o graph-tool cli/main.go && \
    # Downlaod CypherShell JAR.
    # And rename so path doesn't contains version. The version is in ARG, which is not accessible in next stage.
    wget -qO/go/cypher-shell.zip https://dist.neo4j.org/cypher-shell/cypher-shell-${CYPHER_SHELL_VERSION}.zip && \
    unzip /go/cypher-shell.zip -d /go && \
    mv /go/cypher-shell-${CYPHER_SHELL_VERSION} /go/cypher-shell


# Second stage, use JRE to minimize Docker image size
FROM eclipse-temurin:21.0.2_13-jre-alpine

WORKDIR /app/

COPY --from=builder --chown=nonroot \
    /go/src/github.com/indykite/neo4j-graph-tool/graph-tool \
    /go/src/github.com/indykite/neo4j-graph-tool/config.toml \
    /app/
COPY --from=builder --chown=nonroot /go/cypher-shell /app/cypher-shell
COPY ./import /app/import

ENTRYPOINT [ "/app/graph-tool" ]
