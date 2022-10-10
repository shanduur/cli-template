#!/bin/bash -eax

N=$(dasel --plain --file ./.dev/properties.json ".platforms.[#]")
N=$(expr ${N} - 1)

for i in $(seq 0 ${N});
do
    OS=$(dasel --plain --file ./.dev/properties.json ".platforms.[$i].os")
    if [[ "${OS}" == 'linux' ]]; then
        BUILDAH_ARCH=$(dasel --plain --file ./.dev/properties.json ".platforms.[$i].arch") make build
    fi
done
