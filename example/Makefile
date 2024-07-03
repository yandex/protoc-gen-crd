# where local version of protoc-gen-crd module is stored
LOCAL_PATH ?= ..
# where compiled local plugin is stored
LOCAL_PROTOC_GEN_CRD ?= $(LOCAL_PATH)/cmd/protoc-gen-crd/protoc-gen-crd
# command to invoke kubectl
KUBECTL ?= kubectl
# command to invoke go compiler
GO ?= go
# go path
GOPATH ?= $(shell $(GO) env GOPATH)
# go bindir
GOBIN ?= $(GOPATH)/bin
# where installed protoc is stored
PROTOC ?= protoc
# where installed protoc-gen-crd is stored
PROTOC_GEN_CRD ?= $(GOBIN)/protoc-gen-crd
# protoc-gen-crd version to install
PROTOC_GEN_CRD_VERSION ?= latest

PROTO_OPTIONS = vendor/github.com/yandex/protoc-gen-crd/library/go/k8s/protoc_gen_crd/proto/crd.proto
PROTOC_COMMAND = $(PROTOC) --plugin=protoc-gen-crd=$(PROTOC_GEN_CRD) -I=proto -I=vendor --crd_out=paths=source_relative:./proto

all: install-plugin-local vendor-proto compile-crd-server compile-crd-client

# Use this target in this example when protoc plugin module is local.
# In your own projects use install-plugin-remote target.
install-plugin-local:
	cd $(LOCAL_PATH)/cmd/protoc-gen-crd && $(GO) build 
	cp $(LOCAL_PATH)/cmd/protoc-gen-crd/protoc-gen-crd $(GOBIN)/

# Install protoc plugin from github.
install-plugin-remote:
	$(GO) install github.com/yandex/protoc-gen-crd/cmd/protoc-gen-crd@$(PROTOC_GEN_CRD_VERSION)

$(PROTO_OPTIONS):
	$(GO) mod vendor

# Vendor protoc-gen-crd proto definitions
vendor-proto: $(PROTO_OPTIONS)

# Compile server version of CRD.
proto/mycrd.crd.yaml: proto/mycrd.proto
	$(PROTOC_COMMAND) mycrd.proto

compile-crd-server: proto/mycrd.crd.yaml
	
# Compile client version of CRD.
proto/mycrd.kustomize.yaml: proto/mycrd.proto
	$(PROTOC) --plugin=protoc-gen-crd=$(PROTOC_GEN_CRD) -I=proto -I=vendor --crd_opt=client-schema=true --crd_out=paths=source_relative:./proto mycrd.proto

compile-crd-client: proto/mycrd.kustomize.yaml

# install create server yaml to your kubernetes cluster
install: proto/mycrd.crd.yaml
	# NOTE: if you have large model, you may need to avoid 'kubectl apply' since kubectl puts the whole spec copy into metadata.annotations,
	#       and annotation length is limited to 2MB.
	#       In this case you may want to use something like 'kubectl create -f proto/mycrd.crd.yaml || kubectl replace -f proto/mycrd.crd.yaml'
	$(KUBECTL) apply -f proto/mycrd.crd.yaml

# remove generated yaml files
clean:
	rm -f proto/mycrd.crd.yaml proto/mycrd.kustomize.yaml

.PHONY: all install-plugin-local install-plugin-remote vendor-proto compile-crd-server compile-crd-client install clean
