# only set those variables if are not already set
NAME ?= $(shell dasel --plain --file ./.dev/properties.json '.name')
VERSION ?= $(shell dasel --plain --file ./.dev/properties.json '.version')
REVISION ?= $(shell git rev-parse --short HEAD)$(shell git diff --quiet || echo '-dirty')

REGISTRY := $(shell dasel --plain --file ./.dev/properties.json '.registry')
REPOSITORY := $(shell dasel --plain --file ./.dev/properties.json '.repository')

TOOLCHAIN_VERSION := $(shell sed -En 's/^go (.*)$$/\1/p' go.mod)
MODULE_NAME := $(shell sed -En 's/^module (.*)$$/\1/p' go.mod)
GO_SETTINGS += CGO_ENABLED=0

LDFLAGS += -X ${MODULE_NAME}/cmd/version.Version=${VERSION}
LDFLAGS += -X ${MODULE_NAME}/cmd/version.Revision=${REVISION}

OCI_BUILDARGS += --build-arg=TOOLCHAIN_VERSION=${TOOLCHAIN_VERSION}
OCI_BUILDARGS += --build-arg=NAME=${NAME}
OCI_BUILDARGS += --build-arg=REVISION=${REVISION}
OCI_BUILDARGS += --build-arg=VERSION=${VERSION}

OCI_TAGS += --tag=${REGISTRY}/${REPOSITORY}/${NAME}:latest
OCI_TAGS += --tag=${REGISTRY}/${REPOSITORY}/${NAME}:${VERSION}

BUILDAH_ARCH ?= $(shell uname -m)
BUILDAH_MANIFEST := ${NAME}-manifest

.PHONY: test
test:
	go test \
		-cover \
		./...

.PHONY: clean
clean:
	rm -rf ./build
	rm Containerfile

.PHONY: install
install: build
	sudo install ./build/${NAME} /usr/local/bin/${NAME}

.PHONY: build
build: build-platform
	cp ./build/$(shell go env GOOS)/$(shell go env GOARCH)/${NAME} ./build/${NAME}

.PHONY: build-platform
build-platform: test
	${GO_SETTINGS} go build \
		-ldflags="${LDFLAGS}" \
		-o ./build/$(shell ${GO_SETTINGS} go env GOOS)/$(shell ${GO_SETTINGS} go env GOARCH)/${NAME} \
		main.go

.PHONY: build-all
build-all: test
	./scripts/build-matrix.sh

.PHONY: hooks
hooks:
	pre-commit install --hook-type pre-commit
	pre-commit install --hook-type commit-msg

# no phony, as we generate this only once
Containerfile:
	sed "s/{{.APPNAME}}/${NAME}/g" Containerfile.template > Containerfile

.PHONY: docker
docker:
	ENGINE=docker make container-image

.PHONY: docker-clean
docker-clean:
	ENGINE=docker make container-image

.PHONY: podman
podman:
	ENGINE=podman make container-image

.PHONY: podman-clean
podman-clean:
	ENGINE=podman make container-image-clean

.PHONY: container-image
container-image: container-image-clean Containerfile
	${ENGINE} build \
		${OCI_TAGS} \
		${OCI_BUILDARGS} \
		.

# discarding errors, because image might not exist
.PHONY: container-image-clean
container-image-clean:
	@-${ENGINE} image rm -f $(shell ${ENGINE} image ls -aq ${REGISTRY}${NAME}:${VERSION} | xargs -n1 | sort -u | xargs)

# utilizes docker buildx for building multiarch images
.PHONY: docker-push
docker-push: Containerfile
	docker buildx build \
		${OCI_TAGS} \
		${OCI_BUILDARGS} \
		--output=type=registry \
		--platform=linux/amd64,linux/arm64 \
		.

.PHONY: buildah-manifest
buildah-manifest:
	buildah manifest create ${BUILDAH_MANIFEST}

.PHONY: buildah-build
buildah-build: Containerfile
	buildah bud \
		${OCI_TAGS} \
		${OCI_BUILDARGS} \
		--manifest ${BUILDAH_MANIFEST} \
		--arch ${BUILDAH_ARCH} \
		.

# utilizes buildah for building multiarch images
.PHONY: buildah-push
buildah-push: buildah-manifest
	./scripts/buildah-matrix.sh
	buildah manifest push --all \
		${BUILDAH_MANIFEST} \
		"docker://${REGISTRY}/${REPOSITORY}/${NAME}:latest"
	buildah manifest push --all \
		${BUILDAH_MANIFEST} \
		"docker://${REGISTRY}/${REPOSITORY}/${NAME}:${VERSION}"
