#!/bin/sh
set -e

LATEST=$(curl -s https://api.github.com/repos/commander-cli/commander/releases/latest | jq -r .tag_name)

if [ -z "$COMMANDER_VER" ]; then
    COMMANDER_VER=${COMMANDER_VER:-$LATEST}
fi
COMMANDER_DST=${COMMANDER_DST:-/usr/local/bin}
INSTALL_LOC="${COMMANDER_DST%/}/commander"
touch "$INSTALL_LOC" || { echo "ERROR: Cannot write to $COMMANDER_DST set COMMANDER_DST elsewhere or use sudo"; exit 1; }

arch=""
if [ "$(uname -m)" = "x86_64" ]; then
    arch="amd64"
elif [ "$(uname -m)" = "aarch64" ]; then
    arch="arm"
else
    arch="386"
fi

url="https://github.com/commander-cli/commander/releases/download/$COMMANDER_VER/commander-linux-$arch"

echo "Downloading $url"
curl -L "$url" -o "$INSTALL_LOC"
chmod +rx "$INSTALL_LOC"
echo "Commander $COMMANDER_VER has been installed to $INSTALL_LOC"
echo "commander --version"
"$INSTALL_LOC" --version
