#!/bin/sh

set -e


ARCHITECTURE=amd64
if [ "$(uname -m)" = "aarch64" ]; then
  ARCHITECTURE=arm64
fi  
COMPOSE_SWITCH_VERSION="v1.0.2"
COMPOSE_SWITCH_URL="https://github.com/docker/compose-switch/releases/download/${COMPOSE_SWITCH_VERSION}/docker-compose-linux-${ARCHITECTURE}"

error=$(docker compose version 2>&1 >/dev/null)
echo "After error"
if [ $? -ne 0 ]; then
  echo "Docker Compose V2 is not installed"
  exit 1
fi

sh_c='sh -c'
sudo_sh_c='sh -c'
if [ "$user" != 'root' ]; then
    if [ "$(command -v sudo)" ]; then
        sudo_sh_c='sudo -E sh -c'
    elif [ "$(command -v su)" ]; then
        sudo_sh_c='su -c'
    else
        echo "Error: This installer needs the ability to run commands as root."
        exit 1
    fi
fi

$sudo_sh_c "curl -fL $COMPOSE_SWITCH_URL -o /usr/local/bin/compose-switch"
$sudo_sh_c "chmod +x /usr/local/bin/compose-switch"

COMPOSE=$(command -v docker-compose)
if [ "$COMPOSE" = /usr/local/bin/docker-compose ]; then
  # This is a manual installation of docker-compose
  # so, safe for us to rename binary
  $sudo_sh_c "mv /usr/local/bin/docker-compose /usr/local/bin/docker-compose-v1"
  COMPOSE=/usr/local/bin/docker-compose-v1
fi


ALTERNATIVES="$sudo_sh_c update-alternatives"
if ! command -v $ALTERNATIVES; then
  ALTERNATIVES="$sudo_sh_c alternatives"
fi  

echo "Configuring `docker-compose` alternatives"
if [ ! -z $COMPOSE ]; then
  $ALTERNATIVES --install /usr/local/bin/docker-compose docker-compose $COMPOSE 1
fi  
$ALTERNATIVES --install /usr/local/bin/docker-compose docker-compose /usr/local/bin/compose-switch 99

echo "'docker-compose' is now set to run Compose V2"
echo "use '$ALTERNATIVES --config docker-compose' if you want to switch back to Compose V1"
