#!/usr/bin/env bash

# abort on errors and more error-checking.
set -euo pipefail
set -o noclobber    # Avoid overlay files (echo "hi" > foo)
set -o errexit      # Used to exit upon error, avoiding cascading errors
set -o pipefail     # Unveils hidden failures
set -o nounset      # Exposes unset variables
shopt -s nullglob   # Non-matching globs are removed  ('*.foo' => '')
shopt -s failglob   # Non-matching globs throw errors
shopt -s nocaseglob # Case insensitive globs
shopt -s dotglob    # Wildcards match dotfiles ("*.sh" => ".foo.sh")
shopt -s globstar   # Allow ** for recursive matches ('lib/**/*.rb' => 'lib/a/b/c.rb')

DCR_VERSION="latest"
CONFIG_PATH="${CONFIG_PATH:-$(pwd)/configs/config.json}" # If `CONFIG_PATH` not set or null, use default.

printf "%b" "\033[92m" "scripts/docker-run-stable.sh: running openbanking/conformance-dcr:${DCR_VERSION}, CONFIG_PATH=${CONFIG_PATH}" "\033[0m" "\n"
docker --log-level=debug run \
    --rm \
    -it \
    -v "${CONFIG_PATH}":/home/app/.config/conformance-dcr/config.json \
    openbanking/conformance-dcr:"${DCR_VERSION}" \
        -debug \
        -report \
        -config-path=/home/app/.config/conformance-dcr/config.json