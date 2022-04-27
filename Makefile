
CUR_SHA=$(shell git log -n1 --pretty='%h')
CUR_BRANCH=$(shell git branch --show-current)
VERSION=$(shell git describe --exact-match --tags $(CUR_SHA) 2>/dev/null || echo $(CUR_BRANCH)-$(CUR_SHA))

CNT_IMAGE=pathwae/pathwae

CC=podman

image:
	$(CC) build -t $(CNT_IMAGE):$(VERSION) --build-arg VERSION=$(VERSION) -f docker/Dockerfile .
	$(CC) tag $(CNT_IMAGE):$(VERSION) $(CNT_IMAGE):latest

