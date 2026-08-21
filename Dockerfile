# Build supervisor and graph-tool with go
FROM golang:1.26-alpine@sha256:28d89ee9cc0ff9fec75c82ca201e6bf7fdf9a679d4b7b24dfa04f2bb766bb468 AS supervisor
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
FROM neo4j:5.26-enterprise@sha256:09528aeece2d6bc0158655f12035cc142fac5b6afd5b38160a1df27fb3f2a156

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
