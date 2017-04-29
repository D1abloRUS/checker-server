#!/bin/sh

if [ $# -eq 0 ]; then
  checker-server
else
  exec "$@"
fi
