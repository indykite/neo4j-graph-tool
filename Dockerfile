# Build supervisor and graph-tool with go
FROM golang:1.27-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c AS supervisor
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
FROM neo4j:5.26-enterprise@sha256:43e2256107fea80f589066cf71378c1c87c95e8e52b9f6beaed4a45d8f59a177

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
