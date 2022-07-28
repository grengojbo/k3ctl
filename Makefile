
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GO ?= go
GO_MD2MAN ?= go-md2man

all: manager

major:  ## Set major version
	@#git tag -a v0.0.1 -m 'version v0.0.1' && git push --tags
	@git tag $$(svu major)
	@git push --tags
	@#goreleaser --rm-dist

minor:  ##  Set minor version
	@git tag $$(svu minor)
	@git push --tags
	@#goreleaser --rm-dist

patch:  ##  Set patch version
	@git tag $$(svu patch)
	@git push --tags
	@#goreleaser --rm-dist

downloads: ## Download library
	@curl -L -o kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/v2.3.1/kubebuilder_linux_amd64
	@chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
	@#curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

dep: api ## Get the dependencies
	@$(GO) get -v -d ./...

update: api ## Get and update the dependencies
	@$(GO) get -v -d -u ./...

tidy: ## Clean up dependencies
	@$(GO) mod tidy

vendor: dep ## Create vendor directory
	@$(GO) mod vendor

build: # Build k3ctl binary
	@#go build -o k3ctl cli/main.go
	go build -o k3ctl main.go
	@#goreleaser build --single-target

manager: generate fmt vet
	go build -ldflags "-X main.isCli=no" -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kustomize build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
