# Build the kubectl-onmetal binary
FROM golang:1.20 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY bootstraptoken/ bootstraptoken/
COPY cmd/ cmd/
COPY utils/ utils/

ARG TARGETOS
ARG TARGETARCH

# Build
RUN mkdir bin

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -a -o bin/kubectl-onmetal main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/kubectl-onmetal .
USER 65532:65532

ENTRYPOINT ["/kubectl-onmetal"]
