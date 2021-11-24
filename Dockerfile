# Build the kubectl-onmetal binary
FROM golang:1.17 as builder

RUN apt-get update -yq \
    && apt-get install -yq --no-install-recommends libvirt-dev

ARG GOARCH=''
ARG GITHUB_PAT=''

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

COPY hack hack

ENV GOPRIVATE='github.com/onmetal/*'

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=ssh --mount=type=secret,id=github_pat GITHUB_PAT_PATH=/run/secrets/github_pat ./hack/setup-git-redirect.sh \
  && mkdir -p -m 0600 ~/.ssh \
  && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts \
  && go mod download

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY exec/ exec/

# Build
RUN GOMAXPROCS=1 CGO_ENABLED=1 GOOS=linux GOARCH=${GOARCH:-$(go env GOARCH)} go build -a -o kubectl-onmetal main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kubectl-onmetal .
USER 65532:65532

ENTRYPOINT ["/kubectl-onmetal"]
