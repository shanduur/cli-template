VERSION 0.6
FROM golang:1.19
WORKDIR /work

deps:
    RUN go install github.com/tomwright/dasel/cmd/dasel@latest

properties:
    FROM +deps
    COPY --dir .dev .

modules:
    FROM +properties
    COPY go.mod go.sum .
    RUN go mod download

# Copies the code to the buildkit
code:
    FROM +modules

    COPY --dir cmd .
    COPY *.go .
    COPY --dir .git .

test:
    FROM +code
    RUN go test \
        -cover \
        -race \
        -covermode=atomic \
        -coverprofile=coverage.out \
        ./...
    SAVE ARTIFACT coverage.out

codecov:
    FROM scratch
    ARG CODECOV_TOKEN
    COPY +test/coverage.out .
    RUN curl -Os https://uploader.codecov.io/latest/linux/codecov && \
        chmod +x codecov && \
        ./codecov -t ${CODECOV_TOKEN}

lint:
    FROM golangci/golangci-lint:v1.50.1
    WORKDIR /work
    COPY .golangci.yml .
    COPY --dir cmd .
    COPY go.mod go.sum ./*.go .
    RUN golangci-lint run -v

build:
    FROM +code

    ARG GOOS
    ARG GOARCH
    ARG CGO_ENABLED=0

    ARG NAME="$(dasel --plain --file ./.dev/properties.json '.name')"
    ARG VERSION="$(dasel --plain --file ./.dev/properties.json '.version')"
    ARG REVISION="$(git rev-parse --short HEAD)$(git diff --quiet || echo '-dirty')"
    ARG MODULE_NAME="$(sed -En 's/^module (.*)$/\1/p' go.mod)"
    ARG EXTENSION

    ARG LDFLAGS="
        -X ${MODULE_NAME}/cmd/version.Version=${VERSION}
        -X ${MODULE_NAME}/cmd/version.Revision=${REVISION}
    "

    RUN go build \
        -ldflags="${LDFLAGS}" \
        -o build/${NAME}-${GOOS}-${GOARCH}${EXTENSION}

    SAVE ARTIFACT build/${NAME}-${GOOS}-${GOARCH}${EXTENSION}

save:
    FROM +properties

    ARG GOOS
    ARG GOARCH
    ARG EXTENSION

    ARG NAME="$(dasel --plain --file ./.dev/properties.json '.name')"

    COPY +build/${NAME}-${GOOS}-${GOARCH}${EXTENSION} ${NAME}
    SAVE ARTIFACT ${NAME} AS LOCAL ./build/${NAME}

# builds platform specific binaries
for-darwin:
    BUILD +build --GOOS=darwin --GOARCH=arm64
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=darwin --GOARCH=arm64
    END

for-darwin-legacy:
    BUILD +build --GOOS=darwin --GOARCH=amd64
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=darwin --GOARCH=amd64
    END


for-linux-arm64:
    BUILD +build --GOOS=linux --GOARCH=arm64
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=linux --GOARCH=arm64
    END

for-linux-amd64:
    BUILD +build --GOOS=linux --GOARCH=amd64
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=linux --GOARCH=amd64
    END

for-windows:
    BUILD +build --GOOS=windows --GOARCH=amd64 --EXTENSION=.exe
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=windows --GOARCH=amd64 --EXTENSION=.exe
    END


for-windows-on-arm:
    BUILD +build --GOOS=windows --GOARCH=arm64 --EXTENSION=.exe
    IF [ -z ${MATRIX} ]
        BUILD +save --GOOS=windows --GOARCH=arm64 --EXTENSION=.exe
    END

build-all:
    BUILD +for-darwin --MATRIX=1
    BUILD +for-darwin-legacy --MATRIX=1
    BUILD +for-linux-arm64 --MATRIX=1
    BUILD +for-linux-amd64 --MATRIX=1
    BUILD +for-windows --MATRIX=1
    BUILD +for-windows-on-arm --MATRIX=1

images:
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        +image

# Builds container image
image:
    FROM +properties

    ARG GOOS
    ARG GOARCH
    ARG EXTENSION

    ARG REGISTRY="$(dasel --plain --file ./.dev/properties.json '.registry')"
    ARG REPOSITORY="$(dasel --plain --file ./.dev/properties.json '.repository')"
    ARG NAME="$(dasel --plain --file ./.dev/properties.json '.name')"
    ARG VERSION="$(dasel --plain --file ./.dev/properties.json '.version')"

    FROM scratch
    WORKDIR /
    COPY +build/${NAME}-${GOOS}-${GOARCH}${EXTENSION} /usr/bin/app

    ENTRYPOINT [ "/usr/bin/app" ]
    CMD []
    SAVE IMAGE --push ${REGISTRY}/${REPOSITORY}/${NAME}:latest ${REGISTRY}/${REPOSITORY}/${NAME}:${VERSION}
