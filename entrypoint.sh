#!/bin/bash
set -e


if [ "$1" = 'web' -o "$1" = 'addstaff' ]; then
	su-exec nobody staffio "$@"
else
	exec "$@"
fi
