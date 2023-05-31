#!/bin/sh
set -e

# If there are only flags or nothing in args, set "supervisor".
# Otherwise run whatever was specified, sh for example
# ${1#-} strips "-" from the beginning of first argument
if [ "${1#-}" != "$1" ] || [ -z "$1" ]; then
	set -- /app/supervisor "$@"
fi

exec "$@"
