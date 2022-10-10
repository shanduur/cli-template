#!/bin/bash -eax

N=$(dasel --plain --file ./.dev/properties.json ".platforms.[#]")
N=$(expr ${N} - 1)

for i in $(seq 0 ${N});
do
    GOOS=$(dasel --plain --file ./.dev/properties.json ".platforms.[$i].os")
    GOARCH=$(dasel --plain --file ./.dev/properties.json ".platforms.[$i].arch")
    GO_SETTINGS="GOOS=${GOOS} GOARCH=${GOARCH}" make build-platform
done
