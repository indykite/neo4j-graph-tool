# Build supervisor and graph-tool with go
FROM golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS supervisor
# Label is used in makefile to delete intermediate images from multistage build
LABEL stage=supervisor_builder
WORKDIR /go/src/github.com/indykite/neo4j-graph-tool
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -buildvcs=false -ldflags "-w -s -extldflags \"-static\"" -o supervisor && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -buildvcs=false -ldflags "-w -s -extldflags \"-static\"" -o graph-tool cli/main.go && \
    chmod u+x ./supervisor ./graph-tool ./entrypoint.sh

# Build final image
FROM neo4j:5.26-enterprise@sha256:96063a2c477d72aa71a019f9fe1c731fb865fd1e0eec8c94dffe8e02aeee629d

COPY --from=supervisor \
    /go/src/github.com/indykite/neo4j-graph-tool/supervisor \
    /go/src/github.com/indykite/neo4j-graph-tool/graph-tool \
    /go/src/github.com/indykite/neo4j-graph-tool/entrypoint.sh \
    /go/src/github.com/indykite/neo4j-graph-tool/config.toml \
    /app/
COPY ./import /initial-data/import

EXPOSE 7474 7473 7687 8080

ENV NEO4J_ACCEPT_LICENSE_AGREEMENT yes
ENTRYPOINT ["tini", "--", "/app/entrypoint.sh"]
