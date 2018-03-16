==============
monitor-status
==============

``monitor-status`` (please do a PR if you have a better name) is a tool which
allows to watch for udev events. It was initially written to reset my i3
worspaces when a monitor is plugged in, but it can be used for any kind of
device.

It's still work in progress.

Installing and running
======================

There are no binary packages yet, but the tootl can easily be built.
Dependencies are ``go`` (1.8 or higher version) and ``libudev-dev``
(``systemd-devel`` on Fedora).

Once the depencies are installed, simply run ``go get
github.com/lukapeschke/monitor-status`` and you're good to go.

Configuration
=============

The tool will look for a config file in the following places::

    "./config.yml",
    "~/.monitor-status/config.yml",
    "~/.config/monitor-status.yml",
    "~/.config/monitor-status/config.yml",
    "$GOPATH/src/github.com/lukapeschke/monitor-status/config.yml"

If your config file is located elsewhere, you can specify its path with the
``-config`` flag.

The config file has the following format::

    drm: # subsystem
      card0-DP-3: # device
        on_connect: ~/bin/reset_workspaces.sh # command for "add" event
        on_disconnect: xrandr --output DP-1-1 --off # command for "remove" event
      DEVICE2:
        on_connect: echo "hello world!"
        on_disconnect: echo "bye!"
    SUBSYSTEM2:
      [...]

For now, only ``on_connect`` and ``on_disconnect`` are supported, but more
events will be supported in the future.
