#!/bin/sh

set -e

ARCHITECTURE=amd64
if [ "$(uname -m)" = "aarch64" ]; then
  ARCHITECTURE=arm64
fi  
COMPOSE_SWITCH_VERSION=$(curl -sI https://github.com/docker/compose-switch/releases/latest | grep -E '^(location|Location):' | sed 's,.*/,,' | tr -d '\r')
COMPOSE_SWITCH_URL="https://github.com/docker/compose-switch/releases/download/${COMPOSE_SWITCH_VERSION}/docker-compose-linux-${ARCHITECTURE}"

if ! docker compose version >/dev/null 2>&1; then
  echo "Docker Compose V2 is not installed"
  exit 1
fi

curl -fL "$COMPOSE_SWITCH_URL" -o /usr/local/bin/compose-switch
chmod +x /usr/local/bin/compose-switch

COMPOSE=$(command -v docker-compose || true)
if [ "$COMPOSE" = /usr/local/bin/docker-compose ]; then
  # This is a manual installation of docker-compose
  # so, safe for us to rename binary
  mv /usr/local/bin/docker-compose /usr/local/bin/docker-compose-v1
  COMPOSE=/usr/local/bin/docker-compose-v1
fi

ALTERNATIVES="update-alternatives"
if ! command -v $ALTERNATIVES; then
  ALTERNATIVES=alternatives
fi  

echo "Configuring docker-compose alternatives"
if [ -n "$COMPOSE" ]; then
  $ALTERNATIVES --install /usr/local/bin/docker-compose docker-compose $COMPOSE 1
fi  
$ALTERNATIVES --install /usr/local/bin/docker-compose docker-compose /usr/local/bin/compose-switch 99

echo "'docker-compose' is now set to run Compose V2"
echo "use '$ALTERNATIVES --config docker-compose' if you want to switch back to Compose V1"
