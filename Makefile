GIT_REPO=github.com/nutanix/patrao
BIN_NAME=upgradeagent
IMG_NAME=patrao-upgrade-agent
IMG_VERSION=latest
rdir := $(dir $(lastword $(MAKEFILE_LIST)))

.PHONY: clean binaries

all: image

binaries: $(rdir)/bin/$(BIN_NAME)

image: binaries
	-docker rmi $(IMG_NAME):$(IMG_VERSION) >/dev/null 2>&1
	(docker build -t $(IMG_NAME):$(IMG_VERSION) -f $(rdir)/deployments/$(BIN_NAME)/Dockerfile .)

$(rdir)/bin/$(BIN_NAME):
	mkdir -p $(rdir)/bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $(GIT_REPO)/cmd/$(BIN_NAME)
	go test $(GIT_REPO)/internal/app/upgradeagent/

clean:
	-docker rmi $(IMG_NAME):$(IMG_VERSION)
	rm -f $(rdir)/bin/$(BIN_NAME)
