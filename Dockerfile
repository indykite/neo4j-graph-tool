# Build supervisor and graph-tool with go
FROM golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS supervisor
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
FROM neo4j:5.26-enterprise@sha256:739d28a568219b5fd2a0c827c0b43ba54bfc88dbd2b1321d30258ef56c9d07c8

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
