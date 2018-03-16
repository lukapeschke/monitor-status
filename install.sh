#!/bin/sh -x

sudo cp ./monitor-status.service /etc/systemd/system
mkdir -p /home/${USER}/bin/
cp ./monitor-status /home/${USER}/bin/

mkdir -p /home/${USER}/.config/monitor-status/
cp ./config.yml /home/${USER}/.config/monitor-status/
