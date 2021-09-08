#!/bin/sh

COMPOSE_SWITCH_VERSION="1.0.0"
COMPOSE_SWITCH_URL="https://github.com/docker/compose-switch/releases/download/$COMPOSE_SWITCH_VERSION/docker-compose"
DOCKER_COMPOSE_V1_PATH_ORIG="$(which docker-compose)"
DOCKER_COMPOSE_V1_PATH_NEW="$DOCKER_COMPOSE_V1_PATH_ORIG-v1"

if [ ! -f "${DOCKER_COMPOSE_V1_PATH_ORIG}" ]; then
  echo "'docker-compose' V1 could not be found in PATH. Aborting..."
  exit 1
fi

if [ -f "${DOCKER_COMPOSE_V1_PATH_NEW}" ]; then
  echo "Looks like docker-compose V1->V2 switch is already installed. 'docker-compose' binary was already moved to '${DOCKER_COMPOSE_V1_PATH_NEW}'."
  echo "If you are trying to re-enable Docker Compose V2, please run:"
  echo "$ docker-compose enable-v2"
  exit 0
fi

error=$(docker compose version 2>&1 >/dev/null)
if [ $? -ne 0 ]; then
  echo "Docker Compose V2 is not installed"
  exit 1
fi

TMP_FILE=$(mktemp)
echo "Downloading compose V1->V2 switch into $TMP_FILE..."
curl -L -o "$TMP_FILE" "$COMPOSE_SWITCH_URL"
chmod +x "$TMP_FILE"

echo "Switching Docker Compose binary"
mv "$DOCKER_COMPOSE_V1_PATH_ORIG" "$DOCKER_COMPOSE_V1_PATH_NEW"
mv "$TMP_FILE" "$DOCKER_COMPOSE_V1_PATH_ORIG"

echo "Compose V1->V2 installed"

$DOCKER_COMPOSE_V1_PATH_ORIG enable-v2
echo "Docker Compose V2 enabled"
docker-compose version

echo "To switch back to Compose V1, you can run: `docker-compose disable-v2`"
