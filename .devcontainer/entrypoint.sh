#!/bin/bash
set -e  # Exit on any error

# The go command fetches the module's dependencies when it first needs them,
# so there is nothing to download at startup.

# Execute the original command passed to the container
exec "$@"
