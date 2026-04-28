#!/bin/sh

set -e

IMAGE_NAME="safezoneturingpe/garage-webui"
PACKAGE_VERSION=$(cat package.json | grep \"version\" | cut -d'"' -f 4)

echo "Building version $PACKAGE_VERSION"

# El primer argumento del script tiene prioridad, luego la variable de entorno DOCKER_ACTION,
# y por defecto "load"
ACTION="${1:-${DOCKER_ACTION:-load}}"

case "$ACTION" in
    load)
        BUILD_FLAG="--load"
        ;;
    push)
        BUILD_FLAG="--push"
        ;;
    *)
        echo "Acción inválida: $ACTION. Usa 'load' o 'push' (o define la variable DOCKER_ACTION)."
        exit 1
        ;;
esac

docker buildx build --platform linux/amd64 \
 -t "$IMAGE_NAME:latest" -t "$IMAGE_NAME:$PACKAGE_VERSION" $BUILD_FLAG .