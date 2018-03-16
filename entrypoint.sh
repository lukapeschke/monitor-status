#!/bin/bash

if [ -z $1 ]; then
    echo "Missing USER arg"
    exit 1
fi

echo "Downloading and building dependencies..."

go clean

go get -v .

echo "Building..."

go build

echo "Generating service file..."

USER=$1
cat > monitor-status.service <<EOF
[Unit]
Description=A tool to handle udev events

[Service]
ExecStart=/bin/bash -c "/home/${USER}/bin/monitor-status -config /home/${USER}/.config/monitor-status/config.yml"
User=$USER
Environment=DISPLAY=:0

[Install]
WantedBy=multi-user.target
Alias=monitor-status.service
EOF
